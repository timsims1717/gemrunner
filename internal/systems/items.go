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
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func CreateGem(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 10, 10))
	obj.Layer = 11
	gemShimmer := false
	batch := img.Batchers[constants.TileBatch]
	color := tile.Metadata.Color
	gemSpr := reanimator.NewBatchSprite("gem", batch, key, reanimator.Hold)
	shimmer := reanimator.NewBatchAnimationCustom("shimmer", batch, fmt.Sprintf("%s_shimmer", key), []int{0, 0, 1, 1, 2, 2}, reanimator.Tran)
	shimmer.SetEndTrigger(func() {
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
		AddComponent(myecs.OnTouch, CollectGem(color)).
		AddComponent(myecs.LvlElement, struct{}{})
}

func CollectGem(color data.ItemColor) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers || (ch.Color != color && color > data.NonPlayerRed) {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 1
		data.CurrLevelSess.PlayerStats[p].LGems++
		sfx.SoundPlayer.PlaySound(constants.SFXGem, -2.)
		myecs.Manager.DisposeEntity(entity)
	})
}

func CreateDoor(pos pixel.Vec, tile *data.Tile) {
	obj := object.New().WithID(tile.SpriteString())
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 8))
	obj.Layer = 9
	door := &data.Door{}
	door.Object = obj
	door.Color = tile.Metadata.Color
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	switch tile.Block {
	case data.BlockDoorHidden:
		door.DoorType = data.Hidden
	case data.BlockDoorVisible:
		door.DoorType = data.Hidden
		e.AddComponent(myecs.Drawable, img.NewSprite(tile.SpriteString(), constants.TileBatch))
	case data.BlockDoorLocked:
		door.DoorType = data.Locked
		e.AddComponent(myecs.Drawable, img.NewSprite(tile.SpriteString(), constants.TileBatch))
	}
	var interaction *data.Interact
	switch door.Color {
	case data.PlayerBlue, data.PlayerGreen, data.PlayerPurple, data.PlayerOrange:
		interaction = EnterPlayerDoor(door.Color, tile.Metadata.ExitIndex)
	default:
		interaction = EnterDoor(tile.Metadata.ExitIndex)
	}
	door.Entity = e
	var anim *reanimator.Anim
	if door.DoorType == data.Locked {
		anim = reanimator.NewBatchAnimation("unlock", img.Batchers[constants.TileBatch], fmt.Sprintf("door_unlock%s", door.Color.SpriteString()), reanimator.Tran)
		anim.SetEndTrigger(func() {
			e.RemoveComponent(myecs.Animated)
			e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("door_unlocked%s", door.Color.SpriteString()), constants.TileBatch))
			e.AddComponent(myecs.OnTouch, interaction)
		})
	} else {
		anim = reanimator.NewBatchAnimation("open", img.Batchers[constants.TileBatch], fmt.Sprintf("door_visible%s", door.Color.SpriteString()), reanimator.Tran)
		anim.SetEndTrigger(func() {
			e.RemoveComponent(myecs.Animated)
			e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("door_hidden%s", door.Color.SpriteString()), constants.TileBatch))
			e.AddComponent(myecs.OnTouch, interaction)
		})
	}
	e.AddComponent(myecs.Door, door)
	e.AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
		switch door.DoorType {
		case data.Hidden:
			if data.CurrLevel.DoorsOpen && reanimator.FrameSwitch {
				tree := reanimator.NewSimple(anim)
				e.AddComponent(myecs.Drawable, tree)
				e.AddComponent(myecs.Animated, tree)
				e.RemoveComponent(myecs.Update)
			}
		case data.Locked:
			if door.Unlock {
				e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("door_unlock%s", door.Color.SpriteString()), constants.TileBatch))
				door.DoorType = data.Unlocked
			}
		case data.Unlocked:
			if data.CurrLevel.DoorsOpen && reanimator.FrameSwitch {
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

func EnterDoor(exitIndex int) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 12
		sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
		data.CurrLevel.Complete = true
		data.CurrLevel.ExitIndex = exitIndex
	})
}

