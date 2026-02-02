package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/data/death"
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

func RegenerateOrRemove(item *data.BasicItem) {
	if item.Metadata.Regenerate {
		item.Counter = 0
		item.Waiting = true
	} else {
		myecs.Manager.DisposeEntity(item.Entity)
	}
}

func DefaultAnim(item *data.BasicItem) {
	item.Anim = reanimator.New().
		AddAnimation(reanimator.NewBatchSprite("item", img.Batchers[constants.TileBatch], item.Sprite.Key, reanimator.Hold)).
		AddNull("none").
		SetChooseFn(func() string {
			if item.Waiting {
				return "none"
			} else {
				return "item"
			}
		}).SetDefault("item").Finish()
	item.Entity.AddComponent(myecs.Drawable, item.Anim)
	item.Entity.AddComponent(myecs.Animated, item.Anim)
}

func DefaultRegen(item *data.BasicItem) {
	if item.Metadata.Regenerate {
		if item.Delay == 0 {
			item.Delay = constants.ItemRegen * (item.Metadata.RegenDelay + 2)
		}
		item.Entity.AddComponent(myecs.Update, data.NewFn(func() {
			if item.Waiting {
				if reanimator.FrameSwitch {
					item.Counter++
				}
				if item.Counter > item.Delay && data.CurrLevel.FrameChange {
					item.Object.Pos = world.MapToWorld(item.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					item.Regen = true
					item.Waiting = false
					// smoke animation
					e := myecs.Manager.NewEntity()
					obj := object.New()
					obj.Pos = item.Object.Pos
					obj.Layer = 33
					smoke := reanimator.NewBatchAnimation("smoke", img.Batchers[constants.TileBatch], "smoke_regen_anim", reanimator.Tran)
					smoke.SetEndTrigger(func() {
						item.Regen = false
						myecs.Manager.DisposeEntity(e)
					})
					anim := reanimator.NewSimple(smoke)
					e.AddComponent(myecs.Object, obj).
						AddComponent(myecs.Animated, anim).
						AddComponent(myecs.Drawable, anim).
						AddComponent(myecs.Temp, myecs.ClearFlag(false))
				}
			}
		}))
	}
}

func CreateGem(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 10, 10))
	obj.Layer = 11
	if tile.Metadata.Buried {
		obj.Layer = 9
	}
	gemShimmer := false
	batch := img.Batchers[constants.TileBatch]
	color := tile.Metadata.Color
	gemSpr := reanimator.NewBatchSprite("gem", batch, key, reanimator.Hold)
	shimmer := reanimator.NewBatchAnimationCustom("shimmer", batch, fmt.Sprintf("%s_shimmer", key), []int{0, 0, 1, 1, 2, 2}, reanimator.Tran)
	shimmer.SetEndTrigger(func() {
		gemShimmer = false
	})
	anim := reanimator.New().
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
		}).SetDefault("gem").Finish()
	e := myecs.Manager.NewEntity()
	gem := &data.BasicItem{
		Name:     "Gem",
		Key:      key,
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp(5, tile.Metadata.Color),
		Entity:   e,
		Metadata: tile.Metadata,
		Origin:   tile.Coords,
		Layer:    11,
	}
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Drawable, anim).
		AddComponent(myecs.Animated, anim).
		AddComponent(myecs.Gem, struct{}{}).
		AddComponent(myecs.OnTouch, CollectGem(color, tile.Coords)).
		AddComponent(myecs.Item, gem).
		AddComponent(myecs.LvlElement, struct{}{})
}

func CollectGem(color data.ItemColor, c world.Coords) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers ||
			(ch.Color != color && color > data.NonPlayerRed) {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 1
		data.CurrLevelSess.PlayerStats[p].LGems++
		data.CurrLevelSess.GemsCollected = append(data.CurrLevelSess.GemsCollected, c)
		sfx.SoundPlayer.PlaySound(constants.SFXGem, -2.)
		myecs.Manager.DisposeEntity(entity)
	})
}

func CreateDoor(pos pixel.Vec, tile *data.Tile) {
	key := tile.SpriteString()
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 8))
	obj.Layer = 9
	door := &data.Door{
		Item: &data.BasicItem{
			Key:   key,
			Layer: 9,
		},
	}
	door.Item.Object = obj
	door.Item.Color = tile.Metadata.Color
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	switch tile.Block {
	case data.BlockDoorHidden:
		door.DoorType = data.Hidden
	case data.BlockDoorVisible:
		door.DoorType = data.Hidden
		e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	case data.BlockDoorLocked:
		door.DoorType = data.Locked
		e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	}
	interaction := EnterDoor(door.Item.Color, tile.Metadata.ExitIndex)
	//var interaction *data.Interact
	//switch door.Item.Color {
	//case data.PlayerBlue, data.PlayerGreen, data.PlayerPurple, data.PlayerOrange:
	//	interaction = EnterPlayerDoor(door.Item.Color, tile.Metadata.ExitIndex)
	//default:
	//	interaction = EnterDoor(tile.Metadata.ExitIndex)
	//}
	door.Item.Entity = e
	var anim *reanimator.Anim
	if door.DoorType == data.Locked {
		anim = reanimator.NewBatchAnimation("unlock", img.Batchers[constants.TileBatch], fmt.Sprintf("door_unlock%s", door.Item.Color.SpriteString()), reanimator.Tran)
		anim.SetEndTrigger(func() {
			e.RemoveComponent(myecs.Animated)
			e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("door_unlocked%s", door.Item.Color.SpriteString()), constants.TileBatch))
			e.AddComponent(myecs.OnTouch, interaction)
		})
	} else {
		anim = reanimator.NewBatchAnimation("open", img.Batchers[constants.TileBatch], fmt.Sprintf("door_visible%s", door.Item.Color.SpriteString()), reanimator.Tran)
		anim.SetEndTrigger(func() {
			e.RemoveComponent(myecs.Animated)
			e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("door_hidden%s", door.Item.Color.SpriteString()), constants.TileBatch))
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
				e.AddComponent(myecs.Drawable, img.NewSprite(fmt.Sprintf("door_unlock%s", door.Item.Color.SpriteString()), constants.TileBatch))
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

//func EnterPlayerDoor(exitIndex int) *data.Interact {
//	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
//		if p < 0 || p >= constants.MaxPlayers {
//			return
//		}
//		data.CurrLevelSess.PlayerStats[p].LScore += 12
//		sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
//		data.CurrLevel.Open = true
//		data.CurrLevel.ExitIndex = exitIndex
//	})
//}

func EnterDoor(color data.ItemColor, exitIndex int) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers ||
			(color > data.NonPlayerRed && ch.Color != color) ||
			ch.Actions.Direction != data.Up {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 12
		sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
		data.CurrLevel.Complete = true
		data.CurrLevel.ExitIndex = exitIndex
	})
}

