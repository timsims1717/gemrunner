package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
)

// PickUpOrDropItem returns whether an item was dropped or picked up
func PickUpOrDropItem(ch *data.Dynamic, p int) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	if ch.Inventory != nil {
		PickUpAndDropItem(ch, p)
		ch.Flags.PickUpBuff = 0
	} else {
		item := PickUpItem(ch, p)
		if item != nil {
			ch.Inventory = item
			sfx.SoundPlayer.PlaySound(constants.SFXItem, 0.)
		}
	}
}

type pickUpCheck struct {
	priority  int
	cycle     int
	inventory bool
	entity    *ecs.Entity
	obj       *object.Object
}

// PickUpItem returns whether an item was picked up
func PickUpItem(ch *data.Dynamic, p int) *ecs.Entity {
	cx, cy := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
	chCoords := world.Coords{X: cx, Y: cy}
	var heldEntity pickUpCheck
	var sameSpace, cycle bool
	for _, result := range myecs.Manager.Query(myecs.IsPickUp) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		pu, okP := result.Components[myecs.PickUp].(*data.PickUp)
		if okO && okP && !obj.Hidden &&
			obj.ID != ch.Object.ID &&
			!pu.NoInventory &&
			pu.Inventory == -1 {
			x, y := world.WorldToMap(obj.Pos.X+obj.Offset.X, obj.Pos.Y+obj.Offset.Y)
			pickUpCoords := world.Coords{X: x, Y: y}
			if chCoords == pickUpCoords &&
				(heldEntity.entity == nil || !sameSpace ||
					pu.Cycle[p] < heldEntity.cycle ||
					pu.Priority < heldEntity.priority) {
				if sameSpace {
					cycle = true
				} else {
					cycle = false
				}
				sameSpace = true
				heldEntity = pickUpCheck{
					priority: pu.Priority,
					cycle:    pu.Cycle[p],
					entity:   result.Entity,
					obj:      obj,
				}
			}
		}
	}
	if heldEntity.entity != nil {
		if pu, ok := heldEntity.entity.GetComponentData(myecs.PickUp); ok {
			pickUp := pu.(*data.PickUp)
			if cycle {
				pickUp.Cycle[p]++
			} else {
				pickUp.Cycle[p] = 0
			}
			pickUp.Inventory = p
			heldEntity.obj.Hidden = true
			ch.Flags.PickUpBuff = 0
			return heldEntity.entity
		}
	}
	return nil
}

// PickUpAndDropItem returns whether an item was dropped
func PickUpAndDropItem(ch *data.Dynamic, p int) {
	item := PickUpItem(ch, p)
	if ch.Inventory != nil {
		DropItem(ch)
		sfx.SoundPlayer.PlaySound(constants.SFXDrop, 1.)
	}
	if item != nil {
		ch.Inventory = item
		sfx.SoundPlayer.PlaySound(constants.SFXItem, 0.)
	}
}

// DropItem returns whether an item was dropped
func DropItem(ch *data.Dynamic) bool {
	if ch.Inventory == nil {
		return false
	}
	// set the object's new position
	if o, okO := ch.Inventory.GetComponentData(myecs.Object); okO {
		obj := o.(*object.Object)
		tile := data.CurrLevel.Tiles.Get(world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y))
		obj.Hidden = false
		obj.Pos = tile.Object.Pos
		obj.PostPos = tile.Object.Pos
	}
	// set the object's pickup data
	if p, okP := ch.Inventory.GetComponentData(myecs.PickUp); okP {
		pickUp := p.(*data.PickUp)
		pickUp.Inventory = -1
	}
	ch.Inventory = nil
	return true
}

func DoAction(ch *data.Dynamic) bool {
	if ch.Player > -1 && ch.Player < constants.MaxPlayers &&
		ch.Inventory != nil && ch.Inventory.HasComponent(myecs.Action) {
		if fnA, ok := ch.Inventory.GetComponentData(myecs.Action); ok {
			if colFn, okC := fnA.(*data.Interact); okC {
				colFn.Fn(data.CurrLevel, ch.Player, ch, ch.Inventory)
				ch.Flags.ActionBuff = 0
				return true
			}
		}
	}
	return false
}

func UpdateInventory(ch *data.Dynamic) {
	if ch.Inventory != nil {
		o, ok1 := ch.Inventory.GetComponentData(myecs.Object)
		if ok1 {
			if obj, ok := o.(*object.Object); ok {
				obj.SetPos(ch.Object.Pos)
			}
		}
	}
}

