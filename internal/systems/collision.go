package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func CollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC && !obj.Hidden {
			setCollisionFlags(ch)
			chPos := ch.Object.Pos
			x, y := world.WorldToMap(chPos.X, chPos.Y)
			tile := data.CurrLevel.Tiles.Get(x, y)
			left := data.CurrLevel.Tiles.Get(x-1, y)
			right := data.CurrLevel.Tiles.Get(x+1, y)
			up := data.CurrLevel.Tiles.Get(x, y+1)
			down := data.CurrLevel.Tiles.Get(x, y-1)
			if tile == nil {
				outsideOfMap(ch, x, y)
			} else {
				// check each direction for collision
				wallCollisions(ch, tile, left, right, chPos)
				ceilingCollisions(ch, tile, up, chPos)
				floorCollisions(ch, tile, down, chPos)
				// check if the character can use a ladder
				//ch.Flags.LadderHere = tile.Ladder
				//ch.Flags.LadderDown = down != nil && down.Ladder
				//// check if the character is on a ladder
				////ch.Flags.OnLadder = (ch.Flags.OnLadder && !ch.Flags.Floor) || (ch.Flags.LadderHere && !ch.Flags.Floor)
				//ch.Flags.OnLadder = ch.Flags.OnLadder && !ch.Flags.Floor && (ch.Flags.LadderHere || ch.Flags.LadderDown)
				// check to see if the character can run along the floor
				// They can run if they are touching solid ground (Floor)
				//   or are at the top of a ladder or the ladder is on top of Turf
				//   or are holding a ladder and there is floor directly below
				//ch.Flags.CanRun = ((down == nil || down.Solid() || standOnBelow) && touchingFloor) ||
				//	(ch.Flags.LadderDown && touchingFloor && (!ch.Flags.LadderHere || down.Block == data.BlockTurf)) ||
				//	((down == nil || down.Solid() || down.Block == data.BlockTurf) && ch.Flags.OnLadder)
			}
			//debug.AddTruthText("LeftWall:  ", ch.Flags.LeftWall)
			//debug.AddTruthText("RightWall: ", ch.Flags.RightWall)
			//debug.AddTruthText("Ceiling:   ", ch.Flags.Ceiling)
			//debug.AddTruthText("Floor:     ", ch.Flags.Floor)
			//debug.AddTruthText("On Ladder: ", ch.Flags.OnLadder)
			//debug.AddTruthText("Ladder Here: ", ch.Flags.LadderHere)
			//debug.AddTruthText("Ladder Down: ", ch.Flags.LadderDown)
		}
	}
}

func setCollisionFlags(ch *data.Dynamic) {
	ch.Flags.LeftWall = false
	ch.Flags.RightWall = false
	ch.Flags.Ceiling = false
	ch.Flags.Floor = false
}

func outsideOfMap(ch *data.Dynamic, x, y int) {
	// uh oh
	ch.Flags.Ceiling = true
	ch.Flags.Floor = true
	v := world.MapToWorld(world.Coords{X: x, Y: y})
	v = v.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
	ch.Object.Pos.Y = v.Y
}

func wallCollisions(ch *data.Dynamic, tile, left, right *data.Tile, chPos pixel.Vec) {
	// for left and right, we stop the character if the next tile is solid and either
	//   if they run into the tile
	//   or if they are on a ladder (so they stay in the center of the ladder)
	if left == nil || left.SolidH() {
		if ch.State == data.OnLadder || ch.State == data.Falling {
			ch.Flags.LeftWall = true
		} else if chPos.X-ch.Object.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
			ch.Flags.LeftWall = true
			ch.Object.Pos.X = tile.Object.Pos.X - world.HalfSize + ch.Object.HalfWidth
		}
	}
	if right == nil || right.SolidH() {
		if ch.State == data.OnLadder || ch.State == data.Falling {
			ch.Flags.RightWall = true
		} else if chPos.X+ch.Object.HalfWidth >= tile.Object.Pos.X+world.HalfSize {
			ch.Flags.RightWall = true
			ch.Object.Pos.X = tile.Object.Pos.X + world.HalfSize - ch.Object.HalfWidth
		}
	}
}

func ceilingCollisions(ch *data.Dynamic, tile, up *data.Tile, chPos pixel.Vec) {
	// for up, we just make sure they don't enter a solid tile from below
	if (up == nil || up.SolidV()) &&
		chPos.Y+ch.Object.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
		ch.Flags.Ceiling = true
		ch.Object.Pos.Y = tile.Object.Pos.Y + world.HalfSize - ch.Object.HalfHeight
	}
}

func floorCollisions(ch *data.Dynamic, tile, down *data.Tile, chPos pixel.Vec) {
	standOnBelow := !ch.Actions.Down() && ch.State != data.OnLadder && standOnSystem(down)
	touchingFloor := chPos.Y-ch.Object.HalfHeight <= tile.Object.Pos.Y-world.HalfSize && !ch.Flags.HighJump && !ch.Flags.LongJump
	if ((down == nil ||
		down.SolidV() ||
		(down.IsLadder() && !tile.IsLadder() && ch.State != data.OnLadder && !ch.Actions.Down()) ||
		standOnBelow) && touchingFloor) ||
		(down != nil && down.IsLadder() && !tile.IsLadder() && ch.State == data.OnLadder && !touchingFloor) {
		ch.Flags.Floor = true
		ch.Object.Pos.Y = tile.Object.Pos.Y - world.HalfSize + ch.Object.HalfHeight
	}
}

func standOnSystem(downTile *data.Tile) bool {
	if downTile == nil {
		return false
	}
	for _, result := range myecs.Manager.Query(myecs.IsStandOn) {
		if obj, okO := result.Components[myecs.Object].(*object.Object); okO {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			standCoords := world.Coords{X: x, Y: y}
			if downTile.Coords == standCoords {
				return true
			}
		}
	}
	return false
}