func CreateJumpBoots(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	jumpBoots := &data.BasicItem{
		Name:     "Jump Boots",
		Key:      key,
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp(5, metadata.Color),
		Entity:   e,
		Metadata: metadata,
		Origin:   coords,
		Layer:    12,
	}
	jumpBoots.Metadata = data.CopyMetadata(metadata)
	jumpBoots.Metadata.Timer = -1
	jumpBoots.Action = Jump()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	//e.AddComponent(myecs.Drawable, jumpBoots.Sprite)
	e.AddComponent(myecs.PickUp, jumpBoots.PickUp)
	e.AddComponent(myecs.Action, jumpBoots.Action)
	e.AddComponent(myecs.LvlElement, jumpBoots)
	e.AddComponent(myecs.Item, jumpBoots)
	DefaultRegen(jumpBoots)
	DefaultAnim(jumpBoots)
	return jumpBoots
}

func Jump() *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Grounded {
			x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			left := data.CurrLevel.Get(x-1, y)
			right := data.CurrLevel.Get(x+1, y)
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
			tile := data.CurrLevel.Get(x, y)
			ch.LastTile = tile
			ch.State = data.Jumping
			ch.Object.Pos.X = tile.Object.Pos.X
			ch.Object.Pos.Y = tile.Object.Pos.Y
			sfx.SoundPlayer.PlaySound(constants.SFXJump, 0.)
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

func CreateBox(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	s := obj.Pos.Y
	smash := &s
	e := myecs.Manager.NewEntity()
	box := &data.BasicItem{
		Name:     "Box",
		Key:      key,
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp(5, metadata.Color),
		Entity:   e,
		Metadata: metadata,
		Origin:   coords,
		Layer:    12,
	}
	box.Metadata = data.CopyMetadata(metadata)
	box.Metadata.Timer = -1
	e.AddComponent(myecs.Object, box.Object)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	//e.AddComponent(myecs.Drawable, box.Sprite)
	e.AddComponent(myecs.StandOn, struct{}{})
	e.AddComponent(myecs.Smash, smash)
	e.AddComponent(myecs.OnTouch, data.NewInteract(BoxBonk))
	e.AddComponent(myecs.PickUp, box.PickUp)
	e.AddComponent(myecs.Action, data.NewInteract(BoxAction))
	e.AddComponent(myecs.LvlElement, box)
	e.AddComponent(myecs.Item, box)
	dyn := data.NewDynamic()
	dyn.LastTile = data.CurrLevel.Get(coords.X, coords.Y)
	dyn.Object = obj
	dyn.Entity = e
	dyn.Flags.NoLadders = true
	e.AddComponent(myecs.Dynamic, dyn)
	DefaultRegen(box)
	DefaultAnim(box)
	return box
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
		ClearInv(ch)
		if c, okC := entity.GetComponentData(myecs.Dynamic); okC {
			box := c.(*data.Dynamic)
			chX, chY := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			chT := data.CurrLevel.Get(chX, chY)
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
			ch.Flags.Death = death.Dying
			ch.State = data.Hit
		}
	}
}

func CreateKey(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 14))
	obj.Layer = 14
	color := metadata.Color
	theKey := &data.BasicItem{
		Name:   fmt.Sprintf("Key (%s)", color),
		Key:    key,
		Object: obj,
		Sprite: img.NewSprite(key, constants.TileBatch),
		PickUp: data.NewPickUp(5, metadata.Color),
		Action: KeyAction(color),
		Color:  color,
		Origin: coords,
		Layer:  14,
	}
	theKey.Metadata = data.CopyMetadata(metadata)
	theKey.Metadata.Timer = -1
	theKey.Metadata.Regenerate = false
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	//e.AddComponent(myecs.Drawable, theKey.Sprite)
	e.AddComponent(myecs.PickUp, theKey.PickUp)
	e.AddComponent(myecs.Action, theKey.Action)
	e.AddComponent(myecs.LvlElement, theKey)
	e.AddComponent(myecs.Item, theKey)
	theKey.Entity = e
	DefaultRegen(theKey)
	DefaultAnim(theKey)
	return theKey
}

func KeyAction(color data.ItemColor) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if KeyUnlock(ch.Object.Pos, color) {
			ClearInv(ch)
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
				d.Item.Color == color {
				d.Unlock = true
				sfx.SoundPlayer.PlaySound(constants.SFXKey, 0)
				return true
			}
		}
	}
	return false
}

