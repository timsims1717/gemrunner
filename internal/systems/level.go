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
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	levelSeed := random.RandomSeed()
	random.SetLevelSeed(levelSeed)
	data.CurrLevel.Recording = data.CurrReplay == nil && record
	data.CurrLevel.SaveRecord = false
	if data.CurrLevel.Recording {
		data.CurrLevel.LevelReplay = &data.LevelReplay{
			PuzzleSet:  data.CurrPuzzleSet.Metadata.Name,
			Filename:   data.CurrPuzzleSet.Metadata.Filename,
			ReplayFile: content.ReplayFile(data.CurrPuzzleSet.Metadata.Name, data.CurrPuzzleSet.PuzzleIndex),
			PuzzleNum:  data.CurrPuzzleSet.PuzzleIndex,
			Seed:       levelSeed,
		}
		data.CurrLevel.ReplayFrame = data.ReplayFrame{}
	} else if data.CurrReplay != nil {
		data.CurrReplay.FrameIndex = 0
	}

	for _, row := range data.CurrLevel.Tiles.T {
		for _, tile := range row {
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
				CreateGem(obj.Pos, tile)
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
			case data.BlockBomb:
				CreateBomb(obj.Pos, tile)
				tile.Block = data.BlockEmpty
			case data.BlockBombLit:
				key := tile.Block.SpriteString()
				CreateLitBomb(obj.Pos, key, tile.Metadata)
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
