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
		if okO && okC && !obj.Hidden &&
			ch.State != data.Hit && ch.State != data.Dead {
			setCollisionFlags(ch)
			chPos := ch.Object.Pos
			pPos := ch.Object.LastPos
			px, py := world.WorldToMap(pPos.X, pPos.Y)
			pTile := data.CurrLevel.Tiles.Get(px, py)
			pLeft := data.CurrLevel.Tiles.Get(px-1, py)
			pRight := data.CurrLevel.Tiles.Get(px+1, py)
			pUp := data.CurrLevel.Tiles.Get(px, py+1)
			pDown := data.CurrLevel.Tiles.Get(px, py-1)
			var enemyL, enemyR, enemyU, enemyD bool
			enemyL, enemyR, enemyU, enemyD = enemyCollision(ch.Enemy, px, py)
			if pTile == nil {
				outsideOfMap(ch, px, py)
			} else { // check each direction for collision
				wallCollisions(ch, pTile, pLeft, pRight, enemyL, enemyR, chPos)
				ceilingCollisions(ch, pTile, pUp, enemyU, chPos)
				floorCollisions(ch, pTile, pDown, enemyD, chPos)
				chPos = ch.Object.Pos
				x, y := world.WorldToMap(chPos.X, chPos.Y)
				if x != px || y != py { // check again if they changed tiles
					setCollisionFlags(ch)
					enemyL, enemyR, enemyU, enemyD = enemyCollision(ch.Enemy, x, y)
					tile := data.CurrLevel.Tiles.Get(x, y)
					left := data.CurrLevel.Tiles.Get(x-1, y)
					right := data.CurrLevel.Tiles.Get(x+1, y)
					up := data.CurrLevel.Tiles.Get(x, y+1)
					down := data.CurrLevel.Tiles.Get(x, y-1)
					wallCollisions(ch, tile, left, right, enemyL, enemyR, chPos)
					ceilingCollisions(ch, tile, up, enemyU, chPos)
					floorCollisions(ch, tile, down, enemyD, chPos)
				}
			}
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

func wallCollisions(ch *data.Dynamic, tile, left, right *data.Tile, enemyLeft, enemyRight bool, chPos pixel.Vec) {
	// for left and right, we stop the character if the next tile is solid and either
	//   if they run into the tile
	//   or if they are on a ladder (so they stay in the center of the ladder)
	if left == nil || left.IsSolid() || enemyLeft {
		if ch.State == data.OnLadder || ch.State == data.Falling {
			ch.Flags.LeftWall = true
		} else if chPos.X-ch.Object.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
			ch.Flags.LeftWall = true
			ch.Object.Pos.X = tile.Object.Pos.X - world.HalfSize + ch.Object.HalfWidth
		}
	}
	if right == nil || right.IsSolid() || enemyRight {
		if ch.State == data.OnLadder || ch.State == data.Falling {
			ch.Flags.RightWall = true
		} else if chPos.X+ch.Object.HalfWidth >= tile.Object.Pos.X+world.HalfSize {
			ch.Flags.RightWall = true
			ch.Object.Pos.X = tile.Object.Pos.X + world.HalfSize - ch.Object.HalfWidth
		}
	}
}

func ceilingCollisions(ch *data.Dynamic, tile, up *data.Tile, enemyUp bool, chPos pixel.Vec) {
	// for up, we just make sure they don't enter a solid tile from below
	if ch.Flags.Thrown {
		return
	}
	if (up == nil || up.IsSolid() || enemyUp) &&
		chPos.Y+ch.Object.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
		ch.Flags.Ceiling = true
		ch.Object.Pos.Y = tile.Object.Pos.Y + world.HalfSize - ch.Object.HalfHeight
	}
}

func floorCollisions(ch *data.Dynamic, tile, down *data.Tile, enemyDown bool, chPos pixel.Vec) {
	standOn, _ := standOnSystem(ch.Object.ID, down)
	standOnBelow := !ch.Actions.Down() && ch.State != data.OnLadder && standOn
	touchingFloor := chPos.Y-ch.Object.HalfHeight <= tile.Object.Pos.Y-world.HalfSize && !ch.Flags.HighJump && !ch.Flags.LongJump
	if ch.Flags.NoLadders {
		if (standOnBelow || down.IsNilOrSolid() || down.IsRunnable()) && touchingFloor {
			ch.Flags.Floor = true
		}
	} else {
		if ((down == nil ||
			down.IsSolid() ||
			enemyDown ||
			((down.IsLadder() || down.IsNilOrSolid()) && (ch.State != data.OnLadder && ch.State != data.Leaping)) ||
			standOnBelow) &&
			touchingFloor) ||
			(down != nil && down.IsLadder() && !tile.IsLadder() && ch.State == data.OnLadder && !touchingFloor && ch.Actions.Up()) {
			if !enemyDown {
				ch.Flags.Floor = true
			}
			ch.Object.Pos.Y = tile.Object.Pos.Y - world.HalfSize + ch.Object.HalfHeight
		}
	}
}

func standOnSystem(id string, downTile *data.Tile) (bool, bool) {
	if downTile == nil {
		return false, false
	}
	for _, result := range myecs.Manager.Query(myecs.IsStandOn) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC &&
			!obj.Hidden && obj.ID != id {
			pos := obj.PostPos
			// adjustment, maybe add to constants
			if !d.Flags.Floor {
				pos.Y -= world.HalfSize * 0.5
			}
			x, y := world.WorldToMap(pos.X, pos.Y)
			standCoords := world.Coords{X: x, Y: y}
			if downTile.Coords == standCoords {
				return true, !d.Flags.Floor
			}
		}
	}
	return false, false
}

func enemyCollision(eId, x, y int) (bool, bool, bool, bool) {
	if eId < 0 {
		return false, false, false, false
	}
	var enemyL, enemyR, enemyU, enemyD bool
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC && !obj.Hidden &&
			ch.State != data.Hit && ch.State != data.Dead &&
			ch.Enemy != eId && ch.Enemy > -1 {
			pPos := ch.Object.LastPos
			px, py := world.WorldToMap(pPos.X, pPos.Y)
			if px == x-1 && py == y {
				enemyL = true
			}
			if px == x+1 && py == y {
				enemyR = true
			}
			if px == x && py == y+1 {
				enemyU = true
			}
			if px == x && py == y-1 {
				enemyD = true
			}
		}
	}
	return enemyL, enemyR, enemyU, enemyD
}