func CreateJetpack(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	jetpack := &data.BasicItem{
		Name:     "Jetpack",
		Key:      key,
		Object:   obj,
		PickUp:   data.NewPickUp(5, metadata.Color),
		Entity:   e,
		Metadata: metadata,
		Origin:   coords,
		Layer:    12,
	}
	jetpack.Metadata = data.CopyMetadata(metadata)
	jetpack.Action = JetpackAction(jetpack)
	jpa := reanimator.NewBatchAnimation("jetpack", img.Batchers[constants.TileBatch], key, reanimator.Loop)
	jpr := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "jetpack_regen", []int{0, 1, 2, 3, 3, 4, 5, 6, 6}, reanimator.Tran)
	jpr.SetEndTrigger(func() {
		jetpack.Regen = false
	})
	jetpack.Anim = reanimator.New().
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
		}).SetDefault("jetpack").Finish()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, jetpack.Anim)
	e.AddComponent(myecs.Animated, jetpack.Anim)
	e.AddComponent(myecs.PickUp, jetpack.PickUp)
	e.AddComponent(myecs.Action, jetpack.Action)
	e.AddComponent(myecs.LvlElement, jetpack)
	e.AddComponent(myecs.Item, jetpack)
	jetpack.Delay = constants.ItemRegen * (jetpack.Metadata.RegenDelay + 2)
	e.AddComponent(myecs.Update, data.NewFn(func() {
		if jetpack.Waiting {
			if reanimator.FrameSwitch {
				jetpack.Counter++
			}
			if jetpack.Counter > jetpack.Delay && data.CurrLevel.FrameChange {
				jetpack.Object.Pos = world.MapToWorld(jetpack.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
				jetpack.Regen = true
				jetpack.Waiting = false
				jetpack.Entity.AddComponent(myecs.PickUp, jetpack.PickUp)
			}
		}
	}))
	return jetpack
}

func JetpackAction(jetpack *data.BasicItem) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		ClearInv(ch)
		// change the player's state to flying
		ch.State = data.Flying
		ch.Flags.Flying = true
		// set the jetpack vars
		jetpack.Using = true
		jetpack.Counter = 0
		id := sfx.SoundPlayer.PlaySound(constants.SFXJetpackStart, -2.)
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
						tile := data.CurrLevel.Get(x, y+1)
						if tile == nil {
							tile = data.CurrLevel.Get(x, y-1)
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
					if jetpack.Counter+1 <= jetpack.Metadata.Timer*2 {
						sfx.SoundPlayer.KillSound(id)
						id = sfx.SoundPlayer.PlaySound(constants.SFXJetpackEnd, -2.)
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
				if jetpack.Counter > jetpack.Delay && data.CurrLevel.FrameChange {
					jetpack.Object.Pos = world.MapToWorld(jetpack.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					jetpack.Regen = true
					jetpack.Waiting = false
					jetpack.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if jetpack.Waiting {
							if reanimator.FrameSwitch {
								jetpack.Counter++
							}
							if jetpack.Counter > jetpack.Delay && data.CurrLevel.FrameChange {
								jetpack.Object.Pos = world.MapToWorld(jetpack.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
								jetpack.Regen = true
								jetpack.Waiting = false
								jetpack.Entity.AddComponent(myecs.PickUp, jetpack.PickUp)
							}
						}
					}))
					jetpack.Entity.AddComponent(myecs.PickUp, jetpack.PickUp)
				}
			}
		}))
	})
}

func CreateDisguise(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	e := myecs.Manager.NewEntity()
	disguise := &data.Disguise{
		Item: &data.BasicItem{
			Name:     "Disguise",
			Key:      key,
			Object:   obj,
			PickUp:   data.NewPickUp(5, metadata.Color),
			Entity:   e,
			Metadata: metadata,
			Origin:   coords,
			Layer:    12,
		},
	}
	obj.Layer = disguise.Item.Layer
	disguise.Item.Metadata = data.CopyMetadata(metadata)
	disguise.Item.Action = DonDisguise(disguise)
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
		disguise.Item.Regen = false
		disguise.Doff = false
	})
	doff := reanimator.NewBatchAnimation("doff", img.Batchers[constants.TileBatch], "disguise_doff", reanimator.Tran)
	doff.SetEndTrigger(func() {
		if disguise.Item.Metadata.Regenerate {
			disguise.Item.Counter = 0
			disguise.Item.Waiting = true
			disguise.Doff = false
		} else {
			myecs.Manager.DisposeEntity(disguise.Item.Entity)
		}
	})
	disguise.Item.Anim = reanimator.New().
		AddAnimation(regenA).
		AddAnimation(reanimator.NewBatchSprite("disguise", img.Batchers[constants.TileBatch], key, reanimator.Hold)).
		AddAnimation(doff).
		AddNull("none").
		SetChooseFn(func() string {
			if disguise.Item.Waiting || disguise.Item.Using {
				return "none"
			} else if disguise.Doff {
				return "doff"
			} else if disguise.Item.Regen {
				return "regen"
			} else {
				return "disguise"
			}
		}).SetDefault("disguise").Finish()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, disguise.Item.Anim)
	e.AddComponent(myecs.Animated, disguise.Item.Anim)
	e.AddComponent(myecs.PickUp, disguise.Item.PickUp)
	e.AddComponent(myecs.Action, disguise.Item.Action)
	e.AddComponent(myecs.LvlElement, disguise)
	e.AddComponent(myecs.Item, disguise.Item)
	disguise.Item.Delay = constants.ItemRegen * (disguise.Item.Metadata.RegenDelay + 2)
	e.AddComponent(myecs.Update, data.NewFn(func() {
		if disguise.Item.Waiting {
			if reanimator.FrameSwitch {
				disguise.Item.Counter++
			}
			if disguise.Item.Counter > disguise.Item.Delay && data.CurrLevel.FrameChange {
				disguise.Item.Object.Pos = world.MapToWorld(disguise.Item.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
				disguise.Item.Regen = true
				disguise.Item.Waiting = false
				disguise.Item.Entity.AddComponent(myecs.PickUp, disguise.Item.PickUp)
			}
		}
	}))
	return disguise.Item
}

