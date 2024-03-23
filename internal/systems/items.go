package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
	"strings"
)

func CreateGem(pos pixel.Vec, key string) {
	obj := object.New().WithID("gem")
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 10, 10))
	obj.Layer = 11
	gemShimmer := false
	batch := img.Batchers[constants.TileBatch]
	color := strings.Replace(key, "gem_", "", -1)
	gemSpr := reanimator.NewBatchSprite("gem", batch, key, reanimator.Hold)
	shimmer := reanimator.NewBatchAnimation("shimmer", batch, fmt.Sprintf("gem_%s_shimmer", color), reanimator.Tran)
	shimmer.SetTrigger(3, func() {
		gemShimmer = false
	})
	anim := reanimator.New(reanimator.NewSwitch().
		AddAnimation(gemSpr).
		AddAnimation(shimmer).
		SetChooseFn(func() string {
			if gemShimmer {
				return "shimmer"
			} else {
				if !gemShimmer && random.Effects.Intn(constants.IdleFrequency*timing.FPS) == 0 {
					gemShimmer = true
				}
				return "gem"
			}
		}), "gem")
	myecs.Manager.NewEntity().
		AddComponent(myecs.Object, obj).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Animated, anim).
		AddComponent(myecs.Gem, struct{}{}).
		AddComponent(myecs.OnTouch, data.NewInteract(CollectGem)).
		AddComponent(myecs.LvlElement, struct{}{})
}

func CollectGem(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	level.Stats[p].Score += 1
	myecs.Manager.DisposeEntity(entity)
}

func CreateDoor(pos pixel.Vec, key string) {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 6, 6))
	obj.Layer = 9
	door := &data.Door{
		Object: obj,
	}
	switch key {
	case constants.TileDoorYellow, constants.TileDoorOrange,
		constants.TileDoorGray, constants.TileDoorCyan,
		constants.TileDoorBlue, constants.TileDoorGreen,
		constants.TileDoorPurple, constants.TileDoorBrown:
		door.DoorType = data.Opened
		door.Color = strings.Replace(key, "door_", "", -1)
	case constants.TileLockYellow, constants.TileLockOrange,
		constants.TileLockGray, constants.TileLockCyan,
		constants.TileLockBlue, constants.TileLockGreen,
		constants.TileLockPurple, constants.TileLockBrown:
		door.DoorType = data.Locked
		door.Color = strings.Replace(key, "lock_", "", -1)
	}
	var interaction *data.Interact
	switch door.Color {
	case constants.StrColorBlue,
		constants.StrColorGreen,
		constants.StrColorPurple,
		constants.StrColorBrown:
		interaction = EnterPlayerDoor(door.Color)
	default:
		interaction = EnterDoor()
	}
	e := myecs.Manager.NewEntity()
	door.Entity = e
	var anim *reanimator.Anim
	if door.DoorType == data.Locked {
		anim = reanimator.NewBatchAnimation("unlock", img.Batchers[constants.TileBatch], fmt.Sprintf("unlock_%s_open", door.Color), reanimator.Hold)
		anim.SetEndTrigger(func() {
			e.RemoveComponent(myecs.Animated)
			e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("unlock_%s", door.Color), constants.TileBatch))
			e.AddComponent(myecs.OnTouch, interaction)
		})
	} else {
		anim = reanimator.NewBatchAnimation("open", img.Batchers[constants.TileBatch], fmt.Sprintf("door_%s_open", door.Color), reanimator.Hold)
		anim.SetEndTrigger(func() {
			e.RemoveComponent(myecs.Animated)
			e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
			e.AddComponent(myecs.OnTouch, interaction)
		})
	}
	e.AddComponent(myecs.Object, obj)
	if door.DoorType == data.Locked {
		e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	}
	e.AddComponent(myecs.Door, door)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
		switch door.DoorType {
		case data.Opened:
			if data.CurrLevel.DoorsOpen {
				tree := reanimator.NewSimple(anim)
				e.AddComponent(myecs.Drawable, tree)
				e.AddComponent(myecs.Animated, tree)
				e.RemoveComponent(myecs.Update)
			}
		case data.Locked:
			if door.Unlock {
				e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("unlock_%s_open", door.Color), constants.TileBatch))
				door.DoorType = data.Unlocked
			}
		case data.Unlocked:
			if data.CurrLevel.DoorsOpen {
				tree := reanimator.NewSimple(anim)
				e.AddComponent(myecs.Drawable, tree)
				e.AddComponent(myecs.Animated, tree)
				e.RemoveComponent(myecs.Update)
			}
		}
		return false
	}))
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func EnterDoor() *data.Interact {
	return data.NewInteract(func(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers {
			return
		}
		level.Stats[p].Score += 12
		level.Complete = true
	})
}

