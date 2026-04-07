package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/data/death"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/world"
	"math"
)

func CharacterStateSystem() {
	if reanimator.FrameSwitch && data.CurrLevel != nil {
		for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
			_, okO := result.Components[myecs.Object].(*object.Object)
			ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
			if okO && okC {
				currPos := ch.Object.Pos.Add(ch.Object.Offset)
				x, y := world.WorldToMap(currPos.X, currPos.Y)
				tile := data.CurrLevel.Get(x, y)
				below := data.CurrLevel.Get(x, y-1)
				oldState := ch.State
				switch ch.State {
				case data.Grounded:
					if ch.Flags.WallClimb {
						touchingFloor := ch.Flags.Floor
						switch ch.Flags.Orientation {
						case data.Up:
							below = data.CurrLevel.Get(x, y+1)
							touchingFloor = ch.Flags.Ceiling
						case data.Left:
							below = data.CurrLevel.Get(x-1, y)
							touchingFloor = ch.Flags.LeftWall
						case data.Right:
							below = data.CurrLevel.Get(x+1, y)
							touchingFloor = ch.Flags.RightWall
						}
						if below.IsEmpty() || !touchingFloor {
							if ch.BelowTile != nil && ch.BelowTile != below { // just went over the edge
								ch.State = data.AroundCorner
								ch.Flags.NoCollision = true
							} else { // the tile disappeared
								ch.State = data.Falling
								if tile != nil {
									ch.Object.Pos.X = tile.Object.Pos.X
								}
							}
						} else {
							ch.BelowTile = below
						}
					} else {
						if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else if tile != nil && tile.IsLadder() { // a ladder is here
							if ch.Actions.Direction == data.Up { // climbed the ladder
								ch.State = data.OnLadder
								ch.Object.Pos.X = tile.Object.Pos.X
							} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
								ch.State = data.OnLadder
								ch.Object.Pos.X = tile.Object.Pos.X
							} else if below != nil && !below.IsRunnable() && (ch.Actions.Left() || ch.Actions.Right()) { // leaping onto ladder
								ch.State = data.Leaping
								ch.Flags.LeapOn = true
								ch.NextTile = tile
							} else if !ch.Flags.Floor {
								ch.State = data.OnLadder
								ch.Object.Pos.X = tile.Object.Pos.X
							}
						} else if tile != nil && tile.Block == data.BlockHideout && ch.Actions.Direction == data.Up {
							ch.State = data.DoingAction
							ch.Flags.ItemAction = data.Hiding
							ch.Object.Pos.X = tile.Object.Pos.X
						} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
							ch.State = data.OnLadder
							ch.Object.Pos.X = tile.Object.Pos.X
						} else if !ch.Flags.Floor && (below != nil || data.CurrLevel.Continuity) && !below.IsLadder() {
							if ch.Enemy > -1 && below.Block == data.BlockTurf && below.Flags.Collapse {
								ch.State = data.Tripping
								if ch.Object.Flip {
									ch.Object.Pos.X -= 1
								} else {
									ch.Object.Pos.X += 1
								}
								ch.NextTile = below
								ch.NextTile.Flags.Occupied = true
							} else {
								ch.State = data.Falling
								if tile != nil {
									ch.Object.Pos.X = tile.Object.Pos.X
								}
							}
						}
					}
				case data.OnLadder:
					ch.Flags.Landing = false
					if ch.Flags.Floor && ch.Actions.Up() && !tile.IsLadder() { // just got to the top
						ch.State = data.Grounded
						if ch.Actions.Left() { // to the left
							ch.Object.Flip = true
						} else if ch.Actions.Right() { // to the right
							ch.Object.Flip = false
						}
					} else if tile.IsLadder() &&
						!(below.IsSolid() ||
							below.IsRunnable()) { // can only leap if the stuff below isn't solid
						if ch.Actions.Direction == data.Left &&
							!ch.Flags.LeftWall { // leaping to the left
							ch.State = data.Leaping
							ch.Object.Flip = true
							left := data.CurrLevel.Get(x-1, y)
							ch.NextTile = left
							if left != nil && left.IsLadder() { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						} else if ch.Actions.Direction == data.Right &&
							!ch.Flags.RightWall { // leaping to the right
							ch.State = data.Leaping
							ch.Object.Flip = false
							right := data.CurrLevel.Get(x+1, y)
							ch.NextTile = right
							if right != nil && right.IsLadder() { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						}
					} else if tile != nil && tile.Block == data.BlockBar && !ch.Actions.Down() {
						ch.State = data.OnBar
					} else if !tile.IsLadder() && below != nil && below.IsLadder() {
						if ch.Actions.Direction == data.Left { // run to the left
							ch.State = data.Grounded
							ch.Object.Flip = true
						} else if ch.Actions.Direction == data.Right { // run to the right
							ch.State = data.Grounded
							ch.Object.Flip = false
						} else {
							ch.State = data.OnLadder
						}
					} else if ch.Flags.Floor || below.IsRunnable() { // reached the bottom
						if ch.Actions.Direction == data.Left && !ch.Flags.LeftWall { // run to the left
							ch.State = data.Grounded
							ch.Object.Flip = true
						} else if ch.Actions.Direction == data.Right && !ch.Flags.RightWall { // run to the right
							ch.State = data.Grounded
							ch.Object.Flip = false
						} else if ch.Actions.Direction == data.Down {
							if !tile.IsLadder() {
								ch.State = data.Falling
							} else if ch.Flags.Floor {
								ch.State = data.Grounded
							}
						}
					} else if !tile.IsLadder() {
						ch.State = data.Falling
					}
				case data.OnBar:
					ch.Flags.Landing = false
					if tile == nil || tile.Block != data.BlockBar {
						if below == nil || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							if ch.Actions.Left() || ch.Actions.Right() { // leaping onto ladder
								ch.State = data.Leaping
								ch.Flags.LeapOn = true
								ch.NextTile = tile
							} else if !ch.Flags.Floor {
								ch.State = data.OnLadder
							}
						} else if ch.Flags.Floor {
							ch.State = data.Grounded
						} else {
							ch.State = data.Falling
						}
					} else if !ch.Flags.Floor && ch.Actions.Down() {
						ch.State = data.Falling
					} else if below != nil && below.IsLadder() && ch.Actions.Down() {
						ch.State = data.OnLadder
					}
				case data.Falling:
					if ch.Flags.Floor {
						ch.Flags.Orientation = data.Down
						ch.State = data.Grounded
						if ch.LastTile.Coords.Y > tile.Coords.Y+2 {
							ch.Flags.Landing = true
						}
					} else if tile != nil && tile.IsLadder() {
						ch.State = data.OnLadder
					} else if tile != nil &&
						tile.Block == data.BlockBar &&
						tile.Coords != ch.LastTile.Coords {
						ch.State = data.OnBar
					} else if tile != nil &&
						tile.Block == data.BlockTurf &&
						tile.Flags.Collapse && ch.Enemy > -1 {
						ch.State = data.InHole
						ch.ACounter = 0
						tile.Flags.Occupied = true
						ch.NextTile = tile
						ch.Object.Pos = tile.Object.Pos
						DropItem(ch)
					}
					if tile != nil && ch.Pushing == nil {
						ch.Object.Pos.X = tile.Object.Pos.X
					}
				case data.Jumping:
					if !(ch.Flags.HighJump || ch.Flags.LongJump) {
						if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
						ch.Object.Pos.X = tile.Object.Pos.X
					}
				case data.Leaping:
					if !(ch.Flags.LeapOff || ch.Flags.LeapOn || ch.Flags.LeapTo) {
						if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
						ch.Object.Pos.X = tile.Object.Pos.X
					}
				case data.Flying:
					if !ch.Flags.Flying {
						if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
					}
				case data.DoingAction:
					ch.Flags.Landing = false
					if ch.Flags.ItemAction == data.NoItemAction {
						ch.Object.Layer = ch.Layer
						if ch.Options.Flying || ch.Flags.Flying {
							ch.Flags.Flying = true
							ch.State = data.Flying
						} else if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
					}
				case data.Hiding:
					if tile.Block != data.BlockHideout || ch.Actions.Direction == data.Down {
						if ch.Options.Flying || ch.Flags.Flying {
							ch.Flags.Flying = true
							ch.State = data.Flying
						} else if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
					}
				case data.Hit:
					ch.Flags.Landing = false
					if ch.Flags.NextStep {
						if ch.Player > -1 {
							data.CurrLevelSess.PlayerStats[ch.Player].Deaths++
						} else if ch.Enemy > -1 && ch.Inventory != nil {
							if ch.Flags.Death == death.Drowned || ch.Flags.Death == death.Crushed ||
								ch.Flags.Death == death.Dying || SomethingOnTile(tile, ch.Object.ID) || !DropItem(ch) {
								myecs.Manager.DisposeEntity(ch.Inventory.Entity)
								ch.Inventory = nil
							}
						}
						ch.Flags.NextStep = false
						ch.State = data.Dead
						ch.Flags.Death = death.None
					}
				case data.Attack:
					if ch.Flags.NextStep {
						if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
					}
				case data.Tripping:
					if tile.Coords == ch.NextTile.Coords {
						DropItem(ch)
						if !tile.Flags.Collapse {
							ch.Object.Pos.X = tile.Object.Pos.X
							ch.Object.Pos.Y = tile.Object.Pos.Y
							ch.Flags.Death = death.Crushed
							ch.State = data.Hit
							if !ch.Flags.Ignore {
								sfx.SoundPlayer.PlaySound(constants.SFXCrush, 0.)
							}
						}
					}
					if ch.Flags.NextStep {
						ch.State = data.InHole
						ch.Object.Pos = tile.Object.Pos
						DropItem(ch)
					}
				case data.InHole:
					if ch.Flags.NextStep {
						ch.State = data.ClimbingOut
					} else if ch.ACounter > constants.DemonInHoleCounter {
						if above := data.CurrLevel.Get(x, y+1); !above.IsSolidPath() && !above.IsSolid() {
							if aboveL := data.CurrLevel.Get(x-1, y+1); !aboveL.IsSolidPath() && !aboveL.IsSolid() &&
								ch.Actions.Left() {
								ch.Flags.JumpL = true
							} else if aboveR := data.CurrLevel.Get(x+1, y+1); !aboveR.IsSolidPath() && !aboveR.IsSolid() &&
								ch.Actions.Right() {
								ch.Flags.JumpR = true
							}
						}
					}
				case data.ClimbingOut:
					if ch.Flags.NextStep {
						ch.Flags.NoCollision = false
						ch.NextTile = nil
						if ch.Object.Flip {
							ch.Object.Pos.X -= 1
						}
						if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
					}
				case data.AroundCorner:
					if ch.Flags.NextStep {
						ch.Flags.NoCollision = false
						ch.State = data.Grounded
						switch ch.Flags.Orientation {
						case data.Up:
							ch.Object.Rot = math.Pi
							ch.Object.Pos.Y = tile.Object.Pos.Y
						case data.Left:
							ch.Object.Rot = math.Pi * -0.5
							ch.Object.Pos.X = tile.Object.Pos.X
						case data.Right:
							ch.Object.Rot = math.Pi * 0.5
							ch.Object.Pos.X = tile.Object.Pos.X
						default:
							ch.Object.Rot = 0.
							ch.Object.Pos.Y = tile.Object.Pos.Y
						}
					}
				case data.Regen:
					if !ch.Flags.Regen {
						if ch.Options.Flying {
							ch.Flags.Flying = true
							ch.State = data.Flying
						} else if ch.Flags.WallClimb {
							ch.State = data.Grounded
						} else if ch.Flags.Floor || below.IsRunnable() {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else if tile != nil && tile.Block == data.BlockBar {
							ch.State = data.OnBar
						} else {
							ch.State = data.Falling
						}
					}
				case data.Dead:
					ch.Flags.Landing = false
					ch.Object.Layer = ch.Layer
					ch.Flags.Disguised = false
					ch.Flags.Flying = false
					ch.Flags.NextStep = false
					ch.Flags.NoCollision = false
					ch.Flags.JumpR = false
					ch.Flags.JumpL = false
					ch.Flags.Death = death.None
					ch.ACounter = 0
					ch.Control.ClearPrev()
					sfx.SoundPlayer.KillSound(ch.SFX)
					if ch.Enemy > -1 && ch.Options.Regen {
						var t *data.Tile
						if len(ch.Options.LinkedTiles) > 0 {
							t = GetRandomRegenTileFromList(ch.Options.LinkedTiles, ch.Flags.LastRegen)
						} else { // pick a random empty tile
							t = GetRandomRegenTile()
						}
						if t != nil {
							ch.Flags.LastRegen = &t.Coords
							tile = t
							ch.Object.SetPos(t.Object.Pos)
							if ch.Options.RegenFlip {
								ch.Object.Flip = random.Level.Intn(2) == 0
							}
							if ch.Options.RegenOrient {
								ch.Flags.Orientation = tile.Metadata.Orientation
								switch ch.Flags.Orientation {
								case data.Up:
									ch.Object.Rot = math.Pi
								case data.Left:
									ch.Object.Rot = math.Pi * -0.5
								case data.Right:
									ch.Object.Rot = math.Pi * 0.5
								default:
									ch.Object.Rot = 0.
								}
							}
							ch.State = data.Waiting
						}
					} else if ch.Player > -1 &&
						ch.Options.Regen &&
						len(ch.Options.LinkedTiles) > 0 {
						c := ch.Options.LinkedTiles[0]
						t := data.CurrLevel.Get(c.X, c.Y)
						if t != nil {
							tile = t
							ch.Object.SetPos(t.Object.Pos)
							ch.State = data.Waiting
						}
					}
				case data.Waiting:
					if ch.Player > -1 && ch.Actions.Any() {
						ch.State = data.Regen
						sfx.SoundPlayer.PlaySound(constants.SFXRegen, 0.)
						ch.Flags.Regen = true
						PlayerPortal(ch.Object.Layer+1, tile.Object.Pos)
					} else if ch.Enemy > -1 && !EnemyOnTile(tile) {
						ch.State = data.Regen
						ch.Flags.Regen = true
					}
				}
				UpdateInventory(ch)
				if ch.State != data.Dead &&
					ch.State != data.Hit &&
					ch.State != data.Attack &&
					ch.State != data.Tripping &&
					ch.State != data.DoingAction {
					if ch.Actions.PickUp {
						if ch.Player > -1 {
							PickUpOrDropItem(ch, ch.Player)
						} else if ch.Enemy > -1 {
							PickUpOrDropGem(ch, ch.Enemy)
						}
					}
					if ch.Actions.DigLeft && Dig(ch, true) {
					} else if ch.Actions.DigRight && Dig(ch, false) {
					} else if ch.Actions.DigLeft && Place(ch, true) {
					} else if ch.Actions.DigRight && Place(ch, false) {
					} else if ch.Actions.Action && DoAction(ch) {
					} else if ch.Actions.Bomb && PlaceSmallBomb(ch) {
					}
				} else if ch.State == data.Dead {
					ClearInv(ch)
				}
				if oldState != ch.State { // a state change happened
					ch.Flags.NextStep = false
					ch.ACounter = 0
					ch.Control.ClearPrev()
					ch.Actions.PrevDirection = data.NoDirection
					if tile != nil && oldState != data.AroundCorner {
						ch.Object.Pos.Y = tile.Object.Pos.Y
					}
					ch.Flags.Climbed = false
					ch.Flags.GoingUp = false
					ch.Object.PostPos = ch.Object.Pos
					if ch.State != data.Jumping {
						ch.Flags.HighJump = false
						ch.Flags.LongJump = false
					}
					// player sound effects
					if ch.Player > -1 {
						switch oldState {
						case data.Falling:
							sfx.SoundPlayer.KillSound(ch.SFX)
							if ch.State == data.Grounded {
								sfx.SoundPlayer.PlaySound(constants.SFXLand, 0.)
							}
						}
						switch ch.State {
						case data.Falling:
							ch.SFX = sfx.SoundPlayer.PlaySound(constants.SFXFall, 0.)
						}
					}
				}
			}
		}
	}
}