func DonDisguise(disguise *data.Disguise) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Falling || ch.State == data.Jumping {
			return
		}
		ClearInv(ch)
		// remove the pickup component
		disguise.Item.Entity.RemoveComponent(myecs.PickUp)
		// set action
		ch.State = data.DoingAction
		ch.Flags.ItemAction = data.DonDisguise
		// set the disguise vars
		disguise.Item.Using = true
		disguise.Item.Counter = 0
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if disguise.Item.Using {
				if ch.State == data.Dead || (ch.State == data.DoingAction && ch.Flags.ItemAction == data.Hiding) {
					if disguise.Item.Metadata.Regenerate {
						ch.Flags.Disguised = false
						disguise.Item.Using = false
						disguise.Doff = true
						disguise.Item.Counter = 0
						disguise.Item.Object.Pos = ch.Object.Pos
					}
				} else if data.CurrLevel.FrameChange {
					if disguise.Item.Counter%2 == 0 {
						// update timer visuals
						x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
						var txPos pixel.Vec
						tile := data.CurrLevel.Get(x, y+1)
						if tile == nil {
							tile = data.CurrLevel.Get(x, y-1)
							txPos = ch.Object.Pos.Add(pixel.V(0, -world.TileSize))
						} else {
							txPos = ch.Object.Pos.Add(pixel.V(0, world.TileSize))
						}
						timer := disguise.Item.Metadata.Timer - disguise.Item.Counter/2
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
					disguise.Item.Counter++
					if disguise.Item.Counter > disguise.Item.Metadata.Timer*2 {
						ch.Flags.Disguised = false
						disguise.Item.Using = false
						disguise.Doff = true
						disguise.Item.Counter = 0
						disguise.Item.Object.Pos = ch.Object.Pos
					}
				}
			} else if disguise.Item.Waiting {
				if reanimator.FrameSwitch {
					disguise.Item.Counter++
				}
				if disguise.Item.Counter > disguise.Item.Delay && data.CurrLevel.FrameChange {
					disguise.Item.Object.Pos = world.MapToWorld(disguise.Item.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					disguise.Item.Regen = true
					disguise.Item.Waiting = false
					disguise.Item.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if disguise.Item.Waiting {
							if reanimator.FrameSwitch {
								disguise.Item.Counter++
							}
							if disguise.Item.Counter > disguise.Item.Delay && data.CurrLevel.FrameChange {
								disguise.Item.Object.Pos = world.MapToWorld(disguise.Item.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
								disguise.Item.Regen = true
								disguise.Item.Waiting = false
								disguise.Item.Entity.AddComponent(myecs.PickUp, disguise.Item.PickUp)
							}
						}
					}))
					disguise.Item.Entity.AddComponent(myecs.PickUp, disguise.Item.PickUp)
				}
			}
		}))
	})
}

func CreateDrill(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	drill := &data.BasicItem{
		Name:     "Drill",
		Key:      key,
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp(5, metadata.Color),
		Entity:   e,
		Metadata: metadata,
		Origin:   coords,
		Layer:    12,
	}
	drill.Metadata = data.CopyMetadata(metadata)
	drill.Metadata.Timer = -1
	drill.Action = UseDrill(drill)
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, drill.Sprite)
	e.AddComponent(myecs.PickUp, drill.PickUp)
	e.AddComponent(myecs.Action, drill.Action)
	e.AddComponent(myecs.LvlElement, drill)
	e.AddComponent(myecs.Item, drill)
	DefaultRegen(drill)
	return drill
}

func UseDrill(drill *data.BasicItem) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Grounded || ch.State == data.Hiding || ch.State == data.OnLadder {
			// check if on correct ground
			xa, ya := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			tileA := data.CurrLevel.Get(xa, ya)
			drillTileA := data.CurrLevel.Get(xa, ya-1)
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
							tile := data.CurrLevel.Get(x, y)
							drillTile := data.CurrLevel.Get(x, y-1)
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
							DefaultRegen(drill)
						}
					}
				}))
			}
		}
	})
}

func CreateFlamethrower(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	flamethrower := &data.BasicItem{
		Name:     "Flamethrower",
		Key:      key,
		Object:   obj,
		PickUp:   data.NewPickUp(5, metadata.Color),
		Entity:   e,
		Metadata: metadata,
		Origin:   coords,
		Layer:    12,
	}
	flamethrower.Metadata = data.CopyMetadata(metadata)
	flamethrower.Action = FlamethrowerAction(flamethrower)
	fta := reanimator.NewBatchAnimation("flamethrower", img.Batchers[constants.TileBatch], key, reanimator.Loop)
	ftr := reanimator.NewBatchAnimationCustom("regen", img.Batchers[constants.TileBatch], "flamethrower_regen", []int{0, 1, 2, 3, 3, 4, 5, 6, 6}, reanimator.Tran)
	ftr.SetEndTrigger(func() {
		flamethrower.Regen = false
	})
	flamethrower.Anim = reanimator.New().
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
		}).SetDefault("flamethrower").Finish()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, flamethrower.Anim)
	e.AddComponent(myecs.Animated, flamethrower.Anim)
	e.AddComponent(myecs.PickUp, flamethrower.PickUp)
	e.AddComponent(myecs.Action, flamethrower.Action)
	e.AddComponent(myecs.LvlElement, flamethrower)
	e.AddComponent(myecs.Item, flamethrower)
	flamethrower.Delay = constants.ItemRegen * (flamethrower.Metadata.RegenDelay + 2)
	e.AddComponent(myecs.Update, data.NewFn(func() {
		if flamethrower.Waiting {
			if reanimator.FrameSwitch {
				flamethrower.Counter++
			}
			if flamethrower.Counter > flamethrower.Delay && data.CurrLevel.FrameChange {
				flamethrower.Object.Pos = world.MapToWorld(flamethrower.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
				flamethrower.Regen = true
				flamethrower.Waiting = false
				flamethrower.Entity.AddComponent(myecs.PickUp, flamethrower.PickUp)
			}
		}
	}))
	return flamethrower
}

