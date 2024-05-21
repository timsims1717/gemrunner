package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/world"
)

func CharacterStateSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC {
			if reanimator.FrameSwitch {
				currPos := ch.Object.Pos.Add(ch.Object.Offset)
				x, y := world.WorldToMap(currPos.X, currPos.Y)
				tile := data.CurrLevel.Tiles.Get(x, y)
				below := data.CurrLevel.Tiles.Get(x, y-1)
				oldState := ch.State
				switch ch.State {
				case data.Grounded:
					if ch.Flags.HighJump || ch.Flags.LongJump {
						ch.State = data.Jumping
					} else if tile != nil && tile.Block == data.BlockBar {
						ch.State = data.OnBar
					} else if tile != nil && tile.IsLadder() { // a ladder is here
						if ch.Actions.Direction == data.Up { // climbed the ladder
							ch.State = data.OnLadder
						} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
							ch.State = data.OnLadder
						} else if below != nil && !below.IsRunnable() && (ch.Actions.Left() || ch.Actions.Right()) { // leaping onto ladder
							ch.State = data.Leaping
							ch.Flags.LeapOn = true
							ch.NextTile = tile
						} else if !ch.Flags.Floor {
							ch.State = data.OnLadder
						}
					} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
						ch.State = data.OnLadder
					} else if !ch.Flags.Floor && below != nil && !below.IsLadder() {
						ch.State = data.Falling
					}
				case data.OnLadder:
					if ch.Flags.Floor && ch.Actions.Up() { // just got to the top
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
							left := data.CurrLevel.Tiles.Get(x-1, y)
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
							right := data.CurrLevel.Tiles.Get(x+1, y)
							ch.NextTile = right
							if right != nil && right.IsLadder() { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						}
					} else if tile.Block == data.BlockBar && !ch.Actions.Down() {
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
						if ch.Actions.Direction == data.Left { // run to the left
							ch.State = data.Grounded
							ch.Object.Flip = true
						} else if ch.Actions.Direction == data.Right { // run to the right
							ch.State = data.Grounded
							ch.Object.Flip = false
						} else if ch.Flags.Floor && ch.Actions.Direction == data.Down {
							ch.State = data.Grounded
						}
					} else if !tile.IsLadder() {
						ch.State = data.Falling
					}
				case data.OnBar:
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
						ch.State = data.Grounded
					} else if tile != nil && tile.IsLadder() {
						ch.State = data.OnLadder
						ch.Object.Pos.X = tile.Object.Pos.X
					} else if tile != nil &&
						tile.Block == data.BlockBar &&
						tile.Coords != ch.LastTile.Coords {
						ch.State = data.OnBar
					} else {
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
					if ch.Flags.ItemAction == data.NoItemAction {
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
					if !ch.Flags.Hit {
						ch.State = data.Dead
						if ch.Player > -1 {
							data.CurrLevelSess.PlayerStats[ch.Player].Deaths++
						}
					}
				case data.Attack:
					if !ch.Flags.Attack {
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
				case data.Regen:
					if !ch.Flags.Regen {
						if ch.Options.Flying {
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
				case data.Dead:
					ch.Object.Layer = ch.Layer
					ch.Flags.Flying = false
					ch.Flags.Attack = false
					ch.Flags.Hit = false
					if ch.Enemy > -1 &&
						ch.Options.Regen {
						var t *data.Tile
						if len(ch.Options.LinkedTiles) > 0 {
							t = GetRandomRegenTileFromList(ch.Options.LinkedTiles)
						} else { // pick a random empty tile
							t = GetRandomRegenTile()
						}
						if t != nil {
							tile = t
							ch.Object.SetPos(t.Object.Pos)
							if ch.Options.RegenFlip {
								ch.Object.Flip = random.Level.Intn(2) == 0
							}
							ch.State = data.Regen
							ch.Flags.Regen = true
						}
					} else if ch.Player > -1 &&
						ch.Options.Regen &&
						len(ch.Options.LinkedTiles) > 0 {
						c := ch.Options.LinkedTiles[0]
						t := data.CurrLevel.Tiles.Get(c.X, c.Y)
						if t != nil {
							tile = t
							ch.Object.SetPos(t.Object.Pos)
							ch.State = data.Waiting
						}
					}
				case data.Waiting:
					if ch.Actions.Any() {
						ch.State = data.Regen
						sfx.SoundPlayer.PlaySound(constants.SFXRegen, 0.)
						ch.Flags.Regen = true
						PlayerPortal(ch.Object.Layer+1, tile.Object.Pos)
					}
				}
				UpdateInventory(ch)
				if ch.State != data.Dead &&
					ch.State != data.Hit &&
					ch.State != data.Attack &&
					ch.State != data.DoingAction {
					if ch.Actions.PickUp {
						PickUpOrDropItem(ch, ch.Player)
					}
					if ch.State != data.Leaping {
						if ch.Actions.DigLeft && Dig(ch, true) {
						} else if ch.Actions.DigRight && Dig(ch, false) {
						} else if ch.Actions.DigLeft && Place(ch, true) {
						} else if ch.Actions.DigRight && Place(ch, false) {
						} else if ch.Actions.Action && DoAction(ch) {
						}
					}
				} else if ch.State == data.Dead {
					DropItem(ch)
				}
				if oldState != ch.State { // a state change happened
					ch.ACounter = 0
					ch.Control.ClearPrev()
					ch.Actions.PrevDirection = data.None
					ch.Object.Pos.Y = tile.Object.Pos.Y
					if (oldState == data.Grounded &&
						(ch.State == data.OnLadder ||
							ch.State == data.Falling ||
							ch.State == data.Jumping)) ||
						(oldState == data.Falling &&
							(ch.State == data.Grounded ||
								ch.State == data.OnLadder ||
								ch.State == data.OnBar)) ||
						(oldState == data.Leaping) ||
						(oldState == data.Jumping) {
						ch.Object.Pos.X = tile.Object.Pos.X
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
