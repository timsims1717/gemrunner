package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func CreateGem(pos pixel.Vec) {
	obj := object.New()
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 6, 6))
	obj.Layer = 10
	myecs.Manager.NewEntity().
		AddComponent(myecs.Object, obj).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Drawable, img.NewSprite(constants.ItemGem, constants.TileBatch)).
		AddComponent(myecs.Gem, struct{}{}).
		AddComponent(myecs.Interact, data.NewOnTouch(CollectGem))
}

func CollectGem(level *data.Level, ch *data.Character, entity *ecs.Entity) {
	switch ch.PlayerIndex {
	default:
		return
	case 0:
		level.Stats1.Score += 1
	}
	myecs.Manager.DisposeEntity(entity)
}

func CreateEmptyDoor(pos pixel.Vec, key string) {
	obj := object.New()
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 6, 6))
	obj.Layer = 10
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
		noGems := len(myecs.Manager.Query(myecs.IsGem)) < 1
		if noGems {
			e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
			e.AddComponent(myecs.Interact, data.NewOnTouch(EnterDoor))
			e.RemoveComponent(myecs.Update)
		}
		return false
	}))
}

func EnterDoor(level *data.Level, ch *data.Character, entity *ecs.Entity) {
	switch ch.PlayerIndex {
	default:
		return
	case 0:
		level.Stats1.Score += 12
	}
	level.Complete = true
}