func FlamethrowerAction(flamethrower *data.BasicItem) *data.Interact {
	flamethrower.Uses = 0
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Falling || ch.State == data.Jumping || ch.State == data.OnBar {
			return
		}
		// check if not directly behind bedrock
		x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
		x2 := x + 1
		x3 := x + 2
		x4 := x + 3
		flip := false
		if ch.Object.Flip {
			x2 = x - 1
			x3 = x - 2
			x4 = x - 3
			flip = true
		}
		tile := data.CurrLevel.Get(x, y)
		if ch.State == data.Leaping && !tile.IsLadder() && tile.Block != data.BlockBar {
			return // can't do it if leaping off ladders or bar
		}
		firstTile := data.CurrLevel.Get(x2, y)
		if firstTile == nil || firstTile.Block == data.BlockBedrock {
			return // can't do it if directly facing bedrock
		}
		// change the player's state to flamethrower action
		ch.State = data.DoingAction
		ch.Flags.ItemAction = data.FireFlamethrower
		ch.Object.SetPos(tile.Object.Pos)
		// set the flamethrower vars
		flamethrower.Using = true
		flamethrower.Counter = 0
		flamethrower.Uses++
		// create flame animations
		secondTile := data.CurrLevel.Get(x3, y)
		thirdTile := data.CurrLevel.Get(x4, y)
		CreateFlames(firstTile, secondTile, thirdTile, flip)
		// flamethrower sound
		sfx.SoundPlayer.PlaySound(constants.SFXFlamethrower, 0.)
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if flamethrower.Using {
				if reanimator.FrameSwitch {
					flamethrower.Counter++
					if flamethrower.Counter > constants.FlamethrowerCnt {
						flamethrower.Using = false
						flamethrower.Counter = 0
						ch.Flags.ItemAction = data.NoItemAction
						if flamethrower.Metadata.Timer > 0 && flamethrower.Uses >= flamethrower.Metadata.Timer {
							ClearInv(ch)
							flamethrower.Uses = 0
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
							if okCO && okC && !ch.Flags.Ignore && ch2.State != data.Dead && ch2.State != data.Hit && ch2.State != data.Waiting {
								chX, chY := world.WorldToMap(ch2.Object.Pos.X, ch2.Object.Pos.Y)
								ch2Coords := world.NewCoords(chX, chY)
								if ch2Coords == firstTile.Coords ||
									(secondTile != nil && ch2Coords == secondTile.Coords) ||
									(thirdTile != nil && ch2Coords == thirdTile.Coords) {
									tile2 := data.CurrLevel.Get(chX, chY)
									ch2.Object.Pos.X = tile2.Object.Pos.X
									ch2.Object.Pos.Y = tile2.Object.Pos.Y
									ch2.Flags.Death = death.Exploded
									ch2.State = data.Hit
									ch2.Object.Layer = 35
								}
							}
						}
						if flamethrower.Counter == 2 {
							// light bombs
							for _, resultB := range myecs.Manager.Query(myecs.IsBomb) {
								objB, okO := resultB.Components[myecs.Object].(*object.Object)
								b, okB := resultB.Components[myecs.Bomb].(*data.Bomb)
								if okO && okB && !objB.Hidden && !b.Item.Waiting && !b.Item.Regen {
									chX, chY := world.WorldToMap(objB.Pos.X, objB.Pos.Y)
									bCoords := world.NewCoords(chX, chY)
									if bCoords == firstTile.Coords ||
										(secondTile != nil && bCoords == secondTile.Coords) ||
										(thirdTile != nil && bCoords == thirdTile.Coords) {
										tileB := data.CurrLevel.Get(chX, chY)
										objB.Pos.X = tileB.Object.Pos.X
										objB.Pos.Y = tileB.Object.Pos.Y
										LightBomb(b, tileB)
									}
								}
							}
							// destroy turf
							if firstTile.Block == data.BlockTurf ||
								firstTile.Block == data.BlockCracked ||
								firstTile.Block == data.BlockFall ||
								firstTile.Block == data.BlockClose {
								firstTile.Flags.Collapse = true
								firstTile.Counter = constants.CollapseCounter
							}
							if secondTile != nil &&
								(secondTile.Block == data.BlockTurf ||
									secondTile.Block == data.BlockCracked ||
									secondTile.Block == data.BlockFall ||
									secondTile.Block == data.BlockClose) {
								secondTile.Flags.Collapse = true
								secondTile.Counter = constants.CollapseCounter
							}
							if thirdTile != nil &&
								secondTile.Block != data.BlockBedrock &&
								(thirdTile.Block == data.BlockTurf ||
									thirdTile.Block == data.BlockCracked ||
									thirdTile.Block == data.BlockFall ||
									thirdTile.Block == data.BlockClose) {
								thirdTile.Flags.Collapse = true
								thirdTile.Counter = constants.CollapseCounter
							}
						}
					}
				}
			} else if flamethrower.Waiting {
				if reanimator.FrameSwitch {
					flamethrower.Counter++
				}
				if flamethrower.Counter > flamethrower.Delay && data.CurrLevel.FrameChange {
					flamethrower.Object.Pos = world.MapToWorld(flamethrower.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					flamethrower.Regen = true
					flamethrower.Waiting = false
					flamethrower.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if flamethrower.Waiting {
							if reanimator.FrameSwitch {
								flamethrower.Counter++
							}
							if flamethrower.Counter > flamethrower.Delay && data.CurrLevel.FrameChange {
								flamethrower.Object.Pos = world.MapToWorld(flamethrower.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
								flamethrower.Regen = true
								flamethrower.Waiting = false
								flamethrower.Entity.AddComponent(myecs.PickUp, flamethrower.PickUp)
							}
						}
					}))
					flamethrower.Entity.AddComponent(myecs.PickUp, flamethrower.PickUp)
				}
			}
		}))
	})
}

