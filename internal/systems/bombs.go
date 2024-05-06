package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/google/uuid"
	"github.com/gopxl/pixel"
	"math"
)

func CreateBomb(pos pixel.Vec, key string, metadata data.TileMetadata, origin world.Coords) {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 14, 16))
	obj.Layer = 14
	e := myecs.Manager.NewEntity()
	theBomb := &data.Bomb{
		Object:   obj,
		Entity:   e,
		Metadata: metadata,
		Origin:   origin,
	}
	regenA := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "bomb_regen_anim", []int{0, 1, 2, 3, 3, 4, 5, 6, 6, 6, 6}, reanimator.Tran)
	regenA.SetTriggerAll(func() {
		switch regenA.Step {
		case 7, 8:
			obj.Pos.Y = pos.Y + 2
			theBomb.SymSpr.ToggleHidden(false)
		case 9:
			obj.Pos.Y = pos.Y
			theBomb.SymSpr.ToggleHidden(false)
		case 10:
			obj.Pos.Y = pos.Y + 1
			theBomb.SymSpr.ToggleHidden(false)
		default:
			theBomb.SymSpr.ToggleHidden(true)
		}
	})
	regenA.SetEndTrigger(func() {
		obj.Pos.Y = pos.Y
		theBomb.Regen = false
		theBomb.SymSpr.ToggleHidden(false)
	})
	theBomb.Anim = reanimator.New(reanimator.NewSwitch().
		AddAnimation(regenA).
		AddAnimation(reanimator.NewBatchSprite("item", img.Batchers[constants.TileBatch], "bomb", reanimator.Hold)).
		AddNull("none").
		SetChooseFn(func() string {
			if theBomb.Waiting {
				return "none"
			} else if theBomb.Regen {
				return "regen"
			} else {
				return "item"
			}
		}), "fuse1")
	theBomb.Draws = append(theBomb.Draws, theBomb.Anim)
	name := "Bomb"
	litKey := constants.ItemBombLit
	if metadata.BombCross && metadata.Regenerate {
		theBomb.SymSpr = img.NewSprite(constants.ItemBombRegenCross, constants.TileBatch).WithOffset(pixel.V(0, -2))
		name = "Bomb Cross"
		litKey = constants.ItemBombLitCross
	} else if metadata.BombCross {
		theBomb.SymSpr = img.NewSprite(constants.ItemBombCross, constants.TileBatch).WithOffset(pixel.V(0, -2))
		name = "Bomb Cross"
		litKey = constants.ItemBombLitCross
	} else if metadata.Regenerate {
		theBomb.SymSpr = img.NewSprite(constants.ItemBombRegen, constants.TileBatch).WithOffset(pixel.V(0, -2))
	}
	if theBomb.SymSpr != nil {
		theBomb.Draws = append(theBomb.Draws, theBomb.SymSpr)
	}
	theBomb.Name = name
	theBomb.LitKey = litKey
	theBomb.PickUp = data.NewPickUp(name, 5)
	theBomb.Action = BombAction(theBomb)

	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, theBomb.Draws)
	e.AddComponent(myecs.Animated, theBomb.Anim)
	e.AddComponent(myecs.PickUp, theBomb.PickUp)
	e.AddComponent(myecs.Action, theBomb.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
	e.AddComponent(myecs.Bomb, theBomb)
}

func BombAction(theBomb *data.Bomb) *data.Interact {
	return data.NewInteract(func(lvl *data.Level, p int, ch *data.Dynamic, e *ecs.Entity) {
		switch ch.State {
		case data.OnBar, data.Jumping, data.Falling:
			return
		}
		DropItem(ch)
		x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
		tile := lvl.Tiles.Get(x, y)
		LightBomb(theBomb, tile)
	})
}

