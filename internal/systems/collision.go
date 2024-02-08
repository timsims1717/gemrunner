package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
)

func CollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Character].(*data.Character)
		if okO && okC && !obj.Hidden {
			ch.Flags.LeftWall = false
			ch.Flags.RightWall = false
			ch.Flags.Ceiling = false
			ch.Flags.Floor = false
			ch.Flags.CanRun = false
			currPos := ch.FauxObj.Pos
			x, y := world.WorldToMap(currPos.X, currPos.Y)
			currTile := data.CurrLevel.Tiles.Get(x, y)
			leftTile := data.CurrLevel.Tiles.Get(x-1, y)
			rightTile := data.CurrLevel.Tiles.Get(x+1, y)
			upTile := data.CurrLevel.Tiles.Get(x, y+1)
			downTile := data.CurrLevel.Tiles.Get(x, y-1)
			if currTile == nil {
				// uh oh
				ch.Flags.Ceiling = true
				ch.Flags.Floor = true
				ch.Flags.CanRun = true
				v := world.MapToWorld(world.Coords{X: x, Y: y})
				ch.FauxObj.Pos.Y = v.Y
			} else {
				// check each direction for collision
				if (leftTile == nil || leftTile.Solid()) &&
					(currPos.X-obj.HalfWidth <= currTile.Object.Pos.X-world.HalfSize || ch.Flags.OnLadder) {
					ch.Flags.LeftWall = true
					ch.FauxObj.Pos.X = currTile.Object.Pos.X - world.HalfSize + obj.HalfWidth
				}
				if (rightTile == nil || rightTile.Solid()) &&
					(currPos.X+obj.HalfWidth >= currTile.Object.Pos.X+world.HalfSize || ch.Flags.OnLadder) {
					ch.Flags.RightWall = true
					ch.FauxObj.Pos.X = currTile.Object.Pos.X + world.HalfSize - obj.HalfWidth
				}
				if (upTile == nil || upTile.Solid()) && currPos.Y+obj.HalfHeight >= currTile.Object.Pos.Y+world.HalfSize {
					ch.Flags.Ceiling = true
					ch.FauxObj.Pos.Y = currTile.Object.Pos.Y + world.HalfSize - obj.HalfHeight
				}
				touchingFloor := currPos.Y-obj.HalfHeight <= currTile.Object.Pos.Y-world.HalfSize
				if (downTile == nil || downTile.Solid()) && touchingFloor {
					ch.Flags.Floor = true
					ch.FauxObj.Pos.Y = currTile.Object.Pos.Y - world.HalfSize + obj.HalfHeight
				}
				// check if the character can use a ladder
				ch.Flags.LadderHere = currTile.Ladder
				ch.Flags.LadderDown = downTile != nil && downTile.Ladder
				// check if the character is on a ladder
				ch.Flags.OnLadder = ch.Flags.OnLadder && !ch.Flags.Floor && (ch.Flags.LadderHere || ch.Flags.LadderDown)
				// check to see if the character can run along the floor
				// They can run if they are touching solid ground (Floor)
				//   or are at the top of a ladder or the ladder is on top of Turf
				//   or are holding a ladder and there is floor directly below
				ch.Flags.CanRun = ch.Flags.Floor ||
					(ch.Flags.LadderDown && touchingFloor && (!ch.Flags.LadderHere || downTile.Block == data.Turf)) ||
					((downTile == nil || downTile.Solid() || downTile.Block == data.Turf) && ch.Flags.OnLadder)
			}
			debug.AddTruthText("LeftWall:  ", ch.Flags.LeftWall)
			debug.AddTruthText("RightWall: ", ch.Flags.RightWall)
			debug.AddTruthText("Ceiling:   ", ch.Flags.Ceiling)
			debug.AddTruthText("Floor:     ", ch.Flags.Floor)
			debug.AddTruthText("On Ladder: ", ch.Flags.OnLadder)
			debug.AddTruthText("Ladder Here: ", ch.Flags.LadderHere)
			debug.AddTruthText("Ladder Down: ", ch.Flags.LadderDown)
		}
	}
}