func CreateFlames(firstTile, secondTile, thirdTile *data.Tile, flip bool) {
	firstFlame := "flames_l"
	secondFlame := "flames_m"
	thirdFlame := "flames_r"
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
		if thirdTile != nil {
			obj3 := object.New()
			obj3.Layer = 34
			obj3.Pos = thirdTile.Object.Pos
			obj3.SetRect(pixel.R(0, 0, 16, 16))
			obj3.Flip = flip
			e3 := myecs.Manager.NewEntity()
			a3 := reanimator.NewBatchAnimation(thirdFlame, img.Batchers[constants.TileBatch], thirdFlame, reanimator.Tran)
			a3.SetEndTrigger(func() {
				myecs.Manager.DisposeEntity(e3)
			})
			anim3 := reanimator.NewSimple(a3)
			e3.AddComponent(myecs.Object, obj3)
			e3.AddComponent(myecs.Drawable, anim3)
			e3.AddComponent(myecs.Animated, anim3)
			e3.AddComponent(myecs.Temp, myecs.ClearFlag(false))
		}
	}
}

func CreateGoopBucket(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	goopBucket := &data.BasicItem{
		Name:     "Goop Bucket",
		Key:      key,
		Object:   obj,
		Sprite:   img.NewSprite(key, constants.TileBatch),
		PickUp:   data.NewPickUp(5, metadata.Color),
		Entity:   e,
		Metadata: metadata,
		Origin:   coords,
		Layer:    12,
	}
	goopBucket.Metadata = data.CopyMetadata(metadata)
	goopBucket.Action = GoopBucketAction(goopBucket)

	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, goopBucket.Sprite)
	e.AddComponent(myecs.PickUp, goopBucket.PickUp)
	e.AddComponent(myecs.Action, goopBucket.Action)
	e.AddComponent(myecs.LvlElement, goopBucket)
	e.AddComponent(myecs.Item, goopBucket)
	DefaultRegen(goopBucket)
	goopBucket.Delay = constants.ItemRegen * (goopBucket.Metadata.RegenDelay + 2)
	e.AddComponent(myecs.Update, data.NewFn(func() {
		if goopBucket.Waiting {
			if reanimator.FrameSwitch {
				goopBucket.Counter++
			}
			if goopBucket.Counter > goopBucket.Delay && data.CurrLevel.FrameChange {
				goopBucket.Object.Pos = world.MapToWorld(goopBucket.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
				goopBucket.Regen = true
				goopBucket.Waiting = false
				goopBucket.Object.Hidden = false
				goopBucket.Entity.AddComponent(myecs.PickUp, goopBucket.PickUp)
			}
		}
	}))
	return goopBucket
}

func GoopBucketAction(goopBucket *data.BasicItem) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.State == data.Falling || ch.State == data.Jumping || ch.State == data.OnBar {
			return
		}
		x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
		x2 := x + 1
		x3 := x + 2
		x4 := x + 3
		flip := false
		if ch.Object.Flip {
			x2 = x - 1
			x3 = x - 2
			x4 = x - 3
			flip = true
		}
		tile := data.CurrLevel.Get(x, y)
		if ch.State == data.Leaping && !tile.IsLadder() && tile.Block != data.BlockBar {
			return // can't do it if leaping off ladders or bar
		}
		firstTile := data.CurrLevel.Get(x2, y)
		if firstTile.IsSolid() || firstTile.Block == data.BlockLiquid {
			return // can't do it if facing a wall
		}
		// change the player's state to goop bucket action
		ch.State = data.DoingAction
		ch.Flags.ItemAction = data.ThrowingGoop
		ch.Object.SetPos(tile.Object.Pos)
		// set the goop bucket vars
		goopBucket.Using = true
		goopBucket.Counter = 0
		// create goop animations
		secondTile := data.CurrLevel.Get(x3, y)
		thirdTile := data.CurrLevel.Get(x4, y)
		// goop sound
		//sfx.SoundPlayer.PlaySound(constants.SFXFlamethrower, 0.)
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if goopBucket.Using {
				if reanimator.FrameSwitch {
					goopBucket.Counter++
					if goopBucket.Counter > constants.GoopBucketCnt {
						goopBucket.Using = false
						MakeGoopThrow(firstTile, secondTile, thirdTile, flip)
						goopBucket.Counter = 0
						ClearInv(ch)
						goopBucket.Uses = 0
						goopBucket.Object.Hidden = true
						if goopBucket.Metadata.Regenerate {
							// remove the pickup component
							goopBucket.Entity.RemoveComponent(myecs.PickUp)
							goopBucket.Waiting = true
						} else {
							myecs.Manager.DisposeEntity(goopBucket.Entity)
						}
					}
				}
			} else if goopBucket.Waiting {
				if reanimator.FrameSwitch {
					goopBucket.Counter++
				}
				if goopBucket.Counter > goopBucket.Delay && data.CurrLevel.FrameChange {
					goopBucket.Object.Pos = world.MapToWorld(goopBucket.Origin).Add(pixel.V(world.HalfSize, world.HalfSize))
					goopBucket.Regen = true
					goopBucket.Waiting = false
					goopBucket.Object.Hidden = false
					entity.AddComponent(myecs.PickUp, goopBucket.PickUp)
					// smoke animation
					e := myecs.Manager.NewEntity()
					obj := object.New()
					obj.Pos = goopBucket.Object.Pos
					obj.Layer = 33
					smoke := reanimator.NewBatchAnimation("smoke", img.Batchers[constants.TileBatch], "smoke_regen_anim", reanimator.Tran)
					smoke.SetEndTrigger(func() {
						goopBucket.Regen = false
						myecs.Manager.DisposeEntity(e)
					})
					anim := reanimator.NewSimple(smoke)
					e.AddComponent(myecs.Object, obj).
						AddComponent(myecs.Animated, anim).
						AddComponent(myecs.Drawable, anim).
						AddComponent(myecs.Temp, myecs.ClearFlag(false))
				}
			}
		}))
	})
}

