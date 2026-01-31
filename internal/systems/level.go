package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/controllers"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/internal/ui"
	"gemrunner/pkg/object"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gopxl/pixel"
)

func InitLevelTiles(level *data.Level) {
	var gemsCollected []world.Coords
	if completion, ok := data.CurrLevelSess.LevelMap[data.CurrLevelSess.PuzzleIndex]; ok {
		level.DoorsOpen = completion.Completed
		level.Continuity = completion.Continuity || data.CurrLevelSess.PuzzleSet.Metadata.Continuity == data.ContinuityAlwaysOn
		gemsCollected = completion.GemsCollected
	}

	for y, row := range level.Tiles.T {
		for x, tile := range row {
			tile.Coords = world.NewCoords(x, y)
			obj := object.New()
			obj.Pos = world.MapToWorld(tile.Coords)
			obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
			obj.Layer = 10
			tile.Object = obj
			tile.Live = true
			e := myecs.Manager.NewEntity().
				AddComponent(myecs.Object, obj).
				AddComponent(myecs.Tile, tile)
			tile.Entity = e
			// replace reanimator and items
			switch tile.Block {
			case data.BlockPhase:
				switch tile.Metadata.Phase {
				case 1, 2, 3, 4:
					tile.Flags.Collapse = true
					tile.Counter = 10
				case 5, 6, 7, 0:

				}
			case data.BlockBarrier:
				if tile.Metadata.Regenerate {
					if tile.Metadata.Toggle != level.DoorsOpen {
						tile.Flags.Collapse = true
						tile.Counter = 10
					}
				} else if tile.Metadata.Toggle {
					tile.Flags.Collapse = true
					tile.Counter = 10
				}
			case data.BlockClose:
				tile.Flags.Collapse = true
				tile.Counter = constants.CrackedCounter + 1
			case data.BlockPlayer1, data.BlockPlayer2, data.BlockPlayer3, data.BlockPlayer4:
				i := 0
				if tile.Block == data.BlockPlayer2 {
					i = 1
				} else if tile.Block == data.BlockPlayer3 {
					i = 2
				} else if tile.Block == data.BlockPlayer4 {
					i = 3
				}
				PlayerCharacter(obj.Pos, i, data.CurrReplay)
				tile.ToEmpty()
			case data.BlockDemon:
				DemonCharacter(obj.Pos, tile.Metadata)
				if tile.Metadata.Buried {
					tile.Block = data.BlockTurf
					tile.Metadata = data.DefaultMetadata()
				} else {
					tile.ToEmpty()
				}
			case data.BlockFly:
				FlyCharacter(obj.Pos, tile.Metadata)
				if tile.Metadata.Buried {
					tile.Block = data.BlockTurf
					tile.Metadata = data.DefaultMetadata()
				} else {
					tile.ToEmpty()
				}
			case data.BlockGem:
				if !world.CoordsIn(tile.Coords, gemsCollected) {
					CreateGem(obj.Pos, tile)
				}
				if tile.Metadata.Buried {
					tile.Block = data.BlockTurf
					tile.Metadata = data.DefaultMetadata()
				} else {
					tile.ToEmpty()
				}
			case data.BlockDoorHidden, data.BlockDoorVisible, data.BlockDoorLocked:
				CreateDoor(obj.Pos, tile)
				tile.ToEmpty()
			case data.BlockBigBombLit:
				key := tile.Block.SpriteString()
				CreateLitBomb(obj.Pos, key, "big", true, tile.Metadata.Regenerate, false, tile.Metadata.RegenDelay)
				if tile.Metadata.Buried {
					tile.Block = data.BlockTurf
					tile.Metadata = data.DefaultMetadata()
				} else {
					tile.ToEmpty()
				}
			case data.BlockSmallBombLit:
				key := tile.Block.SpriteString()
				CreateLitBomb(obj.Pos, key, "small", false, tile.Metadata.Regenerate, tile.Metadata.Buried, tile.Metadata.RegenDelay)
				if tile.Metadata.Buried {
					tile.Block = data.BlockTurf
					tile.Metadata = data.DefaultMetadata()
				} else {
					tile.ToEmpty()
				}
			case data.BlockJumpBoots, data.BlockBox, data.BlockKey,
				data.BlockBigBomb, data.BlockSmallBomb, data.BlockJetpack,
				data.BlockDisguise, data.BlockDrill, data.BlockFlamethrower,
				data.BlockGoopBucket, data.BlockTransporter:
				CreateItem(tile.Block, obj.Pos, tile.SpriteString(), tile.Metadata, tile.Coords)
				if tile.Metadata.Buried {
					tile.Block = data.BlockTurf
					tile.Metadata = data.DefaultMetadata()
				} else {
					tile.ToEmpty()
				}
			case data.BlockTransporterExit:
				CreateTransporterExit(obj.Pos, tile)
				//case data.BlockGear:
				//	var a *reanimator.Anim
				//	if (tile.Coords.X+tile.Coords.Y)%2 == 0 {
				//		a = reanimator.NewBatchAnimationCustom("gear", img.Batchers[constants.TileBatch], "gear", []int{3, 0, 1, 2}, reanimator.Loop).Reverse()
				//	} else {
				//		a = reanimator.NewBatchAnimation("gear", img.Batchers[constants.TileBatch], "gear", reanimator.Loop)
				//	}
				//	anim := reanimator.NewSimple(a)
				//	tile.Entity.AddComponent(myecs.Drawable, anim)
				//	tile.Entity.AddComponent(myecs.Animated, anim)
			case data.BlockLiquid:
				tile.Object.Layer = 30
			}
			AddLevelTransition(tile, x, y)
		}
	}
}