func EnterPlayerDoor(color data.ItemColor, exitIndex int) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers || ch.Color != color {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 12
		sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
		data.CurrLevel.Complete = true
		data.CurrLevel.ExitIndex = exitIndex
	})
}

func CreateJumpBoots(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	jumpBoots := &data.BasicItem{
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp("Drill", 5, tile.Metadata.Color),
		Entity:   e,
		Metadata: tile.Metadata,
		Origin:   tile.Coords,
	}
	jumpBoots.Action = Jump()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, jumpBoots.Sprite)
	e.AddComponent(myecs.PickUp, jumpBoots.PickUp)
	e.AddComponent(myecs.Action, jumpBoots.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func Jump() *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Grounded {
			x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			left := data.CurrLevel.Tiles.Get(x-1, y)
			right := data.CurrLevel.Tiles.Get(x+1, y)
			// High Jump if:
			//  there is no ceiling here
			//  the character is not going left or right
			//  or they are going left/right and there is a wall left/right
			//  or they are going left/right and there is a wall up left or up right
			// Otherwise, it's a long jump
			if !ch.Flags.Ceiling &&
				((!ch.Actions.Left() && !ch.Actions.Right()) ||
					(ch.Actions.Left() && left.IsSolid()) ||
					(ch.Actions.Right() && right.IsSolid())) {
				ch.Flags.HighJump = true
			} else if (ch.Actions.Left() && !left.IsSolid()) ||
				(ch.Actions.Right() && !right.IsSolid()) {
				ch.Flags.LongJump = true
			} else {
				return
			}
			tile := data.CurrLevel.Tiles.Get(x, y)
			ch.LastTile = tile
			ch.State = data.Jumping
			ch.Object.Pos.X = tile.Object.Pos.X
			ch.Object.Pos.Y = tile.Object.Pos.Y
			// for both kinds of jumps
			ch.ACounter = 0
			if ch.Actions.Left() {
				ch.Flags.JumpL = true
				ch.Object.Flip = true
			} else if ch.Actions.Right() {
				ch.Flags.JumpR = true
				ch.Object.Flip = false
			} else {
				ch.Flags.JumpL = false
				ch.Flags.JumpR = false
			}
		}
	})
}

func CreateBox(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
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
	e.AddComponent(myecs.PickUp, data.NewPickUp("Box", 10, tile.Metadata.Color))
	e.AddComponent(myecs.Action, data.NewInteract(BoxAction))
	e.AddComponent(myecs.LvlElement, struct{}{})
	box := data.NewDynamic(tile)
	box.Object = obj
	box.Entity = e
	box.Flags.NoLadders = true
	e.AddComponent(myecs.Dynamic, box)
}

func BoxAction(p int, ch *data.Dynamic, entity *ecs.Entity) {
	switch ch.State {
	case data.OnLadder, data.OnBar, data.Flying, data.Leaping:
		return
	}
	if (ch.Object.Flip && !ch.Flags.LeftWall) ||
		(!ch.Object.Flip && !ch.Flags.RightWall) {
		// throw if space to throw
		// set action
		ch.State = data.DoingAction
		ch.Flags.ItemAction = data.ThrowBox
		sfx.SoundPlayer.PlaySound(constants.SFXThrow, 2.)
		// update box
		DropItem(ch)
		if c, okC := entity.GetComponentData(myecs.Dynamic); okC {
			box := c.(*data.Dynamic)
			chX, chY := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			chT := data.CurrLevel.Tiles.Get(chX, chY)
			box.Object.Pos.X = chT.Object.Pos.X
			box.Object.Pos.Y = chT.Object.Pos.Y + 1
			box.LastTile = chT
			box.Flags.Thrown = true
			box.ACounter = 0
			if ch.Object.Flip {
				box.Flags.JumpL = true
			} else {
				box.Flags.JumpR = true
			}
		}
	} else {
		// just drop
		DropItem(ch)
	}
}

func BoxBonk(p int, ch *data.Dynamic, entity *ecs.Entity) {
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
		}
	}
}

