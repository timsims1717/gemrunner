package myecs

import (
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
)

var (
	FullCount   = 0
	IDCount     = 0
	LoadedCount = 0
)

type ClearFlag bool

var (
	Manager = ecs.NewManager()

	Object = Manager.NewComponent()
	Parent = Manager.NewComponent()
	Temp   = Manager.NewComponent()
	Update = Manager.NewComponent()

	Drawable   = Manager.NewComponent()
	Animated   = Manager.NewComponent()
	DrawTarget = Manager.NewComponent()

	Tile   = Manager.NewComponent()
	Border = Manager.NewComponent()
	Block  = Manager.NewComponent()

	IsTemp    = ecs.BuildTag(Temp, Object)
	HasUpdate = ecs.BuildTag(Update)

	HasAnimation = ecs.BuildTag(Animated, Object)
	IsDrawable   = ecs.BuildTag(Drawable, Object)

	IsObject  = ecs.BuildTag(Object)
	HasParent = ecs.BuildTag(Object, Parent)

	IsTile    = ecs.BuildTag(Object, Tile)
	HasBorder = ecs.BuildTag(Object, Border)
	IsBlock   = ecs.BuildTag(Object, Block)
)

func UpdateManager() {
	LoadedCount = 0
	IDCount = 0
	FullCount = 0
	for _, result := range Manager.Query(IsObject) {
		if t, ok := result.Components[Object].(*object.Object); ok {
			FullCount++
			if t.ID != "" {
				IDCount++
				if t.Loaded {
					LoadedCount++
				}
			}
		}
	}
}
