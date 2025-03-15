package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/controllers"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems/animations"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func PlayerCharacter(pos pixel.Vec, pIndex int, tile *data.Tile, replay *data.LevelReplay) *data.Dynamic {
	player := data.NewDynamic(tile)
	player.Layer = 27 - pIndex*2
	obj := object.New().WithID(fmt.Sprintf("player_%d", pIndex)).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Layer = player.Layer
	PlayerPortal(obj.Layer+1, pos)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	SetAsPlayer(player, e, pIndex, replay)
	player.Object = obj
	player.Entity = e
	player.Player = pIndex
	player.Vars = data.PlayerVars()
	player.State = data.Regen
	player.Flags.Regen = true
	player.Options.Regen = true
	player.Options.StoredCount = 12
	player.Options.LinkedTiles = []world.Coords{world.NewCoords(world.WorldToMap(pos.X, pos.Y))}
	dsg := animations.PlayerAnimation(player, "disguise", false)
	dsg.Dependent = true
	player.Anims.Set = append(player.Anims.Set, dsg)
	e.AddComponent(myecs.Animated, player.Anims)
	e.AddComponent(myecs.Drawable, player.Anims)
	e.AddComponent(myecs.Dynamic, player)
	e.AddComponent(myecs.Player, player.Player)
	e.AddComponent(myecs.LvlElement, struct{}{})
	return player
}

func SetAsPlayer(ch *data.Dynamic, e *ecs.Entity, p int, replay *data.LevelReplay) {
	switch p {
	case 0:
		if replay != nil {
			ch.Control = controllers.NewReplayController(replay, p, e)
		} else {
			ch.Control = controllers.NewPlayerInput(data.P1Input, e)
		}
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anims.Add(animations.PlayerAnimation(ch, "player1", true))
		ch.Color = data.PlayerBlue
	case 1:
		if replay != nil {
			ch.Control = controllers.NewReplayController(replay, p, e)
		} else {
			ch.Control = controllers.NewPlayerInput(data.P2Input, e)
		}
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anims.Add(animations.PlayerAnimation(ch, "player2", true))
		ch.Color = data.PlayerGreen
	case 2:
		if replay != nil {
			ch.Control = controllers.NewReplayController(replay, p, e)
		} else {
			ch.Control = controllers.NewPlayerInput(data.P3Input, e)
		}
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anims.Add(animations.PlayerAnimation(ch, "player3", true))
		ch.Color = data.PlayerPurple
	case 3:
		if replay != nil {
			ch.Control = controllers.NewReplayController(replay, p, e)
		} else {
			ch.Control = controllers.NewPlayerInput(data.P4Input, e)
		}
		e.AddComponent(myecs.Controller, ch.Control)
		ch.Anims.Add(animations.PlayerAnimation(ch, "player4", true))
		ch.Color = data.PlayerOrange
	}
	data.CurrLevel.Players[p] = ch
	data.CurrLevel.PControls[p] = ch.Control
}

func DemonCharacter(pos pixel.Vec, tile *data.Tile) *data.Dynamic {
	metadata := tile.Metadata
	demon := data.NewDynamic(tile)
	obj := object.New().WithID("demon").SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 12, 16))
	demon.Layer = 29
	obj.Layer = demon.Layer
	demon.Object = obj
	demon.Anims = reanimator.NewSet().Add(animations.DemonAnimation(demon))
	demon.State = data.Regen
	demon.Flags.Regen = true
	demon.Options.Regen = metadata.Regenerate
	demon.Options.LinkedTiles = metadata.LinkedTiles
	demon.Vars = data.DemonVars()
	e := myecs.Manager.NewEntity()
	demon.Entity = e
	demon.Control = controllers.NewLRChase(demon, e)
	e.AddComponent(myecs.Object, demon.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Animated, demon.Anims)
	e.AddComponent(myecs.Drawable, demon.Anims)
	e.AddComponent(myecs.Dynamic, demon)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, demon.Control)
	e.AddComponent(myecs.LvlElement, struct{}{})
	demon.Enemy = len(data.CurrLevel.Enemies)
	data.CurrLevel.Enemies = append(data.CurrLevel.Enemies, demon)
	return demon
}

func KillPlayer(p int, ch *data.Dynamic, entity *ecs.Entity) {
	if p < 0 ||
		p >= constants.MaxPlayers ||
		ch.Flags.Disguised ||
		ch.Flags.ItemAction == data.TransportIn ||
		ch.State == data.Hiding ||
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
			enemy.State == data.Flying) &&
			(ch.State == data.Grounded ||
				ch.State == data.OnLadder ||
				ch.State == data.Leaping ||
				ch.State == data.Jumping ||
				(ch.State == data.Falling && enemy.Flags.Flying) ||
				ch.State == data.Flying ||
				ch.State == data.DoingAction) {
			ch.Flags.Hit = true
			ch.State = data.Hit
			enemy.Flags.Attack = true
			enemy.State = data.Attack
		}
	}
}

func FlyCharacter(pos pixel.Vec, tile *data.Tile) *data.Dynamic {
	metadata := tile.Metadata
	fly := data.NewDynamic(tile)
	obj := object.New().WithID("fly").SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 12, 12))
	obj.Flip = metadata.Flipped
	fly.Layer = 29
	obj.Layer = fly.Layer
	fly.State = data.Flying
	fly.Options.RegenFlip = true
	fly.Options.Flying = true
	fly.Flags.Flying = true
	fly.Options.Regen = metadata.Regenerate
	fly.Options.LinkedTiles = metadata.LinkedTiles
	fly.State = data.Regen
	fly.Flags.Regen = true
	fly.Object = obj
	fly.Anims = reanimator.NewSet().Add(animations.FlyAnimation(fly))
	fly.Vars = data.FlyVars()
	e := myecs.Manager.NewEntity()
	fly.Entity = e
	fly.Control = controllers.NewBackAndForth(fly, e, metadata.Flipped)
	e.AddComponent(myecs.Object, fly.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Animated, fly.Anims)
	e.AddComponent(myecs.Drawable, fly.Anims)
	e.AddComponent(myecs.Dynamic, fly)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, fly.Control)
	e.AddComponent(myecs.LvlElement, struct{}{})
	fly.Enemy = len(data.CurrLevel.Enemies)
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

func PlayerPortal(layer int, pos pixel.Vec) {
	obj := object.New()
	obj.Pos = pos
	obj.Layer = layer
	first := true
	m := myecs.Manager.NewEntity()
	a := reanimator.NewBatchAnimation("portal", img.Batchers[constants.TileBatch], "portal_magic", reanimator.Tran)
	a.SetEndTrigger(func() {
		first = false
	})
	b := reanimator.NewBatchAnimation("portalClose", img.Batchers[constants.TileBatch], "portal_magic_close", reanimator.Tran)
	b.SetEndTrigger(func() {
		myecs.Manager.DisposeEntity(m)
	})
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(a).
		AddAnimation(b).
		SetChooseFn(func() string {
			if first {
				return "portal"
			} else {
				return "portalClose"
			}
		}), "portal")
	m.AddComponent(myecs.Object, obj)
	m.AddComponent(myecs.Animated, anim)
	m.AddComponent(myecs.Drawable, anim)
	m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
}
