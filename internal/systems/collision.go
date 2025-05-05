package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/util"
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
			pTile := data.CurrLevel.Get(px, py)
			pLeft := data.CurrLevel.Get(px-1, py)
			pRight := data.CurrLevel.Get(px+1, py)
			pUp := data.CurrLevel.Get(px, py+1)
			pDown := data.CurrLevel.Get(px, py-1)
			var enemyL, enemyR, enemyU, enemyD bool
			enemyL, enemyR, enemyU, enemyD = enemyCollision(ch, px, py)
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
					enemyL, enemyR, enemyU, enemyD = enemyCollision(ch, x, y)
					tile := data.CurrLevel.Get(x, y)
					left := data.CurrLevel.Get(x-1, y)
					right := data.CurrLevel.Get(x+1, y)
					up := data.CurrLevel.Get(x, y+1)
					down := data.CurrLevel.Get(x, y-1)
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
	v = v.Add(pixel.V(world.HalfSize, world.HalfSize))
	ch.Object.Pos.Y = v.Y
}

func wallCollisions(ch *data.Dynamic, tile, left, right *data.Tile, enemyLeft, enemyRight bool, chPos pixel.Vec) {
	// for left and right, we stop the character if the next tile is solid and either
	//   if they run into the tile
	//   or if they are on a ladder (so they stay in the center of the ladder)
	if ch.State == data.OnLadder || ch.State == data.Falling {
		ch.Flags.LeftWall = left.IsSolid() || left.Block == data.BlockLiquid || enemyLeft
		ch.Flags.RightWall = right.IsSolid() || right.Block == data.BlockLiquid || enemyRight
	} else {
		if chPos.X-ch.Object.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
			if left.IsSolid() || left.Block == data.BlockLiquid {
				ch.Flags.LeftWall = true
				ch.Object.Pos.X = tile.Object.Pos.X - world.HalfSize + ch.Object.HalfWidth
			} else if enemyLeft && ch.Object.Pos.X < ch.Object.LastPos.X {
				ch.Flags.LeftWall = true
				ch.Object.Pos.X = ch.Object.LastPos.X
			}
		}
		if chPos.X+ch.Object.HalfWidth >= tile.Object.Pos.X+world.HalfSize {
			if right.IsSolid() || right.Block == data.BlockLiquid {
				ch.Flags.RightWall = true
				ch.Object.Pos.X = tile.Object.Pos.X + world.HalfSize - ch.Object.HalfWidth
			} else if enemyRight && ch.Object.Pos.X > ch.Object.LastPos.X {
				ch.Flags.RightWall = true
				ch.Object.Pos.X = ch.Object.LastPos.X
			}
		}
	}
}

func ceilingCollisions(ch *data.Dynamic, tile, up *data.Tile, enemyUp bool, chPos pixel.Vec) {
	// for up, we just make sure they don't enter a solid tile from below
	if ch.Flags.Thrown || ch.Flags.LongJump || ch.Flags.ItemAction == data.Drilling {
		return
	}
	if chPos.Y+ch.Object.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
		if up.IsSolid() || up.Block == data.BlockLiquid {
			ch.Flags.Ceiling = true
			ch.Object.Pos.Y = tile.Object.Pos.Y + world.HalfSize - ch.Object.HalfHeight
		} else if enemyUp && ch.Object.Pos.Y > ch.Object.LastPos.Y {
			ch.Flags.Ceiling = true
		}
	}
}

func floorCollisions(ch *data.Dynamic, tile, down *data.Tile, enemyDown bool, chPos pixel.Vec) {
	standOn, _ := standOnSystem(ch.Object.ID, down)
	standOnBelow := !ch.Actions.Down() && ch.State != data.OnLadder && standOn
	touchingFloor := chPos.Y-ch.Object.HalfHeight <= tile.Object.Pos.Y-world.HalfSize && !ch.Flags.HighJump && !ch.Flags.LongJump
	if ch.Flags.NoLadders {
		if (standOnBelow || down.IsSolid() || down.IsRunnable()) && touchingFloor {
			ch.Flags.Floor = true
		}
	} else {
		if enemyDown && ch.State == data.Falling && touchingFloor {
			ch.Object.Pos.Y = tile.Object.Pos.Y - world.HalfSize + ch.Object.HalfHeight
		} else if enemyDown && ch.State == data.OnLadder {
			if touchingFloor {
				ch.Flags.Climbed = false
				ch.Object.Pos.Y = tile.Object.Pos.Y - world.HalfSize + ch.Object.HalfHeight
			} else {
				ch.Flags.Floor = true
			}
		} else if ((down.IsSolid() ||
			(down.IsLadder() && (ch.State != data.Flying && ch.State != data.OnLadder && ch.State != data.Leaping)) ||
			standOnBelow) &&
			touchingFloor) ||
			(down != nil && down.IsLadder() && !tile.IsLadder() && ch.State == data.OnLadder && !touchingFloor && ch.Actions.Up()) {
			ch.Flags.Floor = true
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

func enemyCollision(ch *data.Dynamic, x, y int) (bool, bool, bool, bool) {
	if ch.Enemy < 0 {
		return false, false, false, false
	}
	var enemyL, enemyR, enemyU, enemyD bool
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch2, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC && !obj.Hidden &&
			ch2.State != data.Hit && ch2.State != data.Dead &&
			ch2.Enemy != ch.Enemy && ch2.Enemy > -1 {
			pPos := ch2.Object.LastPos
			px, py := world.WorldToMap(pPos.X, pPos.Y)
			yd := util.Abs(py - y)
			xd := util.Abs(px - x)
			if xd > 1 || yd > 1 {
				continue // too far away to matter
			}
			overlap := ch.Object.Rect.Moved(ch.Object.Pos).Intersects(ch2.Object.Rect.Moved(ch2.Object.Pos))
			if px == x-1 && (py == y || (yd == 1 && overlap)) {
				enemyL = true
			}
			if px == x+1 && (py == y || (yd == 1 && overlap)) {
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