func CreateKey(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 14))
	obj.Layer = 14
	color := tile.Metadata.Color
	theKey := &data.BasicItem{
		Object: obj,
		Sprite: img.NewSprite(key, constants.TileBatch),
		PickUp: data.NewPickUp(fmt.Sprintf("Key (%s)", color), 5, tile.Metadata.Color),
		Action: KeyAction(color),
		Color:  color,
	}
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, theKey.Sprite)
	e.AddComponent(myecs.PickUp, theKey.PickUp)
	e.AddComponent(myecs.Action, theKey.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func KeyAction(color data.ItemColor) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if KeyUnlock(ch.Object.Pos, color) {
			DropItem(ch)
			myecs.Manager.DisposeEntity(entity)
		}
	})
}

func KeyUnlock(chPos pixel.Vec, color data.ItemColor) bool {
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
				sfx.SoundPlayer.PlaySound(constants.SFXKey, 0)
				return true
			}
		}
	}
	return false
}

func CreateJetpack(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	jetpack := &data.Jetpack{
		Object:   obj,
		PickUp:   data.NewPickUp("Jetpack", 5, tile.Metadata.Color),
		Entity:   e,
		Metadata: tile.Metadata,
		Origin:   tile.Coords,
	}
	jetpack.Action = JetpackAction(jetpack)
	jpa := reanimator.NewBatchAnimation("jetpack", img.Batchers[constants.TileBatch], key, reanimator.Loop)
	jpr := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "jetpack_regen", []int{0, 1, 2, 3, 3, 4, 5, 6, 6}, reanimator.Tran)
	jpr.SetEndTrigger(func() {
		jetpack.Regen = false
	})
	jetpack.Anim = reanimator.New(reanimator.NewSwitch().
		AddAnimation(jpa).
		AddAnimation(jpr).
		AddNull("none").
		SetChooseFn(func() string {
			if jetpack.Waiting {
				return "none"
			} else if jetpack.Regen {
				return "regen"
			} else {
				return "jetpack"
			}
		}), "jetpack")
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, jetpack.Anim)
	e.AddComponent(myecs.Animated, jetpack.Anim)
	e.AddComponent(myecs.PickUp, jetpack.PickUp)
	e.AddComponent(myecs.Action, jetpack.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func JetpackAction(jetpack *data.Jetpack) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		DropItem(ch)
		// change the player's state to flying
		ch.State = data.Flying
		ch.Flags.Flying = true
		// set the jetpack vars
		jetpack.Using = true
		jetpack.Counter = 0
		// remove the pickup component
		jetpack.Entity.RemoveComponent(myecs.PickUp)
		// set the same animation frame for the player as the jetpack
		ch.AnInt = jetpack.Anim.GetCurrentAnim().Step + 1
		// add the player as a parent
		jetpack.Entity.AddComponent(myecs.Parent, ch.Object)
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if jetpack.Using {
				if data.CurrLevel.FrameChange {
					if jetpack.Counter%2 == 0 {
						// update timer visuals
						x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
						var txPos pixel.Vec
						tile := data.CurrLevel.Tiles.Get(x, y+1)
						if tile == nil {
							tile = data.CurrLevel.Tiles.Get(x, y-1)
							txPos = ch.Object.Pos.Add(pixel.V(0, -world.TileSize))
						} else {
							txPos = ch.Object.Pos.Add(pixel.V(0, world.TileSize))
						}
						timer := jetpack.Metadata.Timer - jetpack.Counter/2
						pColor := constants.Player1Color
						switch p {
						case 1:
							pColor = constants.Player2Color
						case 2:
							pColor = constants.Player3Color
						case 3:
							pColor = constants.Player4Color
						}
						data.NewFloatingText().
							WithPos(txPos).
							WithColor(pixel.ToRGBA(constants.ColorWhite)).
							WithShadow(pixel.ToRGBA(pColor)).
							WithText(fmt.Sprintf("%d", timer)).
							WithTimer(1)
					}
					jetpack.Counter++
					if jetpack.Counter > jetpack.Metadata.Timer*2 {
						ch.Flags.Flying = false
						if jetpack.Metadata.Regenerate {
							jetpack.Using = false
							jetpack.Waiting = true
							jetpack.Counter = 0
							jetpack.Entity.RemoveComponent(myecs.Parent)
						} else {
							myecs.Manager.DisposeEntity(jetpack.Entity)
						}
					}
				}
			} else if jetpack.Waiting {
				if reanimator.FrameSwitch {
					jetpack.Counter++
				}
				delay := constants.ItemRegen * (jetpack.Metadata.RegenDelay + 2)
				if jetpack.Counter > delay && data.CurrLevel.FrameChange {
					jetpack.Object.Pos = world.MapToWorld(jetpack.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					jetpack.Regen = true
					jetpack.Waiting = false
					jetpack.Entity.RemoveComponent(myecs.Update)
					jetpack.Entity.AddComponent(myecs.PickUp, jetpack.PickUp)
				}
			}
		}))
	})
}