func OutsideMapSystem() {
	if reanimator.FrameSwitch && data.CurrLevel != nil {
		for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
			_, okO := result.Components[myecs.Object].(*object.Object)
			ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
			if okO && okC {
				if ch.State != data.Dead && ch.State != data.Waiting && ch.Flags.OutsideMap {
					ch.Flags.OutsideMap = false
					if ch.Player > -1 && ch.Player < constants.MaxPlayers {
						currPos := ch.Object.Pos.Add(ch.Object.Offset)
						x, y := world.WorldToMap(currPos.X, currPos.Y)
						lastTile := ch.LastTile
						if lastTile != nil {
							if x < lastTile.Coords.X {
								if t, ok := lastTile.Transitions[data.Left]; ok && (t.Open || data.CurrLevel.Continuity) {
									sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
									data.CurrLevel.Complete = true
									data.CurrLevel.ExitIndex = t.ExitIndex
									data.CurrLevel.StartCoords = &t.ExitTile
									return
								}
							} else if x > lastTile.Coords.X {
								if t, ok := lastTile.Transitions[data.Right]; ok && (t.Open || data.CurrLevel.Continuity) {
									sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
									data.CurrLevel.Complete = true
									data.CurrLevel.ExitIndex = t.ExitIndex
									data.CurrLevel.StartCoords = &t.ExitTile
									return
								}
							} else if y < lastTile.Coords.Y {
								if t, ok := lastTile.Transitions[data.Down]; ok && (t.Open || data.CurrLevel.Continuity) {
									sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
									data.CurrLevel.Complete = true
									data.CurrLevel.ExitIndex = t.ExitIndex
									data.CurrLevel.StartCoords = &t.ExitTile
									return
								}
							} else if y > lastTile.Coords.Y {
								if t, ok := lastTile.Transitions[data.Up]; ok && (t.Open || data.CurrLevel.Continuity) {
									sfx.SoundPlayer.PlaySound(constants.SFXExitLevel, 0.)
									data.CurrLevel.Complete = true
									data.CurrLevel.ExitIndex = t.ExitIndex
									data.CurrLevel.StartCoords = &t.ExitTile
									return
								}
							}
						}
					}
					// kill them
					ch.State = data.Dead
					if ch.Player > -1 {
						data.CurrLevelSess.PlayerStats[ch.Player].Deaths++
					}
				}
			}
		}
	}
}