func MakeGoopThrow(firstTile, secondTile, thirdTile *data.Tile, flip bool) {
	firstGoop := "goop_throw_l"
	secondGoop := "goop_throw_m"
	thirdGoop := "goop_throw_r"
	obj1 := object.New()
	obj1.Layer = 34
	obj1.Pos = firstTile.Object.Pos
	obj1.SetRect(pixel.R(0, 0, 16, 16))
	obj1.Flip = flip
	e1 := myecs.Manager.NewEntity()
	a1 := reanimator.NewBatchAnimation(firstGoop, img.Batchers[constants.TileBatch], firstGoop, reanimator.Tran)
	a1.SetEndTrigger(func() {
		CreateFallingGoop(firstTile, flip)
		myecs.Manager.DisposeEntity(e1)
	})
	anim1 := reanimator.NewSimple(a1)
	e1.AddComponent(myecs.Object, obj1)
	e1.AddComponent(myecs.Drawable, anim1)
	e1.AddComponent(myecs.Animated, anim1)
	e1.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	if secondTile != nil && !secondTile.IsSolid() && secondTile.Block != data.BlockLiquid {
		obj2 := object.New()
		obj2.Layer = 34
		obj2.Pos = secondTile.Object.Pos
		obj2.SetRect(pixel.R(0, 0, 16, 16))
		obj2.Flip = flip
		e2 := myecs.Manager.NewEntity()
		a2 := reanimator.NewBatchAnimation(secondGoop, img.Batchers[constants.TileBatch], secondGoop, reanimator.Tran)
		a2.SetEndTrigger(func() {
			CreateFallingGoop(secondTile, flip)
			myecs.Manager.DisposeEntity(e2)
		})
		anim2 := reanimator.NewSimple(a2)
		e2.AddComponent(myecs.Object, obj2)
		e2.AddComponent(myecs.Drawable, anim2)
		e2.AddComponent(myecs.Animated, anim2)
		e2.AddComponent(myecs.Temp, myecs.ClearFlag(false))
		if thirdTile != nil && !thirdTile.IsSolid() && thirdTile.Block != data.BlockLiquid {
			obj3 := object.New()
			obj3.Layer = 34
			obj3.Pos = thirdTile.Object.Pos
			obj3.SetRect(pixel.R(0, 0, 16, 16))
			obj3.Flip = flip
			e3 := myecs.Manager.NewEntity()
			a3 := reanimator.NewBatchAnimation(thirdGoop, img.Batchers[constants.TileBatch], thirdGoop, reanimator.Tran)
			a3.SetEndTrigger(func() {
				CreateFallingGoop(thirdTile, flip)
				myecs.Manager.DisposeEntity(e3)
			})
			anim3 := reanimator.NewSimple(a3)
			e3.AddComponent(myecs.Object, obj3)
			e3.AddComponent(myecs.Drawable, anim3)
			e3.AddComponent(myecs.Animated, anim3)
			e3.AddComponent(myecs.Temp, myecs.ClearFlag(false))
		}
	}
}

func CreateFallingGoop(tile *data.Tile, flip bool) {
	obj := object.New()
	obj.Layer = 34
	obj.Pos = tile.Object.Pos
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Flip = flip
	d := data.NewDynamic()
	d.LastTile = tile
	e := myecs.Manager.NewEntity()
	fall := reanimator.NewBatchAnimation("goop_fall", img.Batchers[constants.TileBatch], "goop_fall", reanimator.Loop)
	land := reanimator.NewBatchAnimation("goop_land", img.Batchers[constants.TileBatch], "goop_land", reanimator.Tran)
	land.Offset.Y--
	land.SetTrigger(3, func() {
		x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
		lt := data.CurrLevel.Get(x, y-1)
		if lt != nil {
			switch lt.Block {
			case data.BlockTurf, data.BlockBedrock, data.BlockClose,
				data.BlockCracked, data.BlockFall, data.BlockPhase, data.BlockSpike:
				lt.Block = data.BlockGoop
				lt.Flags.Collapse = false
				lt.Counter = 0
				lt.Flags.Cracked = false
			}
		}
	})
	land.SetEndTrigger(func() {
		myecs.Manager.DisposeEntity(e)
	})
	anim := reanimator.New().
		AddAnimation(fall).
		AddAnimation(land).
		AddNull("none").
		SetChooseFn(func() string {
			if d.Flags.Floor {
				return "goop_land"
			} else {
				return "goop_fall"
			}
		}).SetDefault("flamethrower").Finish()
	anim.Update()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Drawable, anim)
	e.AddComponent(myecs.Animated, anim)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	d.Object = obj
	d.Entity = e
	d.Flags.NoLadders = true
	e.AddComponent(myecs.Dynamic, d)
	e.AddComponent(myecs.Update, data.NewFn(func() {
		x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
		lt := data.CurrLevel.Get(x, y)
		if lt != nil && lt.Block == data.BlockLiquid {
			myecs.Manager.DisposeEntity(e)
		}
	}))
}

