package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
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
	} else {
		ch.Inventory = PickUpItem(ch, p)
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
			//// check if item is held by someone else, and drop it for them
			//if pickUp.Held > -1 {
			//	DropLift(data.CurrLevel.Players[pickUp.Held], false)
			//}
			if cycle {
				pickUp.Cycle[p]++
			} else {
				pickUp.Cycle[p] = 0
			}
			pickUp.Inventory = p
			heldEntity.obj.Hidden = true
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
	}
	ch.Inventory = item
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

func DoAction(ch *data.Dynamic) {
	if ch.Player > -1 && ch.Player < constants.MaxPlayers &&
		ch.Inventory != nil && ch.Inventory.HasComponent(myecs.Action) {
		if fnA, ok := ch.Inventory.GetComponentData(myecs.Action); ok {
			if colFn, okC := fnA.(*data.Interact); okC {
				colFn.Fn(data.CurrLevel, int(ch.Player), ch, ch.Inventory)
			}
		}
	}
}