func CreateDisguise(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	disguise := &data.Disguise{
		Object:   obj,
		PickUp:   data.NewPickUp("Jetpack", 5, tile.Metadata.Color),
		Entity:   e,
		Metadata: tile.Metadata,
		Origin:   tile.Coords,
	}
	disguise.Action = DonDisguise(disguise)
	regenA := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "disguise_regen", []int{0, 1, 2, 3, 3, 4, 5, 6, 6, 6, 6}, reanimator.Tran)
	regenA.SetTriggerAll(func() {
		switch regenA.Step {
		case 7, 8:
			obj.Pos.Y = pos.Y + 2
		case 9:
			obj.Pos.Y = pos.Y
		case 10:
			obj.Pos.Y = pos.Y + 1
		}
	})
	regenA.SetEndTrigger(func() {
		obj.Pos.Y = pos.Y
		disguise.Regen = false
	})
	doff := reanimator.NewBatchAnimation("doff", img.Batchers[constants.TileBatch], "disguise_doff", reanimator.Tran)
	doff.SetEndTrigger(func() {
		if disguise.Metadata.Regenerate {
			disguise.Counter = 0
			disguise.Waiting = true
			disguise.Doff = false
		} else {
			myecs.Manager.DisposeEntity(disguise.Entity)
		}
	})
	disguise.Anim = reanimator.New(reanimator.NewSwitch().
		AddAnimation(regenA).
		AddAnimation(reanimator.NewBatchSprite("disguise", img.Batchers[constants.TileBatch], key, reanimator.Hold)).
		AddAnimation(doff).
		AddNull("none").
		SetChooseFn(func() string {
			if disguise.Waiting || disguise.Using {
				return "none"
			} else if disguise.Doff {
				return "doff"
			} else if disguise.Regen {
				return "regen"
			} else {
				return "disguise"
			}
		}), "disguise")
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, disguise.Anim)
	e.AddComponent(myecs.Animated, disguise.Anim)
	e.AddComponent(myecs.PickUp, disguise.PickUp)
	e.AddComponent(myecs.Action, disguise.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func DonDisguise(disguise *data.Disguise) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Falling || ch.State == data.Jumping {
			return
		}
		DropItem(ch)
		// remove the pickup component
		disguise.Entity.RemoveComponent(myecs.PickUp)
		// set action
		ch.State = data.DoingAction
		ch.Flags.ItemAction = data.DonDisguise
		// set the disguise vars
		disguise.Using = true
		disguise.Counter = 0
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if disguise.Using {
				if ch.State == data.Dead || (ch.State == data.DoingAction && ch.Flags.ItemAction == data.Hiding) {
					if disguise.Metadata.Regenerate {
						ch.Flags.Disguised = false
						disguise.Using = false
						disguise.Doff = true
						disguise.Counter = 0
						disguise.Object.Pos = ch.Object.Pos
					}
				} else if data.CurrLevel.FrameChange {
					if disguise.Counter%2 == 0 {
						// update timer visuals
						x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
						var txPos pixel.Vec
						tile := data.CurrLevel.Tiles.Get(x, y+1)
						if tile == nil {
							tile = data.CurrLevel.Tiles.Get(x, y-1)
							txPos = ch.Object.Pos.Add(pixel.V(0, -world.TileSize))
						} else {
							txPos = ch.Object.Pos.Add(pixel.V(0, world.TileSize))
						}
						timer := disguise.Metadata.Timer - disguise.Counter/2
						pColor := constants.Player1Color
						switch p {
						case 1:
							pColor = constants.Player2Color
						case 2:
							pColor = constants.Player3Color
						case 3:
							pColor = constants.Player4Color
						}
						data.NewFloatingText().
							WithPos(txPos).
							WithColor(pixel.ToRGBA(constants.ColorWhite)).
							WithShadow(pixel.ToRGBA(pColor)).
							WithText(fmt.Sprintf("%d", timer)).
							WithTimer(1)
					}
					disguise.Counter++
					if disguise.Counter > disguise.Metadata.Timer*2 {
						ch.Flags.Disguised = false
						disguise.Using = false
						disguise.Doff = true
						disguise.Counter = 0
						disguise.Object.Pos = ch.Object.Pos
					}
				}
			} else if disguise.Waiting {
				if reanimator.FrameSwitch {
					disguise.Counter++
				}
				delay := constants.ItemRegen * (disguise.Metadata.RegenDelay + 2)
				if disguise.Counter > delay && data.CurrLevel.FrameChange {
					disguise.Object.Pos = world.MapToWorld(disguise.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					disguise.Regen = true
					disguise.Waiting = false
					disguise.Entity.RemoveComponent(myecs.Update)
					disguise.Entity.AddComponent(myecs.PickUp, disguise.PickUp)
				}
			}
		}))
	})
}