func Dig(ch *data.Dynamic, isLeft bool) bool {
	if ch.Player > -1 && ch.Player < constants.MaxPlayers {
		if ch.State == data.OnLadder ||
			ch.State == data.OnBar ||
			ch.State == data.Grounded ||
			ch.State == data.Flying {
			var sideTile, digTile *data.Tile
			x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			tile := data.CurrLevel.Tiles.Get(x, y)
			if tile != nil {
				if isLeft {
					sideTile = data.CurrLevel.Tiles.Get(x-1, y)
					digTile = data.CurrLevel.Tiles.Get(x-1, y-1)
				} else {
					sideTile = data.CurrLevel.Tiles.Get(x+1, y)
					digTile = data.CurrLevel.Tiles.Get(x+1, y-1)
				}
				if sideTile != nil && digTile != nil {
					if !sideTile.IsLadder() && !sideTile.IsSolid() && digTile.CanDig() {
						// set action
						ch.State = data.DoingAction
						ch.Flags.ItemAction = data.MagicDig
						ch.Object.Flip = isLeft
						// start digging the tile
						//things := ThingsOnTile(sideTile)
						//collapse := true
						//for _, thing := range things {
						//	d, okD := thing.GetComponentData(myecs.Dynamic)
						//	if okD {
						//		if chE, ok := d.(*data.Dynamic); ok {
						//			if chE.Enemy > -1 {
						//				collapse = false
						//			}
						//		}
						//	}
						//}
						digTile.Flags.Collapse = true
						digTile.Metadata.Regenerate = true
						AddMask(digTile, "dig_mask", isLeft, false)
						digTile.Counter = 0
						obj := object.New()
						obj.Pos = digTile.Object.Pos
						obj.Pos.Y += 4
						if isLeft {
							obj.Flip = true
							obj.Pos.X += 4
						} else {
							obj.Pos.X -= 4
						}
						obj.Layer = 33
						m := myecs.Manager.NewEntity()
						anim := reanimator.NewSimple(reanimator.NewBatchAnimation("digMagic", img.Batchers[constants.TileBatch], "dig_magic", reanimator.Done))
						m.AddComponent(myecs.Object, obj)
						m.AddComponent(myecs.Animated, anim)
						m.AddComponent(myecs.Drawable, anim)
						m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
						m.AddComponent(myecs.Update, data.NewFn(func() {
							if anim.Done || ch.State != data.DoingAction {
								myecs.Manager.DisposeEntity(m)
								ch.StoredBlocks = append(ch.StoredBlocks, digTile)
							}
						}))
						return true
					}
				}
			}
		}
	}
	return false
}

func Place(ch *data.Dynamic, isLeft bool) bool {
	if ch.Player > -1 && false &&
		ch.Player < constants.MaxPlayers &&
		len(ch.StoredBlocks) > 0 {
		if ch.State == data.OnLadder ||
			ch.State == data.Grounded ||
			ch.State == data.Flying {
			var sideTile, digTile *data.Tile
			x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
			tile := data.CurrLevel.Tiles.Get(x, y)
			if tile != nil {
				if isLeft {
					sideTile = data.CurrLevel.Tiles.Get(x-1, y)
					digTile = data.CurrLevel.Tiles.Get(x-1, y-1)
				} else {
					sideTile = data.CurrLevel.Tiles.Get(x+1, y)
					digTile = data.CurrLevel.Tiles.Get(x+1, y-1)
				}
				oldTile := ch.StoredBlocks[0]
				if oldTile != nil && sideTile != nil && digTile != nil {
					if !sideTile.IsLadder() && !sideTile.IsSolid() && digTile.IsEmpty() {
						// set action
						ch.State = data.DoingAction
						ch.Flags.ItemAction = data.MagicPlace
						ch.Object.Flip = isLeft

						// set tile
						oldTile.CopyInto(digTile)
						oldTile.ToEmpty()
						if len(ch.StoredBlocks) > 1 {
							ch.StoredBlocks = ch.StoredBlocks[1:]
						} else {
							ch.StoredBlocks = []*data.Tile{}
						}
						AddMask(digTile, "dig_mask", isLeft, true)

						// add reversed wand
						obj := object.New()
						obj.Pos = digTile.Object.Pos
						obj.Pos.Y += 4
						if isLeft {
							obj.Flip = true
							obj.Pos.X += 4
						} else {
							obj.Pos.X -= 4
						}
						obj.Layer = 33
						m := myecs.Manager.NewEntity()
						anim := reanimator.NewSimple(reanimator.NewBatchAnimation("digMagic", img.Batchers[constants.TileBatch], "dig_magic", reanimator.Done).Reverse())
						m.AddComponent(myecs.Object, obj)
						m.AddComponent(myecs.Animated, anim)
						m.AddComponent(myecs.Drawable, anim)
						m.AddComponent(myecs.Temp, myecs.ClearFlag(false))
						m.AddComponent(myecs.Update, data.NewFn(func() {
							if anim.Done || ch.State != data.DoingAction {
								myecs.Manager.DisposeEntity(m)
							}
						}))
						return true
					}
				}
			}
		}
	}
	return false
}
