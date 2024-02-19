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
	obj := object.New().WithID(fmt.Sprintf("player_%d", pIndex))
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Pos = pos
	obj.Layer = 27 - pIndex*2
	player := data.NewDynamic()
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	switch pIndex {
	case 0:
		player.Control = controllers.NewPlayerInput(data.P1Input)
		e.AddComponent(myecs.Controller, player.Control)
		player.Anim = reanimator.HumanoidAnimation(player, "player1")
		player.Color = constants.StrColorBlue
	case 1:
		player.Control = controllers.NewPlayerInput(data.P2Input)
		e.AddComponent(myecs.Controller, player.Control)
		player.Anim = reanimator.HumanoidAnimation(player, "player2")
		player.Color = constants.StrColorGreen
	case 2:
		player.Control = controllers.NewPlayerInput(data.P3Input)
		e.AddComponent(myecs.Controller, player.Control)
		player.Anim = reanimator.HumanoidAnimation(player, "player3")
		player.Color = constants.StrColorPurple
	case 3:
		player.Control = controllers.NewPlayerInput(data.P4Input)
		e.AddComponent(myecs.Controller, player.Control)
		player.Anim = reanimator.HumanoidAnimation(player, "player4")
		player.Color = constants.StrColorBrown
	}
	player.Object = obj
	player.Entity = e
	player.Player = data.Player(pIndex)
	player.Vars = data.PlayerVars()
	e.AddComponent(myecs.Animated, player.Anim)
	e.AddComponent(myecs.Drawable, player.Anim)
	e.AddComponent(myecs.Dynamic, player)
	e.AddComponent(myecs.Player, player.Player)
	return player
}

func DemonCharacter(pos pixel.Vec) *data.Dynamic {
	obj := object.New().WithID("demon")
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Pos = pos
	obj.Layer = 29
	demon := data.NewDynamic()
	demon.Control = controllers.NewLRChase(demon)
	demon.Object = obj
	demon.Anim = reanimator.HumanoidAnimation(demon, "demon")
	demon.Vars = data.DemonVars()
	e := myecs.Manager.NewEntity()
	demon.Entity = e
	e.AddComponent(myecs.Object, demon.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Animated, demon.Anim)
	e.AddComponent(myecs.Drawable, demon.Anim)
	e.AddComponent(myecs.Dynamic, demon)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, demon.Control)
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
		if (enemy.State == data.Grounded || enemy.State == data.Ladder || enemy.State == data.Leaping || enemy.State == data.Flying) &&
			(ch.State == data.Grounded || ch.State == data.Ladder || ch.State == data.Leaping || ch.State == data.Jumping) {
			ch.Flags.Hit = true
			ch.State = data.Hit
			enemy.Flags.Attack = true
			enemy.State = data.Attack
		}
	}
}

func FlyCharacter(pos pixel.Vec, left bool) *data.Dynamic {
	obj := object.New().WithID("fly")
	obj.SetRect(pixel.R(0, 0, 12, 12))
	obj.Pos = pos
	obj.Flip = left
	obj.Layer = 29
	fly := data.NewDynamic()
	fly.Control = controllers.NewBackAndForth(fly, left)
	fly.State = data.Flying
	fly.Flags.Flying = true
	fly.Object = obj
	fly.Anim = reanimator.FlyAnimation(fly)
	fly.Vars = data.FlyVars()
	e := myecs.Manager.NewEntity()
	fly.Entity = e
	e.AddComponent(myecs.Object, fly.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Animated, fly.Anim)
	e.AddComponent(myecs.Drawable, fly.Anim)
	e.AddComponent(myecs.Dynamic, fly)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, fly.Control)
	return fly
}