func CreateDrill(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	drill := &data.BasicItem{
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp("Drill", 5, tile.Metadata.Color),
		Entity:   e,
		Metadata: tile.Metadata,
		Origin:   tile.Coords,
	}
	drill.Action = UseDrill(drill)
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, drill.Sprite)
	e.AddComponent(myecs.PickUp, drill.PickUp)
	e.AddComponent(myecs.Action, drill.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func UseDrill(drill *data.BasicItem) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Grounded {
			// check if on correct ground
			xa, ya := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			tileA := data.CurrLevel.Tiles.Get(xa, ya)
			drillTileA := data.CurrLevel.Tiles.Get(xa, ya-1)
			if drillTileA != nil && drillTileA.Block == data.BlockBedrock {
				// set action
				ch.State = data.DoingAction
				ch.Flags.ItemAction = data.DrillStart
				ch.Object.SetPos(tileA.Object.Pos)
				// set vars
				drill.Using = true
				entity.AddComponent(myecs.Update, data.NewFn(func() {
					if ch.Flags.CheckAction {
						ch.Flags.CheckAction = false
						if ch.Flags.ItemAction == data.DrillStart || (ch.Flags.ItemAction == data.Drilling && ch.Actions.Action) {
							x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
							tile := data.CurrLevel.Tiles.Get(x, y)
							drillTile := data.CurrLevel.Tiles.Get(x, y-1)
							if drillTile != nil && drillTile.Block == data.BlockBedrock {
								ch.Flags.ItemAction = data.Drilling
								ch.Object.SetPos(tile.Object.Pos.Add(pixel.V(0., 1.)))
								drillTile.Flags.Collapse = true
								drillTile.Metadata.Regenerate = false
								drillTile.Counter = constants.DrillCounter
								AddMaskWithTrigger(drillTile, "drill_mask", false, false, func() {
									ch.Flags.CheckAction = true
									RemoveMask(drillTile)
								})
							}
						} else {
							ch.Flags.ItemAction = data.NoItemAction
							entity.RemoveComponent(myecs.Update)
						}
					}
				}))
			}
		}
	})
}

