package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
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

func LevelInit(record bool) {
	if data.CurrPuzzleSet.CurrPuzzle == nil {
		panic("no puzzle loaded to start level")
	}
	LevelDispose()
	PuzzleViewInit()
	data.CurrLevel = &data.Level{}
	data.CurrLevel.Tiles = data.CurrPuzzleSet.CurrPuzzle.CopyTiles()
	data.CurrLevel.Metadata = data.CurrPuzzleSet.CurrPuzzle.Metadata
	data.CurrLevel.Puzzle = data.CurrPuzzleSet.CurrPuzzle
	FloatingTextStartLevel()
	SetPuzzleTitle()
	UpdatePuzzleTimer()
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	levelSeed := random.RandomSeed()
	random.SetLevelSeed(levelSeed)
	data.CurrLevel.Recording = data.CurrReplay == nil && record
	data.CurrLevel.SaveRecord = constants.Configuration.Gameplay.AlwaysRecord
	if data.CurrLevel.Recording {
		data.CurrLevel.LevelReplay = &data.LevelReplay{
			PuzzleSet:   data.CurrPuzzleSet.Metadata.Name,
			Filename:    data.CurrPuzzleSet.Metadata.Filename,
			ReplayFile:  content.ReplayFile(data.CurrPuzzleSet.Metadata.Name, data.CurrPuzzleSet.PuzzleIndex),
			PuzzleNum:   data.CurrPuzzleSet.PuzzleIndex,
			StartCoords: data.CurrLevelSess.StartCoords,
			Seed:        levelSeed,
		}
		data.CurrLevel.ReplayFrame = data.ReplayFrame{}
	} else if data.CurrReplay != nil {
		data.CurrReplay.FrameIndex = 0
		data.CurrLevelSess.StartCoords = data.CurrReplay.StartCoords
	}

	var gemsCollected []world.Coords
	if completion, ok := data.CurrLevelSess.LevelMap[data.CurrLevelSess.PuzzleIndex]; ok {
		data.CurrLevel.DoorsOpen = completion.Completed
		data.CurrLevel.Continuity = completion.Continuity || data.CurrLevelSess.PuzzleSet.Metadata.Continuity == data.ContinuityAlwaysOn
		gemsCollected = completion.GemsCollected
	}

	for y, row := range data.CurrLevel.Tiles.T {
		for x, tile := range row {
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
				if tile.Metadata.Toggle {
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
				PlayerCharacter(obj.Pos, i, tile, data.CurrReplay)
				tile.Block = data.BlockEmpty
			case data.BlockDemon:
				DemonCharacter(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockFly:
				FlyCharacter(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockGem:
				if !world.CoordsIn(tile.Coords, gemsCollected) {
					CreateGem(obj.Pos, tile)
				}
				tile.Block = data.BlockEmpty
			case data.BlockDoorHidden, data.BlockDoorVisible, data.BlockDoorLocked:
				CreateDoor(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockJumpBoots:
				CreateJumpBoots(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockBox:
				CreateBox(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockKey:
				CreateKey(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockBigBomb:
				CreateBomb(obj.Pos, tile, "big", true)
				tile.Block = data.BlockEmpty
			case data.BlockBigBombLit:
				key := tile.Block.SpriteString()
				CreateLitBomb(obj.Pos, key, "big", true, tile.Metadata.Regenerate, tile.Metadata.RegenDelay)
				tile.Block = data.BlockEmpty
			case data.BlockSmallBomb:
				CreateBomb(obj.Pos, tile, "small", false)
				tile.Block = data.BlockEmpty
			case data.BlockSmallBombLit:
				key := tile.Block.SpriteString()
				CreateLitBomb(obj.Pos, key, "small", false, tile.Metadata.Regenerate, tile.Metadata.RegenDelay)
				tile.Block = data.BlockEmpty
			case data.BlockJetpack:
				CreateJetpack(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockDisguise:
				CreateDisguise(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockDrill:
				CreateDrill(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockFlamethrower:
				CreateFlamethrower(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockTransporter:
				CreateTransporter(obj.Pos, tile)
				tile.Block = data.BlockEmpty
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
	if data.CurrLevelSess.StartCoords != nil && (data.CurrLevel.DoorsOpen || data.CurrLevelSess.PuzzleSet.Metadata.Continuity == data.ContinuityAlwaysOn) {
		if t := data.CurrLevel.Get(data.CurrLevelSess.StartCoords.X, data.CurrLevelSess.StartCoords.Y); t != nil && !t.IsSolid() {
			for _, p := range data.CurrLevel.Players {
				p.Object.SetPos(t.Object.Pos)
			}
		}
	}
	CreateFakePlayer()
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
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.CurrLevel.Puzzle.Metadata.PrimaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.CurrLevel.Puzzle.Metadata.PrimaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.CurrLevel.Puzzle.Metadata.PrimaryColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.CurrLevel.Puzzle.Metadata.SecondaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.CurrLevel.Puzzle.Metadata.SecondaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.CurrLevel.Puzzle.Metadata.SecondaryColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.CurrLevel.Puzzle.Metadata.DoodadColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.CurrLevel.Puzzle.Metadata.DoodadColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.CurrLevel.Puzzle.Metadata.DoodadColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedGoop", float32(data.CurrLevel.Puzzle.Metadata.GoopColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenGoop", float32(data.CurrLevel.Puzzle.Metadata.GoopColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueGoop", float32(data.CurrLevel.Puzzle.Metadata.GoopColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedLiquidPrimary", float32(data.CurrLevel.Puzzle.Metadata.LiquidPrimaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenLiquidPrimary", float32(data.CurrLevel.Puzzle.Metadata.LiquidPrimaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueLiquidPrimary", float32(data.CurrLevel.Puzzle.Metadata.LiquidPrimaryColor.B))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uRedLiquidSecondary", float32(data.CurrLevel.Puzzle.Metadata.LiquidSecondaryColor.R))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uGreenLiquidSecondary", float32(data.CurrLevel.Puzzle.Metadata.LiquidSecondaryColor.G))
			ui.Dialogs[dlgKey].ViewPort.Canvas.SetUniform("uBlueLiquidSecondary", float32(data.CurrLevel.Puzzle.Metadata.LiquidSecondaryColor.B))
			data.CurrLevel.PLoc[p] = &mgl32.Vec2{}
			data.CurrLevel.PLoc[p][0] = float32(data.CurrLevel.Players[p].Object.Pos.X)
			data.CurrLevel.PLoc[p][1] = float32(data.CurrLevel.Players[p].Object.Pos.Y)
		} else {
			data.CurrLevel.PLoc[p] = &mgl32.Vec2{}
			data.CurrLevel.PLoc[p][0] = -1
			data.CurrLevel.PLoc[p][1] = -1
		}
	}
	UpdatePuzzleShaders()
	ChangeWorldShader(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode)
}

func LevelDispose() {
	if data.CurrLevel != nil {
		for _, row := range data.CurrLevel.Tiles.T {
			for _, tile := range row {
				if tile.FloatingText != nil {
					myecs.Manager.DisposeEntity(tile.FloatingText.Entity)
					myecs.Manager.DisposeEntity(tile.FloatingText.ShEntity)
				}
				myecs.Manager.DisposeEntity(tile.Entity)
			}
		}
		for _, player := range data.CurrLevel.Players {
			if player != nil {
				sfx.SoundPlayer.KillSound(player.SFX)
				myecs.Manager.DisposeEntity(player.Entity)
			}
		}
		for _, enemy := range data.CurrLevel.Enemies {
			myecs.Manager.DisposeEntity(enemy.Entity)
		}
		if data.CurrLevel.FakePlayer != nil {
			myecs.Manager.DisposeEntity(data.CurrLevel.FakePlayer.Entity)
		}
		data.CurrLevel = nil
	}
}

func AddLevelTransition(tile *data.Tile, x, y int) {
	// level transitions
	tile.Transitions = make(map[data.Direction]*data.LevelTransition)
	if data.CurrPuzzleSet.Metadata.Adventure && data.CurrPuzzleSet.Metadata.Continuity != data.NoContinuity {
		var ht *data.Tile
		var hi int
		var hc bool
		var hd data.Direction
		if x == 0 {
			hd = data.Left
			l := data.CurrLevel.Puzzle.Grid
			l.X--
			if left := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(l); left != nil {
				hi = data.CurrLevelSess.PuzzleSet.GetGrid(l)
				m := left.Metadata.Height
				n := data.CurrLevel.Puzzle.Metadata.Height
				if m == n { // levels are the same height
					ht = left.Get(left.Metadata.Width-1, y)
				} else if m < n { // next level is shorter than this one
					ht = left.Get(left.Metadata.Width-1, y-(n-m)/2)
				} else { // next level is taller
					ht = left.Get(left.Metadata.Width-1, y+(n-m)/2)
				}
			}
		} else if x == data.CurrLevel.Metadata.Width-1 {
			hd = data.Right
			r := data.CurrLevel.Puzzle.Grid
			r.X++
			if right := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(r); right != nil {
				hi = data.CurrLevelSess.PuzzleSet.GetGrid(r)
				m := right.Metadata.Height
				n := data.CurrLevel.Puzzle.Metadata.Height
				if m == n { // levels are the same width
					ht = right.Get(0, y)
				} else if m < n { // next level is shorter than this one
					ht = right.Get(0, y-(n-m)/2)
				} else { // next level is taller
					ht = right.Get(0, y+(n-m)/2)
				}
			}
		}
		if lc, ok := data.CurrLevelSess.LevelMap[hi]; ok {
			hc = lc.Completed
		}
		if ht != nil && (!ht.IsSolid() || (ht.Block == data.BlockBarrier && ht.Metadata.Toggle != hc)) {
			tile.Transitions[hd] = &data.LevelTransition{
				ExitIndex: hi,
				ExitTile:  ht.Coords,
			}
		}
		var vt *data.Tile
		var vi int
		var vc bool
		var vd data.Direction
		if y == 0 {
			vd = data.Down
			b := data.CurrLevel.Puzzle.Grid
			b.Y--
			if below := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(b); below != nil {
				vi = data.CurrLevelSess.PuzzleSet.GetGrid(b)
				m := below.Metadata.Width
				n := data.CurrLevel.Puzzle.Metadata.Width
				if m == n { // levels are the same width
					vt = below.Get(x, below.Metadata.Height-1)
				} else if m < n { // next level is less wide than this one
					vt = below.Get(x-(n-m)/2, below.Metadata.Height-1)
				} else { // next level is wider
					vt = below.Get(x+(n-m)/2, below.Metadata.Height-1)
				}
			}
		} else if y == data.CurrLevel.Metadata.Height-1 {
			vd = data.Up
			a := data.CurrLevel.Puzzle.Grid
			a.Y++
			if above := data.CurrLevelSess.PuzzleSet.GetGridPuzzle(a); above != nil {
				vi = data.CurrLevelSess.PuzzleSet.GetGrid(a)
				m := above.Metadata.Width
				n := data.CurrLevel.Puzzle.Metadata.Width
				if m == n { // levels are the same width
					vt = above.Get(x, 0)
				} else if m < n { // next level is less wide than this one
					vt = above.Get(x-(n-m)/2, 0)
				} else { // next level is wider
					vt = above.Get(x+(n-m)/2, 0)
				}
			}
		}
		if lc, ok := data.CurrLevelSess.LevelMap[vi]; ok {
			vc = lc.Completed
		}
		if vt != nil && (!vt.IsSolid() || (vt.Block == data.BlockBarrier && vt.Metadata.Toggle != vc)) {
			tile.Transitions[vd] = &data.LevelTransition{
				ExitIndex: vi,
				ExitTile:  vt.Coords,
			}
		}
	}
}

func CreateFakePlayer() {
	if data.CurrLevel == nil {
		return
	}
	tile := GetRandomRegenTile()
	if tile == nil {
		x := random.Level.Intn(data.CurrLevel.Metadata.Width)
		y := random.Level.Intn(data.CurrLevel.Metadata.Height)
		tile = data.CurrLevel.Get(x, y)
	}
	ch := data.NewDynamic(tile)
	ch.Layer = 0
	e := myecs.Manager.NewEntity()
	obj := object.New().WithID("fake_player").SetPos(tile.Object.Pos)
	obj.Pos = world.MapToWorld(tile.Coords)
	obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
	obj.Layer = 0
	obj.SetRect(pixel.R(0., 0., 16., 16.))
	ch.Object = obj
	ch.State = data.Grounded
	ch.Vars = data.DemonVars()
	ch.Control = controllers.NewRandomWalk(ch, e)
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Temp, myecs.ClearFlag(false)).
		AddComponent(myecs.Dynamic, ch).
		AddComponent(myecs.Controller, ch.Control)
	ch.Entity = e
	data.CurrLevel.FakePlayer = ch
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

func GetRandomRegenTileFromList(coords []world.Coords) *data.Tile {
	var tiles []*data.Tile
	for _, c := range coords {
		tile := data.CurrLevel.Get(c.X, c.Y)
		if !SomethingOnTile(tile) {
			tiles = append(tiles, tile)
		}
	}
	return GetBestRegenTile(tiles)
}

func GetRandomRegenTile() *data.Tile {
	var tiles []*data.Tile
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tile, ok := result.Components[myecs.Tile].(*data.Tile)
		if okO && ok && tile.Live {
			if tile.IsEmpty() && !SomethingOnTile(tile) && tile.Block != data.BlockDemonRegen {
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
