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
	shimmer := reanimator.NewBatchAnimationCustom("shimmer", batch, fmt.Sprintf("gem_%s_shimmer", color), []int{0, 0, 1, 1, 2, 2}, reanimator.Tran)
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
		AddComponent(myecs.OnTouch, data.NewInteract(CollectGem)).
		AddComponent(myecs.LvlElement, struct{}{})
}

func CollectGem(p int, ch *data.Dynamic, entity *ecs.Entity) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	data.CurrLevelSess.PlayerStats[p].LScore += 1
	data.CurrLevelSess.PlayerStats[p].LGems++
	sfx.SoundPlayer.PlaySound(constants.SFXGem, -2.)
	myecs.Manager.DisposeEntity(entity)
}

func CreateDoor(pos pixel.Vec, key string) {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 8))
	obj.Layer = 9
	door := &data.Door{
		Object: obj,
	}
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	openKey := key
	switch key {
	case constants.TileDoorYellow, constants.TileDoorOrange,
		constants.TileDoorGray, constants.TileDoorCyan,
		constants.TileDoorBlue, constants.TileDoorGreen,
		constants.TileDoorPurple, constants.TileDoorBrown:
		door.DoorType = data.Opened
		door.Color = strings.Replace(key, "door_", "", -1)
	case constants.TileClosedYellow, constants.TileClosedOrange,
		constants.TileClosedGray, constants.TileClosedCyan,
		constants.TileClosedBlue, constants.TileClosedGreen,
		constants.TileClosedPurple, constants.TileClosedBrown:
		door.DoorType = data.Opened
		openKey = strings.Replace(key, "_open", "", -1)
		door.Color = strings.Replace(openKey, "door_", "", -1)
		e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
	case constants.TileLockYellow, constants.TileLockOrange,
		constants.TileLockGray, constants.TileLockCyan,
		constants.TileLockBlue, constants.TileLockGreen,
		constants.TileLockPurple, constants.TileLockBrown:
		door.DoorType = data.Locked
		door.Color = strings.Replace(key, "lock_", "", -1)
		e.AddComponent(myecs.Drawable, img.NewSprite(key, constants.TileBatch))
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
			e.AddComponent(myecs.Drawable, img.NewSprite(openKey, constants.TileBatch))
			e.AddComponent(myecs.OnTouch, interaction)
		})
	}
	e.AddComponent(myecs.Door, door)
	e.AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
		switch door.DoorType {
		case data.Opened:
			if data.CurrLevel.DoorsOpen && reanimator.FrameSwitch {
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

func EnterDoor() *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 12
		sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
		data.CurrLevel.Complete = true
	})
}

func EnterPlayerDoor(color string) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if p < 0 || p >= constants.MaxPlayers || ch.Color != color {
			return
		}
		data.CurrLevelSess.PlayerStats[p].LScore += 12
		sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
		data.CurrLevel.Complete = true
	})
}

func CreateBox(pos pixel.Vec, tile *data.Tile) {
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
	e.AddComponent(myecs.PickUp, data.NewPickUp("Box", 10))
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
	case data.OnLadder, data.OnBar, data.Flying:
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

func CreateKey(pos pixel.Vec, key string) {
	obj := object.New().WithID(key)
	obj.Pos = pos
	obj.SetRect(pixel.R(0, 0, 8, 14))
	obj.Layer = 14
	color := strings.Replace(key, "key_", "", -1)
	theKey := &data.Key{
		Object: obj,
		Sprite: img.NewSprite(key, constants.TileBatch),
		PickUp: data.NewPickUp(fmt.Sprintf("Key (%s)", color), 5),
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
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		if KeyUnlock(ch.Object.Pos, color) {
			DropItem(ch)
			myecs.Manager.DisposeEntity(entity)
		}
	})
}

func KeyUnlock(chPos pixel.Vec, color string) bool {
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

func CreateJetpack(pos pixel.Vec, metadata data.TileMetadata, origin world.Coords) {
	key := constants.ItemJetpack
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	jetpack := &data.Jetpack{
		Object:   obj,
		PickUp:   data.NewPickUp("Jetpack", 5),
		Entity:   e,
		Metadata: metadata,
		Origin:   origin,
	}
	jetpack.Action = JetpackAction(jetpack)
	jpa := reanimator.NewBatchAnimation("jetpack", img.Batchers[constants.TileBatch], "jetpack", reanimator.Loop)
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
						data.NewFloatingText(fmt.Sprintf("%d", timer), true, false, false, txPos, pixel.ToRGBA(constants.ColorWhite), pixel.ToRGBA(pColor), nil)
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

func CreateDisguise(pos pixel.Vec, metadata data.TileMetadata, origin world.Coords) {
	key := constants.ItemDisguise
	obj := object.New().WithID(key).SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	obj.Layer = 12
	e := myecs.Manager.NewEntity()
	disguise := &data.Disguise{
		Object:   obj,
		PickUp:   data.NewPickUp("Jetpack", 5),
		Entity:   e,
		Metadata: metadata,
		Origin:   origin,
	}
	disguise.Action = DisguiseAction(disguise)
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
	disguise.Anim = reanimator.New(reanimator.NewSwitch().
		AddAnimation(regenA).
		AddAnimation(reanimator.NewBatchSprite("disguise", img.Batchers[constants.TileBatch], "disguise", reanimator.Hold)).
		AddNull("none").
		SetChooseFn(func() string {
			if disguise.Waiting || disguise.Using {
				return "none"
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

func DisguiseAction(disguise *data.Disguise) *data.Interact {
	return data.NewInteract(func(p int, ch *data.Dynamic, entity *ecs.Entity) {
		DropItem(ch)
		// change the player to disguised
		ch.Flags.Disguised = true
		ch.Anims[1].SetAnim(ch.Anims[0].GetCurrentAnim().Key, ch.Anims[0].GetCurrentFrame())
		// set the disguise vars
		disguise.Using = true
		disguise.Counter = 0
		// remove the pickup component
		disguise.Entity.RemoveComponent(myecs.PickUp)
		entity.AddComponent(myecs.Update, data.NewFn(func() {
			if disguise.Using {
				if data.CurrLevel.FrameChange {
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
						data.NewFloatingText(fmt.Sprintf("%d", timer), true, false, false, txPos, pixel.ToRGBA(constants.ColorWhite), pixel.ToRGBA(pColor), nil)
					}
					disguise.Counter++
					if disguise.Counter > disguise.Metadata.Timer*2 {
						ch.Flags.Disguised = false
						if disguise.Metadata.Regenerate {
							disguise.Using = false
							disguise.Waiting = true
							disguise.Counter = 0
						} else {
							myecs.Manager.DisposeEntity(disguise.Entity)
						}
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