func CreateFlamethrower(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	flamethrower := &data.BasicItem{
		Object:   obj,
		PickUp:   data.NewPickUp("Flamethrower", 5, tile.Metadata.Color),
		Entity:   e,
		Metadata: tile.Metadata,
		Origin:   tile.Coords,
	}
	flamethrower.Action = FlamethrowerAction(flamethrower)
	fta := reanimator.NewBatchAnimation("flamethrower", img.Batchers[constants.TileBatch], key, reanimator.Loop)
	ftr := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "flamethrower_regen", []int{0, 1, 2, 3, 3, 4, 5, 6, 6}, reanimator.Tran)
	ftr.SetEndTrigger(func() {
		flamethrower.Regen = false
	})
	flamethrower.Anim = reanimator.New(reanimator.NewSwitch().
		AddAnimation(fta).
		AddAnimation(ftr).
		AddNull("none").
		SetChooseFn(func() string {
			if flamethrower.Waiting || flamethrower.Using {
				return "none"
			} else if flamethrower.Regen {
				return "regen"
			} else {
				return "flamethrower"
			}
		}), "flamethrower")
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, flamethrower.Anim)
	e.AddComponent(myecs.Animated, flamethrower.Anim)
	e.AddComponent(myecs.PickUp, flamethrower.PickUp)
	e.AddComponent(myecs.Action, flamethrower.Action)
	e.AddComponent(myecs.LvlElement, struct{}{})
}