func InitContinuity(level *data.Level) {
	if data.CurrLevelSess.PuzzleSet.Metadata.Continuity != data.NoContinuity {
		var t *data.Tile
		if data.CurrLevelSess.StartCoords != nil {
			t = level.Get(data.CurrLevelSess.StartCoords.X, data.CurrLevelSess.StartCoords.Y)
		}
		for i, p := range level.Players {
			if p != nil {
				if t != nil && !t.IsSolid() {
					p.Object.SetPos(t.Object.Pos)
				}
				p.SmallBombs = data.CurrLevelSess.PlayerStats[i].CurrBombs
				inv := data.CurrLevelSess.PlayerStats[i].Inventory
				if inv != nil {
					if inv.Entity == nil || inv.Object == nil { // create item
						item := CreateItem(inv.Block, pixel.ZV, inv.Key, inv.Metadata, inv.Origin)
						p.Inventory = item
						p.Inventory.PickUp.Inventory = i
					} else {
						p.Inventory = inv
						p.Inventory.Entity.AddComponent(myecs.Temp, myecs.ClearFlag(false))
					}
				}
			}
		}
	}
}

func CreateFakePlayer(level *data.Level) {
	if level == nil {
		return
	}
	tile := GetRandomRegenTile()
	if tile == nil {
		x := random.Level.Intn(level.Metadata.Width)
		y := random.Level.Intn(level.Metadata.Height)
		tile = level.Get(x, y)
	}
	ch := data.NewDynamic()
	ch.LastTile = tile
	ch.Layer = 0
	e := myecs.Manager.NewEntity()
	obj := object.New().WithFixedID("fake_player").SetPos(tile.Object.Pos)
	obj.Pos = world.MapToWorld(tile.Coords)
	obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
	obj.Layer = 0
	obj.SetRect(pixel.R(0., 0., 16., 16.))
	ch.Object = obj
	ch.State = data.Grounded
	ch.Vars = data.DemonVars()
	ch.Control = controllers.NewRandomWalk(ch, e)
	ch.Flags.Ignore = true
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Dynamic, ch).
		AddComponent(myecs.Controller, ch.Control)
	ch.Entity = e
	level.FakePlayer = ch
}