func CreateTransporter(pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 8))
	obj.Layer = 9
	e := myecs.Manager.NewEntity()
	trans := &data.Transporter{
		Item: &data.BasicItem{
			Name:     "Transporter",
			Key:      key,
			Object:   obj,
			Sprite:   img.NewSprite(key, constants.TileBatch),
			Entity:   e,
			Color:    metadata.Color,
			Metadata: metadata,
			Origin:   coords,
			Layer:    9,
		},
	}
	trans.Item.Metadata = data.CopyMetadata(metadata)

	if len(metadata.LinkedTiles) > 0 {
		trans.Dest = data.CurrLevel.Get(metadata.LinkedTiles[0].X, metadata.LinkedTiles[0].Y)
	}
	if trans.Dest != nil && trans.Dest.Block != data.BlockTransporterExit {
		trans.Dest = nil
	}
	trans.BarO = object.New()
	trans.BarO.Pos = pos
	trans.BarO.Layer = 16
	trans.BarUp = false
	barDown := reanimator.NewBatchAnimationCustom("bar_down", img.Batchers[constants.TileBatch], "transporter_bar", []int{0, 1, 2, 3, 4, 5, 6}, reanimator.Hold)
	barUp := reanimator.NewBatchAnimationCustom("bar_up", img.Batchers[constants.TileBatch], "transporter_bar", []int{6, 5, 4, 3, 2, 1}, reanimator.Tran)
	barUp.SetEndTrigger(func() {
		trans.Item.Using = false
		trans.BarUp = false
	})
	trans.BarT = reanimator.New().
		AddNull("none").
		AddAnimation(barDown).
		AddAnimation(barUp).
		SetChooseFn(func() string {
			if trans.Item.Using {
				if trans.BarUp {
					return "bar_up"
				} else {
					return "bar_down"
				}
			}
			trans.BarUp = false
			return "none"
		}).SetDefault("none").Finish()
	trans.BarE = myecs.Manager.NewEntity()
	trans.BarE.AddComponent(myecs.Object, trans.BarO)
	trans.BarE.AddComponent(myecs.Drawable, trans.BarT)
	trans.BarE.AddComponent(myecs.Animated, trans.BarT)
	trans.BarE.AddComponent(myecs.Temp, myecs.ClearFlag(false))

	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, trans.Item.Sprite)
	e.AddComponent(myecs.Item, trans.Item)
	if trans.Dest != nil {
		trans.Item.Action = EnterTransporter(trans)
		e.AddComponent(myecs.OnTouch, trans.Item.Action)
	}
	e.AddComponent(myecs.LvlElement, trans)
	return trans.Item
}

func EnterTransporter(trans *data.Transporter) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if ch.Layer == 0 {
			return
		}
		// return if being used
		if trans.Item.Using {
			return
		}
		// return if player doesn't match the color of transporter
		switch trans.Item.Color {
		case data.NonPlayerRed:
			if p >= 0 && p < constants.MaxPlayers {
				return
			}
		case data.PlayerBlue, data.PlayerGreen, data.PlayerPurple, data.PlayerOrange:
			if ch.Color != trans.Item.Color && p >= 0 && p < constants.MaxPlayers {
				return
			}
		}
		// return if no linked tiles or if linked tile is empty
		if trans.Dest == nil || trans.Dest.Block != data.BlockTransporterExit {
			return
		}
		trans.Item.Using = true
		trans.BarUp = false
		ch.State = data.DoingAction
		ch.Flags.ItemAction = data.TransportIn
		ch.Flags.Transport = false
		ch.Object.Pos = trans.Item.Object.Pos
		ch.Object.Layer = 15
		sfx.SoundPlayer.PlaySound(constants.SFXTransIn, -2.)
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if ch.State == data.Dead {
				trans.Item.Using = false
				entity.RemoveComponent(myecs.Update)
				ch.Object.Layer = ch.Layer
				return
			}
			if ch.Flags.ItemAction == data.TransportIn &&
				ch.Flags.Transport && !trans.Dest.Flags.Using {
				// teleport the player
				sfx.SoundPlayer.PlaySound(constants.SFXTransOut, -1.)
				ch.Flags.ItemAction = data.TransportExit
				ch.Flags.Transport = false
				ch.Object.Pos = trans.Dest.Object.Pos
				ch.Object.PostPos = ch.Object.Pos
				trans.Dest.Flags.Using = true
				trans.BarUp = true
				entity.RemoveComponent(myecs.Update)
			}
		}))
	})
}

func CreateTransporterExit(pos pixel.Vec, tile *data.Tile) {
	zap := reanimator.NewBatchAnimationCustom("zap", img.Batchers[constants.TileBatch], "transporter_exit_zap", []int{0, 1, 0, 1, 0, 1, 0, 1}, reanimator.Tran)
	zap.SetEndTrigger(func() {
		tile.Flags.Using = false
	})
	tree := reanimator.New().
		AddNull("none").
		AddAnimation(zap).
		SetChooseFn(func() string {
			if tile.Flags.Using {
				return "zap"
			}
			return "none"
		}).SetDefault("none").Finish()
	e := myecs.Manager.NewEntity()
	obj := object.New()
	obj.Pos = pos
	obj.Layer = 16
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Drawable, tree)
	e.AddComponent(myecs.Animated, tree)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
}
