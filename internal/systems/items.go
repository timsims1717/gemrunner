package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
	"math/rand"
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
				if !gemShimmer && rand.Intn(constants.IdleFrequency*timing.FPS) == 0 {
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
		AddComponent(myecs.OnTouch, data.NewInteract(CollectGem))
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
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	e.AddComponent(myecs.StandOn, struct{}{})
	e.AddComponent(myecs.PickUp, data.NewPickUp(10, true))
	box := data.NewDynamic()
	box.Object = obj
	box.Entity = e
	e.AddComponent(myecs.Dynamic, box)
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
}

func KeyAction(color string) *data.Interact {
	return data.NewInteract(func(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
		if o, okO := entity.GetComponentData(myecs.Object); okO {
			obj := o.(*object.Object)
			if KeyUnlock(level, obj.Pos.Add(obj.Offset), ch.Object.Pos, color) {
				myecs.Manager.DisposeEntity(entity)
			}
		}
	})
}

func KeyUnlock(level *data.Level, pos1, pos2 pixel.Vec, color string) bool {
	x1, y1 := world.WorldToMap(pos1.X, pos1.Y)
	x2, y2 := world.WorldToMap(pos2.X, pos2.Y)
	for _, result := range myecs.Manager.Query(myecs.IsDoor) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		d, okD := result.Components[myecs.Door].(*data.Door)
		if okO && okD {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			if ((x == x1 && y == y1) || (x == x2 && y == y2)) &&
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