func InitLevelDialogs(level *data.Level) {
	for p := 0; p < constants.MaxPlayers; p++ {
		if p < data.CurrPuzzleSet.Metadata.NumPlayers {
			var dlgKey string
			switch p {
			case 0:
				dlgKey = constants.DialogPlayer1Inv
			case 1:
				dlgKey = constants.DialogPlayer2Inv
			case 2:
				dlgKey = constants.DialogPlayer3Inv
			case 3:
				dlgKey = constants.DialogPlayer4Inv
			}
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedPrimary", float32(level.Puzzle.Metadata.PrimaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenPrimary", float32(level.Puzzle.Metadata.PrimaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBluePrimary", float32(level.Puzzle.Metadata.PrimaryColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedSecondary", float32(level.Puzzle.Metadata.SecondaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenSecondary", float32(level.Puzzle.Metadata.SecondaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueSecondary", float32(level.Puzzle.Metadata.SecondaryColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedDoodad", float32(level.Puzzle.Metadata.DoodadColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenDoodad", float32(level.Puzzle.Metadata.DoodadColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueDoodad", float32(level.Puzzle.Metadata.DoodadColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedGoop", float32(level.Puzzle.Metadata.GoopColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenGoop", float32(level.Puzzle.Metadata.GoopColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueGoop", float32(level.Puzzle.Metadata.GoopColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedLiquidPrimary", float32(level.Puzzle.Metadata.LiquidPrimaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenLiquidPrimary", float32(level.Puzzle.Metadata.LiquidPrimaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueLiquidPrimary", float32(level.Puzzle.Metadata.LiquidPrimaryColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedLiquidSecondary", float32(level.Puzzle.Metadata.LiquidSecondaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenLiquidSecondary", float32(level.Puzzle.Metadata.LiquidSecondaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueLiquidSecondary", float32(level.Puzzle.Metadata.LiquidSecondaryColor.B))
			level.PLoc[p] = &mgl32.Vec2{}
			level.PLoc[p][0] = float32(level.Players[p].Object.Pos.X)
			level.PLoc[p][1] = float32(level.Players[p].Object.Pos.Y)
		} else {
			level.PLoc[p] = &mgl32.Vec2{}
			level.PLoc[p][0] = -1
			level.PLoc[p][1] = -1
		}
	}
}

func DisposeCurrLevel() {
	DisposeLevel(data.CurrLevel)
	data.CurrLevel = nil
	if data.CurrentPlayArea != nil {
		DisposeLevel(data.CurrentPlayArea.Level)
		data.CurrentPlayArea.Level = nil
	}
}

func DisposeLevel(level *data.Level) {
	if level == nil {
		return
	}
	for _, row := range level.Tiles.T {
		for _, tile := range row {
			if tile != nil {
				if tile.FloatingText != nil {
					myecs.Manager.DisposeEntity(tile.FloatingText.Entity)
					myecs.Manager.DisposeEntity(tile.FloatingText.ShEntity)
					data.RemoveFloatingText(tile)
				}
				if tile.Entity != nil {
					myecs.Manager.DisposeEntity(tile.Entity)
					tile.Entity = nil
				}
			}
		}
	}
	for i, player := range level.Players {
		if player != nil {
			sfx.SoundPlayer.KillSound(player.SFX)
			myecs.Manager.DisposeEntity(player.Entity)
		}
		level.Players[i] = nil
	}
	for _, enemy := range level.Enemies {
		myecs.Manager.DisposeEntity(enemy.Entity)
	}
	level.Enemies = []*data.Dynamic{}
	if level.FakePlayer != nil {
		myecs.Manager.DisposeEntity(level.FakePlayer.Entity)
		level.FakePlayer = nil
	}
	for _, entity := range level.AllEntities {
		myecs.Manager.DisposeEntity(entity)
	}
	level.AllEntities = []*ecs.Entity{}
	level = nil
}

func AddLevelTransition(tile *data.Tile, x, y int) {
	// level transitions
	tile.Transitions = make(map[data.Direction]*data.LevelTransition)
	if data.CurrPuzzleSet.Metadata.Adventure && data.CurrPuzzleSet.Metadata.Continuity != data.NoContinuity {
		var ht *data.Tile
		hi := -1
		var hx, hy int
		var ho, hc bool
		var hd data.Direction
		if x == 0 {
			hd = data.Left
			l := data.CurrLevel.Puzzle.Grid
			l.X--
			if left := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(l); left != nil {
				hi = data.CurrLevelSess.PuzzleSet.GetGrid(l)
				m := left.Metadata.Height
				n := data.CurrLevel.Puzzle.Metadata.Height
				hx = left.Metadata.Width - 1
				if m == n { // levels are the same height
					hy = y
				} else if m < n { // next level is shorter than this one
					hy = y - (n-m)/2
				} else { // next level is taller
					hy = y + (n-m)/2
				}
				ht = left.Get(hx, hy)
			}
		} else if x == data.CurrLevel.Metadata.Width-1 {
			hd = data.Right
			r := data.CurrLevel.Puzzle.Grid
			r.X++
			if right := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(r); right != nil {
				hi = data.CurrLevelSess.PuzzleSet.GetGrid(r)
				m := right.Metadata.Height
				n := data.CurrLevel.Puzzle.Metadata.Height
				hx = 0
				if m == n { // levels are the same width
					hy = y
				} else if m < n { // next level is shorter than this one
					hy = y - (n-m)/2
				} else { // next level is taller
					hy = y + (n-m)/2
				}
				ht = right.Get(hx, hy)
			}
		}
		if lc, ok := data.CurrLevelSess.LevelMap[hi]; ok {
			ho = data.CurrLevelSess.PuzzleSet.Metadata.Continuity == data.ContinuityAlwaysOn ||
				hi == data.CurrLevelSess.LastPuzzle ||
				(lc.Completed && data.CurrLevel.DoorsOpen)
			hc = lc.Completed
		}
		if ht != nil && !ht.IsSolidLevelTrans(hc) {
			tile.Transitions[hd] = &data.LevelTransition{
				Open:      ho,
				ExitIndex: hi,
				ExitTile:  world.Coords{X: hx, Y: hy},
			}
		}
		var vt *data.Tile
		vi := -1
		var vx, vy int
		var vo, vc bool
		var vd data.Direction
		if y == 0 {
			vd = data.Down
			b := data.CurrLevel.Puzzle.Grid
			b.Y--
			if below := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(b); below != nil {
				vi = data.CurrLevelSess.PuzzleSet.GetGrid(b)
				m := below.Metadata.Width
				n := data.CurrLevel.Puzzle.Metadata.Width
				vy = below.Metadata.Height - 1
				if m == n { // levels are the same width
					vx = x
				} else if m < n { // next level is less wide than this one
					vx = x - (n-m)/2
				} else { // next level is wider
					vx = x + (n-m)/2
				}
				vt = below.Get(vx, vy)
			}
		} else if y == data.CurrLevel.Metadata.Height-1 {
			vd = data.Up
			a := data.CurrLevel.Puzzle.Grid
			a.Y++
			if above := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(a); above != nil {
				vi = data.CurrLevelSess.PuzzleSet.GetGrid(a)
				m := above.Metadata.Width
				n := data.CurrLevel.Puzzle.Metadata.Width
				vy = 0
				if m == n { // levels are the same width
					vx = x
				} else if m < n { // next level is less wide than this one
					vx = x - (n-m)/2
				} else { // next level is wider
					vx = x + (n-m)/2
				}
				vt = above.Get(vx, vy)
			}
		}
		if lc, ok := data.CurrLevelSess.LevelMap[vi]; ok {
			vo = data.CurrLevelSess.PuzzleSet.Metadata.Continuity == data.ContinuityAlwaysOn ||
				vi == data.CurrLevelSess.LastPuzzle ||
				(lc.Completed && data.CurrLevel.DoorsOpen)
			vc = lc.Completed
		}
		if vt != nil && !vt.IsSolidLevelTrans(vc) {
			tile.Transitions[vd] = &data.LevelTransition{
				Open:      vo,
				ExitIndex: vi,
				ExitTile:  world.Coords{X: vx, Y: vy},
			}
		}
	}
}

func GetBestRegenTile(tiles []*data.Tile) *data.Tile {
	if len(tiles) < 1 {
		return nil
	}
	var bestTiles []*data.Tile
	for _, t := range tiles {
		dist := DistanceToClosestPlayer(t)
		if dist == -1 || dist > 2 {
			bestTiles = append(bestTiles, t)
		}
	}
	if len(bestTiles) > 0 {
		return bestTiles[random.Level.Intn(len(bestTiles))]
	}
	return tiles[random.Level.Intn(len(tiles))]
}

func GetRandomRegenTileFromList(coords []world.Coords, exclude *world.Coords) *data.Tile {
	var tiles []*data.Tile
	in := false
	for _, c := range coords {
		if exclude == nil || c != *exclude {
			tile := data.CurrLevel.Get(c.X, c.Y)
			tiles = append(tiles, tile)
		} else if exclude != nil && c == *exclude {
			in = true
		}
	}
	if in && len(tiles) == 0 {
		tile := data.CurrLevel.Get(exclude.X, exclude.Y)
		tiles = append(tiles, tile)
	}
	return GetBestRegenTile(tiles)
}

func GetRandomRegenTile() *data.Tile {
	var tiles []*data.Tile
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok && tile.Live {
			if tile.IsEmpty() && tile.Block != data.BlockDemonRegen {
				tiles = append(tiles, tile)
			}
		}
	}
	return GetBestRegenTile(tiles)
}

func DistanceToClosestPlayer(tile *data.Tile) int {
	dist := -1
	for _, p := range data.CurrLevel.Players {
		if p == nil {
			continue
		}
		px, py := world.WorldToMap(p.Object.Pos.X, p.Object.Pos.Y)
		d := world.DistanceOrthogonal(world.Coords{X: px, Y: py}, tile.Coords)
		if dist == -1 || d < dist {
			dist = d
		}
	}
	return dist
}

func EnemyOnTile(tile *data.Tile) bool {
	for _, result := range myecs.Manager.Query(myecs.IsEnemy) {
		obj, ok := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if ok && okC {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			if ch.State != data.Dead && ch.State != data.Waiting &&
				x == tile.Coords.X && y == tile.Coords.Y {
				return true
			}
		}
	}
	return false
}

func SomethingOnTile(tile *data.Tile) bool {
	for _, result := range myecs.Manager.Query(myecs.IsLvlElement) {
		obj, ok := result.Components[myecs.Object].(*object.Object)
		if ok {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			if x == tile.Coords.X && y == tile.Coords.Y {
				return true
			}
		}
	}
	return false
}

func ThingsOnTile(tile *data.Tile) []*ecs.Entity {
	var things []*ecs.Entity
	for _, result := range myecs.Manager.Query(myecs.IsLvlElement) {
		obj, ok := result.Components[myecs.Object].(*object.Object)
		if ok {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			if x == tile.Coords.X && y == tile.Coords.Y {
				things = append(things, result.Entity)
			}
		}
	}
	return things
}