func LightBomb(theBomb *data.Bomb, tile *data.Tile) {
	if theBomb.Metadata.Regenerate {
		theBomb.SymSpr.ToggleHidden(true)
		counter := 0
		theBomb.Waiting = true
		theBomb.Entity.AddComponent(myecs.Update, data.NewFn(func() {
			if reanimator.FrameSwitch {
				counter++
			}
			delay := constants.BombFuse + (constants.BombRegen * (theBomb.Metadata.RegenDelay + 2))
			if counter > delay && data.CurrLevel.FrameChange {
				theBomb.Object.Pos = world.MapToWorld(theBomb.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
				theBomb.Regen = true
				theBomb.Waiting = false
				theBomb.Entity.RemoveComponent(myecs.Update)
			}
		}))
	} else {
		myecs.Manager.DisposeEntity(theBomb.Entity)
	}
	CreateLitBomb(tile.Object.Pos, theBomb.LitKey, data.TileMetadata{BombCross: theBomb.Metadata.BombCross})
}

func CreateLitBomb(pos pixel.Vec, key string, metadata data.TileMetadata) {
	counter := 0
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 14, 16))
	obj.Layer = 30
	waiting := false
	regen := false
	var symbolSprite *img.Sprite
	if metadata.BombCross && metadata.Regenerate {
		symbolSprite = img.NewSprite(constants.ItemBombRegenCross, constants.TileBatch).WithOffset(pixel.V(0, -2))
	} else if metadata.BombCross {
		symbolSprite = img.NewSprite(constants.ItemBombCross, constants.TileBatch).WithOffset(pixel.V(0, -2))
	} else if metadata.Regenerate {
		symbolSprite = img.NewSprite(constants.ItemBombRegen, constants.TileBatch).WithOffset(pixel.V(0, -2))
	}
	//fuseSound := sfx.SoundPlayer.PlaySound(constants.SFXBombLight, 0)
	e := myecs.Manager.NewEntity()
	fuse1 := reanimator.NewBatchAnimation("fuse1", img.Batchers[constants.TileBatch], "bomb_fuse1", reanimator.Loop)
	fuse1.SetTriggerAll(func() {
		counter++
	})
	fuse2 := reanimator.NewBatchAnimation("fuse2", img.Batchers[constants.TileBatch], "bomb_fuse2", reanimator.Tran)
	fuse2.SetEndTrigger(func() {
		if metadata.Regenerate {
			obj.Hidden = true
			waiting = true
			counter = 0
			e.AddComponent(myecs.Update, data.NewFn(func() {
				if reanimator.FrameSwitch {
					counter++
				}
				if counter > constants.BombRegen*(metadata.RegenDelay+2) && data.CurrLevel.FrameChange {
					obj.Hidden = false
					counter = 0
					waiting = false
					regen = true
					e.RemoveComponent(myecs.Update)
				}
			}))
		} else {
			myecs.Manager.DisposeEntity(e)
		}
		CreateExplosion(pos, metadata.BombCross, nil)
	})
	regenA := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "bomb_regen_anim", []int{0, 1, 2, 3, 3, 4, 5, 6, 6, 6, 6}, reanimator.Tran)
	regenA.SetTriggerAll(func() {
		switch regenA.Step {
		case 7, 8:
			obj.Pos.Y = pos.Y + 2
			symbolSprite.ToggleHidden(false)
		case 9:
			obj.Pos.Y = pos.Y
			symbolSprite.ToggleHidden(false)
		case 10:
			obj.Pos.Y = pos.Y + 1
			symbolSprite.ToggleHidden(false)
		default:
			symbolSprite.ToggleHidden(true)
		}
	})
	regenA.SetEndTrigger(func() {
		obj.Pos.Y = pos.Y
		regen = false
		symbolSprite.ToggleHidden(false)
	})
	tree := reanimator.New(reanimator.NewSwitch().
		AddAnimation(fuse1).
		AddAnimation(fuse2).
		AddAnimation(regenA).
		AddNull("none").
		SetChooseFn(func() string {
			if waiting {
				return "none"
			} else if regen {
				return "regen"
			} else if data.CurrLevel.Start && counter > constants.BombFuse {
				return "fuse2"
			} else {
				return "fuse1"
			}
		}), "fuse1")
	draws := []interface{}{tree}
	if symbolSprite != nil {
		draws = append(draws, symbolSprite)
	}
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, draws)
	e.AddComponent(myecs.Animated, tree)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func CreateExplosion(pos pixel.Vec, cross bool, fuseSound *uuid.UUID) {
	var coords []world.Coords
	x, y := world.WorldToMap(pos.X, pos.Y)
	if cross {
		for n := y - 3; n < y+4; n++ {
			if n == y {
				for m := x - 3; m < x+4; m++ {
					coords = append(coords, world.NewCoords(m, n))
				}
			} else {
				coords = append(coords, world.NewCoords(x, n))
			}
		}
	} else {
		for n := y - 2; n < y+3; n++ {
			if n == y {
				for m := x - 2; m < x+3; m++ {
					coords = append(coords, world.NewCoords(m, n))
				}
			} else if n == y-1 || n == y+1 {
				for m := x - 1; m < x+2; m++ {
					coords = append(coords, world.NewCoords(m, n))
				}
			} else {
				coords = append(coords, world.NewCoords(x, n))
			}
		}
	}
	for _, c := range coords {
		key := "exp_end"
		r := 0.
		flip := false
		flop := false
		if cross {
			if c.X == x && c.Y == y {
				key = "exp_cross"
			} else if (c.X == x && util.Abs(c.Y-y) < 3) ||
				(c.Y == y && util.Abs(c.X-x) < 3) {
				key = "exp_line"
			}
			if c.X == x && c.Y != y { // vert line
				if c.Y < y { // bottom of cross
					flop = true
				}
			} else if c.Y == y && c.X != x { // horiz line
				r = math.Pi * -0.5
				if c.X < x { // left of cross
					r = math.Pi * 0.5
				}
			}
		} else {
			if c.X == x && c.Y == y {
				key = "exp_center"
			} else if (c.X == x && util.Abs(c.Y-y) == 1) ||
				(c.Y == y && util.Abs(c.X-x) == 1) {
				key = "exp_tee"
			} else if util.Abs(c.X-x) == 1 && util.Abs(c.Y-y) == 1 {
				key = "exp_corner"
			}
			if c.X == x && c.Y != y { // vert line
				if c.Y < y { // bottom
					flop = true
				}
			} else if c.Y == y && c.X != x { // horiz line
				r = math.Pi * -0.5
				if c.X < x {
					r = math.Pi * 0.5
				}
			} else if c.X > x && c.Y > y { // top right
				flip = true
			} else if c.X > x && c.Y < y { // bottom right
				flip = true
				flop = true
			} else if c.X < x && c.Y < y { // bottom left
				flop = true
			}
		}
		// destroy turf
		tile := data.CurrLevel.Tiles.Get(c.X, c.Y)
		if tile != nil {
			switch tile.Block {
			case data.BlockTurf, data.BlockCracked, data.BlockFall:
				tile.Flags.Collapse = true
				tile.Counter = constants.CollapseCounter
			}
		}
		// explosion
		obj := object.New()
		obj.Layer = 34
		obj.Pos = world.MapToWorld(c).Add(pixel.V(world.HalfSize, world.HalfSize))
		obj.SetRect(pixel.R(0, 0, 16, 16))
		obj.Flip = flip
		obj.Flop = flop
		obj.Rot = r
		e := myecs.Manager.NewEntity()
		a := reanimator.NewBatchAnimation(key, img.Batchers[constants.TileBatch], key, reanimator.Done)
		a.SetEndTrigger(func() {
			myecs.Manager.DisposeEntity(e)
		})
		anim := reanimator.NewSimple(a)
		e.AddComponent(myecs.Object, obj)
		e.AddComponent(myecs.Drawable, anim)
		e.AddComponent(myecs.Animated, anim)
		e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	}
	// kill players and enemies
	for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
		_, okCO := resultC.Components[myecs.Object].(*object.Object)
		ch, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
		if okCO && okC && ch.State != data.Dead {
			chX, chY := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			tile := data.CurrLevel.Tiles.Get(chX, chY)
			if world.CoordsIn(world.NewCoords(chX, chY), coords) {
				ch.Object.Pos.X = tile.Object.Pos.X
				ch.Object.Pos.Y = tile.Object.Pos.Y
				ch.Flags.Blow = true
				ch.Flags.Hit = true
				ch.State = data.Hit
				ch.Object.Layer = 35
			}
		}
	}
	// Light other bombs
	for _, resultB := range myecs.Manager.Query(myecs.IsBomb) {
		obj, okO := resultB.Components[myecs.Object].(*object.Object)
		b, okB := resultB.Components[myecs.Bomb].(*data.Bomb)
		if okO && okB && !obj.Hidden && !b.Waiting && !b.Regen {
			chX, chY := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			tile := data.CurrLevel.Tiles.Get(chX, chY)
			if world.CoordsIn(world.NewCoords(chX, chY), coords) {
				obj.Pos.X = tile.Object.Pos.X
				obj.Pos.Y = tile.Object.Pos.Y
				LightBomb(b, tile)
			}
		}
	}
	// sfx
	sfx.SoundPlayer.KillSound(fuseSound)
	sfx.SoundPlayer.PlaySound(constants.SFXBombBlow, 0)
}