func EnterPlayerDoor(color string) *data.Interact {
	return data.NewInteract(func(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers || ch.Color != color {
			return
		}
		level.Stats[p].Score += 12
		level.Complete = true
	})
}

func CreateBox(pos pixel.Vec) {
	key := constants.ItemBox
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	s := obj.Pos.Y
	smash := &s
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	e.AddComponent(myecs.StandOn, struct{}{})
	e.AddComponent(myecs.Smash, smash)
	e.AddComponent(myecs.OnTouch, data.NewInteract(BoxBonk))
	e.AddComponent(myecs.PickUp, data.NewPickUp(10, true))
	e.AddComponent(myecs.LvlElement, struct{}{})
	box := data.NewDynamic()
	box.Object = obj
	box.Entity = e
	box.Flags.NoLadders = true
	e.AddComponent(myecs.Dynamic, box)
}

func BoxBonk(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
	s, ok := entity.GetComponentData(myecs.Smash)
	d, okD := entity.GetComponentData(myecs.Dynamic)
	if ok && okD {
		f := s.(*float64)
		box := d.(*data.Dynamic)
		if *f-box.Object.Pos.Y >= constants.SmashDistance &&
			ch.Object.Pos.Y < box.Object.Pos.Y &&
			(ch.State != data.Falling || ch.Vars.Gravity < box.Vars.Gravity) {
			ch.Flags.Hit = true
			ch.State = data.Hit
			//for i := 0; i < 5; i++ {
			//	pObj := object.New()
			//	pObj.Pos = box.Object.Pos
			//	pObj.Layer = 12
			//	switch i {
			//	case 0:
			//		pObj.Pos = pObj.Pos.Add(pixel.V(-5, 5))
			//	case 1:
			//		pObj.Pos = pObj.Pos.Add(pixel.V(-7, -3))
			//	case 2:
			//		pObj.Pos = pObj.Pos.Add(pixel.V(6, 6))
			//	case 3:
			//		pObj.Pos = pObj.Pos.Add(pixel.V(5, 0))
			//	case 4:
			//		pObj.Pos = pObj.Pos.Add(pixel.V(4, -5))
			//	}
			//	end := 7
			//	count := 0
			//	e := myecs.Manager.NewEntity()
			//	e.AddComponent(myecs.Object, pObj)
			//	e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("%s%d", constants.ItemBoxPiece, i), constants.TileBatch))
			//	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
			//	e.AddComponent(myecs.Update, data.NewFn(func() {
			//		if reanimator.FrameSwitch {
			//			count++
			//			if count >= end {
			//				myecs.Manager.DisposeEntity(e)
			//			}
			//		}
			//	}))
			//}
			//myecs.Manager.DisposeEntity(entity)
		}
	}
}

func CreateKey(pos pixel.Vec, key string) {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 14))
	obj.Layer = 14
	color := strings.Replace(key, "key_", "", -1)
	theKey := &data.Key{
		Object: obj,
		Sprite: img.NewSprite(key, constants.TileBatch),
		PickUp: data.NewPickUp(5, false),
		Action: KeyAction(color),
	}
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, theKey.Sprite)
	e.AddComponent(myecs.PickUp, theKey.PickUp)
	e.AddComponent(myecs.Action, theKey.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func KeyAction(color string) *data.Interact {
	return data.NewInteract(func(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
		if KeyUnlock(level, ch.Object.Pos, color) {
			DropLift(ch, false)
			myecs.Manager.DisposeEntity(entity)
		}
		ch.Flags.Using = false
	})
}

func KeyUnlock(level *data.Level, chPos pixel.Vec, color string) bool {
	chX, chY := world.WorldToMap(chPos.X, chPos.Y)
	for _, result := range myecs.Manager.Query(myecs.IsDoor) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		d, okD := result.Components[myecs.Door].(*data.Door)
		if okO && okD {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			if x == chX && y == chY &&
				!d.Unlock &&
				d.DoorType == data.Locked &&
				d.Color == color {
				d.Unlock = true
				return true
			}
		}
	}
	return false
}
