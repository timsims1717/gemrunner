package systems

import (
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
)

func CreateGem(pos pixel.Vec) {
	obj := object.New().WithID("gem")
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 10, 10))
	obj.Layer = 11
	gemShimmer := false
	batch := img.Batchers[constants.TileBatch]
	gemSpr := reanimator.NewBatchSprite("gem", batch, constants.ItemGem, reanimator.Hold)
	shimmer := reanimator.NewBatchAnimation("shimmer", batch, "gem_shimmer", reanimator.Tran)
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
	obj.Layer = 10
	door := &data.Door{
		Object: obj,
	}
	switch key {
	case constants.TileDoorPink:
		door.DoorType = data.PinkOpen
	case constants.TileDoorBlue:
		door.DoorType = data.BlueOpen
	case constants.TileLockPink:
		door.DoorType = data.PinkLock
	case constants.TileLockBlue:
		door.DoorType = data.BlueLock
	}
	e := myecs.Manager.NewEntity()
	door.Entity = e
	e.AddComponent(myecs.Object, obj)
	if key == constants.TileLockBlue || key == constants.TileLockPink {
		e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	}
	e.AddComponent(myecs.Door, door)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
		switch door.DoorType {
		case data.PinkOpen, data.BlueOpen:
			noGems := len(myecs.Manager.Query(myecs.IsGem)) < 1
			if noGems {
				e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
				e.AddComponent(myecs.OnTouch, data.NewInteract(EnterDoor))
				e.RemoveComponent(myecs.Update)
			}
		case data.PinkLock, data.BlueLock:
			if door.Unlock {
				if door.DoorType%data.KeyParity == 0 {
					e.AddComponent(myecs.Drawable, img.NewSprite(constants.TileUnlockPink, constants.TileBatch))
					door.DoorType = data.PinkUnlock
				} else {
					e.AddComponent(myecs.Drawable, img.NewSprite(constants.TileUnlockBlue, constants.TileBatch))
					door.DoorType = data.BlueUnlock
				}
			}
		case data.PinkUnlock, data.BlueUnlock:
			noGems := len(myecs.Manager.Query(myecs.IsGem)) < 1
			if noGems {
				if door.DoorType%data.KeyParity == 0 {
					e.AddComponent(myecs.Drawable, img.NewSprite(constants.TileDoorPink, constants.TileBatch))
					door.DoorType = data.PinkOpen
				} else {
					e.AddComponent(myecs.Drawable, img.NewSprite(constants.TileDoorBlue, constants.TileBatch))
					door.DoorType = data.BlueOpen
				}
				e.AddComponent(myecs.OnTouch, data.NewInteract(EnterDoor))
				e.RemoveComponent(myecs.Update)
			}
		}
		return false
	}))
}

func EnterDoor(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	level.Stats[p].Score += 12
	level.Complete = true
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
	e.AddComponent(myecs.PickUp, data.NewPickUp(constants.PickUpPriority[key], true))
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
	var kt data.KeyType
	switch key {
	case constants.ItemKeyPink:
		kt = data.PinkKey
	case constants.ItemKeyBlue:
		kt = data.BlueKey
	}
	theKey := &data.Key{
		Object:  obj,
		Sprite:  img.NewSprite(key, constants.TileBatch),
		PickUp:  data.NewPickUp(constants.PickUpPriority[key], false),
		Action:  KeyAction(kt),
		KeyType: kt,
	}
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, theKey.Sprite)
	e.AddComponent(myecs.PickUp, theKey.PickUp)
	e.AddComponent(myecs.Action, theKey.Action)
}

func KeyAction(keyType data.KeyType) *data.Interact {
	return data.NewInteract(func(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
		if o, okO := entity.GetComponentData(myecs.Object); okO {
			obj := o.(*object.Object)
			if KeyUnlock(level, obj.Pos.Add(obj.Offset), ch.Object.Pos, keyType) {
				myecs.Manager.DisposeEntity(entity)
			}
		}
	})
}

func KeyUnlock(level *data.Level, pos1, pos2 pixel.Vec, keyType data.KeyType) bool {
	x1, y1 := world.WorldToMap(pos1.X, pos1.Y)
	x2, y2 := world.WorldToMap(pos2.X, pos2.Y)
	for _, result := range myecs.Manager.Query(myecs.IsDoor) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		d, okD := result.Components[myecs.Door].(*data.Door)
		if okO && okD {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			if ((x == x1 && y == y1) || (x == x2 && y == y2)) &&
				!d.Unlock && (d.DoorType == data.PinkLock || d.DoorType == data.BlueLock) &&
				int(d.DoorType%data.KeyParity) == int(keyType%data.KeyParity) {
				d.Unlock = true
				return true
			}
		}
	}
	return false
}
