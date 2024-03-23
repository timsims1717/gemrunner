package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/controllers"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems/reanimator"
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func PlayerCharacter(pos pixel.Vec, pIndex int) *data.Dynamic {
	obj := object.New().WithID(fmt.Sprintf("player_%d", pIndex)).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Layer = 27 - pIndex*2
	player := data.NewDynamic()
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	SetAsPlayer(player, e, pIndex)
	player.Object = obj
	player.Entity = e
	player.Player = pIndex
	player.Vars = data.PlayerVars()
	e.AddComponent(myecs.Animated, player.Anim)
	e.AddComponent(myecs.Drawable, player.Anim)
	e.AddComponent(myecs.Dynamic, player)
	e.AddComponent(myecs.Player, player.Player)
	e.AddComponent(myecs.LvlElement, struct{}{})
	pickUp := data.NewPickUp(11, true)
	pickUp.NoInventory = true
	e.AddComponent(myecs.PickUp, pickUp)
	return player
}

func SetAsPlayer(ch *data.Dynamic, e *ecs.Entity, p int) {
	switch p {
	case 0:
		ch.Control = controllers.NewPlayerInput(data.P1Input, e)
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anim = reanimator.HumanoidAnimation(ch, "player1")
		ch.Color = constants.StrColorBlue
	case 1:
		ch.Control = controllers.NewPlayerInput(data.P2Input, e)
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anim = reanimator.HumanoidAnimation(ch, "player2")
		ch.Color = constants.StrColorGreen
	case 2:
		ch.Control = controllers.NewPlayerInput(data.P3Input, e)
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anim = reanimator.HumanoidAnimation(ch, "player3")
		ch.Color = constants.StrColorPurple
	case 3:
		ch.Control = controllers.NewPlayerInput(data.P4Input, e)
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anim = reanimator.HumanoidAnimation(ch, "player4")
		ch.Color = constants.StrColorBrown
	}
	data.CurrLevel.Players[p] = ch
	data.CurrLevel.Stats[p] = data.NewStats()
	data.CurrLevel.PControls[p] = ch.Control
}

func DemonCharacter(pos pixel.Vec, metadata data.TileMetadata) *data.Dynamic {
	obj := object.New().WithID("demon").SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Layer = 29
	demon := data.NewDynamic()
	demon.Object = obj
	demon.Anim = reanimator.HumanoidAnimation(demon, "demon")
	demon.State = data.Regen
	demon.Flags.Regen = true
	demon.Options.Regen = metadata.Regenerate
	demon.Vars = data.DemonVars()
	e := myecs.Manager.NewEntity()
	demon.Entity = e
	demon.Control = controllers.NewLRChase(demon, e)
	e.AddComponent(myecs.Object, demon.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Animated, demon.Anim)
	e.AddComponent(myecs.Drawable, demon.Anim)
	e.AddComponent(myecs.Dynamic, demon)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, demon.Control)
	e.AddComponent(myecs.LvlElement, struct{}{})
	data.CurrLevel.Enemies = append(data.CurrLevel.Enemies, demon)
	return demon
}

func KillPlayer(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
	if p < 0 ||
		p >= constants.MaxPlayers ||
		ch.State == data.Hit ||
		ch.State == data.Dead {
		return
	}
	bg, ok := entity.GetComponentData(myecs.Dynamic)
	if ok {
		enemy := bg.(*data.Dynamic)
		if (enemy.State == data.Grounded ||
			enemy.State == data.OnLadder ||
			enemy.State == data.Leaping ||
			enemy.State == data.Flying ||
			enemy.State == data.Carried) &&
			(ch.State == data.Grounded ||
				ch.State == data.OnLadder ||
				ch.State == data.Leaping ||
				ch.State == data.Jumping ||
				ch.State == data.Flying ||
				ch.State == data.Carried ||
				ch.State == data.Thrown) {
			ch.Flags.Hit = true
			ch.State = data.Hit
			enemy.Flags.Attack = true
			enemy.State = data.Attack
		}
	}
}

func FlyCharacter(pos pixel.Vec, left bool) *data.Dynamic {
	obj := object.New().WithID("fly").SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 12, 12))
	obj.Flip = left
	obj.Layer = 29
	fly := data.NewDynamic()
	fly.State = data.Flying
	fly.Flags.Flying = true
	fly.Object = obj
	fly.Anim = reanimator.FlyAnimation(fly)
	fly.Vars = data.FlyVars()
	e := myecs.Manager.NewEntity()
	fly.Entity = e
	fly.Control = controllers.NewBackAndForth(fly, e, left)
	e.AddComponent(myecs.Object, fly.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Animated, fly.Anim)
	e.AddComponent(myecs.Drawable, fly.Anim)
	e.AddComponent(myecs.Dynamic, fly)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, fly.Control)
	e.AddComponent(myecs.LvlElement, struct{}{})
	data.CurrLevel.Enemies = append(data.CurrLevel.Enemies, fly)
	return fly
}

func SetEmptyControl(ch *data.Dynamic) {
	ch.Entity.AddComponent(myecs.Controller, controllers.NewEmpty(ch, ch.Entity))
}

func SetPlayerControl(ch *data.Dynamic, p int) {
	if data.CurrLevel.PControls[p] == nil {
		var pc *controllers.PlayerInput
		switch p {
		case 0:
			pc = controllers.NewPlayerInput(data.P1Input, ch.Entity)
		case 1:
			pc = controllers.NewPlayerInput(data.P2Input, ch.Entity)
		case 2:
			pc = controllers.NewPlayerInput(data.P3Input, ch.Entity)
		case 3:
			pc = controllers.NewPlayerInput(data.P4Input, ch.Entity)
		}
		ch.Entity.AddComponent(myecs.Controller, pc)
		data.CurrLevel.PControls[p] = pc
	} else {
		SetEmptyControl(data.CurrLevel.Players[p])
		ch.Entity.AddComponent(myecs.Controller, data.CurrLevel.PControls[p])
	}
}

func ResetControl(ch *data.Dynamic) {
	ch.Entity.AddComponent(myecs.Controller, ch.Control)
}
