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

	Object        = Manager.NewComponent()
	Parent        = Manager.NewComponent()
	Temp          = Manager.NewComponent()
	Update        = Manager.NewComponent()
	Interpolation = Manager.NewComponent()

	Drawable   = Manager.NewComponent()
	Animated   = Manager.NewComponent()
	DrawTarget = Manager.NewComponent()

	Tile       = Manager.NewComponent()
	Border     = Manager.NewComponent()
	Block      = Manager.NewComponent()
	Dynamic    = Manager.NewComponent()
	Controller = Manager.NewComponent()
	Player     = Manager.NewComponent()
	OnTouch    = Manager.NewComponent()
	Gem        = Manager.NewComponent()
	Door       = Manager.NewComponent()
	StandOn    = Manager.NewComponent()
	Smash      = Manager.NewComponent()
	PickUp     = Manager.NewComponent()
	Action     = Manager.NewComponent()
	LvlElement = Manager.NewComponent()
	Bomb       = Manager.NewComponent()
	Text       = Manager.NewComponent()
	Item       = Manager.NewComponent()

	IsTemp    = ecs.BuildTag(Temp)
	HasUpdate = ecs.BuildTag(Update)

	HasAnimation = ecs.BuildTag(Animated, Object)
	IsDrawable   = ecs.BuildTag(Drawable, Object)

	IsObject         = ecs.BuildTag(Object)
	HasParent        = ecs.BuildTag(Object, Parent)
	HasInterpolation = ecs.BuildTag(Object, Interpolation)

	IsTile       = ecs.BuildTag(Object, Tile)
	HasBorder    = ecs.BuildTag(Object, Border)
	IsBlock      = ecs.BuildTag(Object, Block)
	IsDynamic    = ecs.BuildTag(Object, Dynamic)
	IsPlayer     = ecs.BuildTag(Object, Dynamic, Player, Controller)
	IsCharacter  = ecs.BuildTag(Object, Dynamic, Controller)
	IsTouch      = ecs.BuildTag(Object, OnTouch)
	IsGem        = ecs.BuildTag(Object, Gem)
	IsStandOn    = ecs.BuildTag(Object, Dynamic, StandOn)
	IsSmash      = ecs.BuildTag(Object, Dynamic, Smash)
	IsPickUp     = ecs.BuildTag(Object, PickUp, Item)
	HasAction    = ecs.BuildTag(Object, PickUp, Action)
	IsDoor       = ecs.BuildTag(Object, Door)
	IsLvlElement = ecs.BuildTag(Object, LvlElement)
	IsBomb       = ecs.BuildTag(Object, Bomb)
	IsText       = ecs.BuildTag(Object, Drawable, Text)
	IsItem       = ecs.BuildTag(Object, Item)
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
				if !t.Unloaded {
					LoadedCount++
				}
			}
		}
	}
}