func FlamethrowerAction(flamethrower *data.BasicItem) *data.Interact {
	numUses := 0
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Falling || ch.State == data.Jumping || ch.State == data.OnBar {
			return
		}
		// check if not directly behind bedrock
		x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
		x2 := x + 1
		x3 := x + 2
		flip := false
		if ch.Object.Flip {
			x2 = x - 1
			x3 = x - 2
			flip = true
		}
		tile := data.CurrLevel.Tiles.Get(x, y)
		if ch.State == data.Leaping && !tile.IsLadder() && tile.Block != data.BlockBar {
			return // can't do it if leaping off ladders or bar
		}
		firstTile := data.CurrLevel.Tiles.Get(x2, y)
		if firstTile != nil && firstTile.Block != data.BlockBedrock {
			// change the player's state to flamethrower action
			ch.State = data.DoingAction
			ch.Flags.ItemAction = data.FireFlamethrower
			ch.Object.SetPos(tile.Object.Pos)
			// set the flamethrower vars
			flamethrower.Using = true
			flamethrower.Counter = 0
			numUses++
			// create flame animations
			secondTile := data.CurrLevel.Tiles.Get(x3, y)
			firstFlame := "flames_l"
			secondFlame := "flames_r"
			obj1 := object.New()
			obj1.Layer = 34
			obj1.Pos = firstTile.Object.Pos
			obj1.SetRect(pixel.R(0, 0, 16, 16))
			obj1.Flip = flip
			e1 := myecs.Manager.NewEntity()
			a1 := reanimator.NewBatchAnimation(firstFlame, img.Batchers[constants.TileBatch], firstFlame, reanimator.Tran)
			a1.SetEndTrigger(func() {
				myecs.Manager.DisposeEntity(e1)
			})
			anim1 := reanimator.NewSimple(a1)
			e1.AddComponent(myecs.Object, obj1)
			e1.AddComponent(myecs.Drawable, anim1)
			e1.AddComponent(myecs.Animated, anim1)
			e1.AddComponent(myecs.Temp, myecs.ClearFlag(false))
			if secondTile != nil {
				obj2 := object.New()
				obj2.Layer = 34
				obj2.Pos = secondTile.Object.Pos
				obj2.SetRect(pixel.R(0, 0, 16, 16))
				obj2.Flip = flip
				e2 := myecs.Manager.NewEntity()
				a2 := reanimator.NewBatchAnimation(secondFlame, img.Batchers[constants.TileBatch], secondFlame, reanimator.Tran)
				a2.SetEndTrigger(func() {
					myecs.Manager.DisposeEntity(e2)
				})
				anim2 := reanimator.NewSimple(a2)
				e2.AddComponent(myecs.Object, obj2)
				e2.AddComponent(myecs.Drawable, anim2)
				e2.AddComponent(myecs.Animated, anim2)
				e2.AddComponent(myecs.Temp, myecs.ClearFlag(false))
			}
			entity.AddComponent(myecs.Update, data.NewFn(func() {
				if flamethrower.Using {
					if reanimator.FrameSwitch {
						flamethrower.Counter++
						if flamethrower.Counter > constants.FlamethrowerCnt {
							flamethrower.Using = false
							flamethrower.Counter = 0
							ch.Flags.ItemAction = data.NoItemAction
							if flamethrower.Metadata.Timer > 0 && numUses >= flamethrower.Metadata.Timer {
								DropItem(ch)
								numUses = 0
								if flamethrower.Metadata.Regenerate {
									// remove the pickup component
									flamethrower.Entity.RemoveComponent(myecs.PickUp)
									flamethrower.Waiting = true
								} else {
									myecs.Manager.DisposeEntity(flamethrower.Entity)
								}
							}
						} else {
							// kill players and enemies
							for _, resultC := range myecs.Manager.Query(myecs.IsCharacter) {
								_, okCO := resultC.Components[myecs.Object].(*object.Object)
								ch2, okC := resultC.Components[myecs.Dynamic].(*data.Dynamic)
								if okCO && okC && ch2.State != data.Dead {
									chX, chY := world.WorldToMap(ch2.Object.Pos.X, ch2.Object.Pos.Y)
									ch2Coords := world.NewCoords(chX, chY)
									if ch2Coords == firstTile.Coords || (secondTile != nil && ch2Coords == secondTile.Coords) {
										tile2 := data.CurrLevel.Tiles.Get(chX, chY)
										ch2.Object.Pos.X = tile2.Object.Pos.X
										ch2.Object.Pos.Y = tile2.Object.Pos.Y
										ch2.Flags.Blow = true
										ch2.Flags.Hit = true
										ch2.State = data.Hit
										ch2.Object.Layer = 35
									}
								}
							}
							if flamethrower.Counter == 1 {
								// light bombs
								for _, resultB := range myecs.Manager.Query(myecs.IsBomb) {
									objB, okO := resultB.Components[myecs.Object].(*object.Object)
									b, okB := resultB.Components[myecs.Bomb].(*data.Bomb)
									if okO && okB && !objB.Hidden && !b.Waiting && !b.Regen {
										chX, chY := world.WorldToMap(objB.Pos.X, objB.Pos.Y)
										bCoords := world.NewCoords(chX, chY)
										if bCoords == firstTile.Coords || (secondTile != nil && bCoords == secondTile.Coords) {
											tileB := data.CurrLevel.Tiles.Get(chX, chY)
											objB.Pos.X = tileB.Object.Pos.X
											objB.Pos.Y = tileB.Object.Pos.Y
											LightBomb(b, tileB)
										}
									}
								}
								// destroy turf
								if firstTile.Block == data.BlockTurf || firstTile.Block == data.BlockCracked || firstTile.Block == data.BlockFall || firstTile.Block == data.BlockClose {
									firstTile.Flags.Collapse = true
									firstTile.Counter = constants.CollapseCounter
								}
								if secondTile != nil && (secondTile.Block == data.BlockTurf || secondTile.Block == data.BlockCracked || secondTile.Block == data.BlockFall || secondTile.Block == data.BlockClose) {
									secondTile.Flags.Collapse = true
									secondTile.Counter = constants.CollapseCounter
								}
							}
						}
					}
				} else if flamethrower.Waiting {
					if reanimator.FrameSwitch {
						flamethrower.Counter++
					}
					delay := constants.ItemRegen * (flamethrower.Metadata.RegenDelay + 2)
					if flamethrower.Counter > delay && data.CurrLevel.FrameChange {
						flamethrower.Object.Pos = world.MapToWorld(flamethrower.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
						flamethrower.Regen = true
						flamethrower.Waiting = false
						flamethrower.Entity.RemoveComponent(myecs.Update)
						flamethrower.Entity.AddComponent(myecs.PickUp, flamethrower.PickUp)
					}
				}
			}))
		}
	})
}
