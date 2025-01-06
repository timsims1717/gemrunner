package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/controllers"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func LevelInit() {
	if data.CurrPuzzleSet.CurrPuzzle == nil {
		panic("no puzzle loaded to start level")
	}
	LevelDispose()
	data.CurrLevel = &data.Level{}
	data.CurrLevel.Tiles = data.CurrPuzzleSet.CurrPuzzle.CopyTiles()
	data.CurrLevel.Metadata = data.CurrPuzzleSet.CurrPuzzle.Metadata
	data.CurrLevel.Puzzle = data.CurrPuzzleSet.CurrPuzzle
	FloatingTextStartLevel()
	SetPuzzleTitle()
	data.CurrLevelSess.PuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
	random.SetLevelSeed(random.RandomSeed())

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
			case data.BlockPlayer1, data.BlockPlayer2, data.BlockPlayer3, data.BlockPlayer4:
				i := 0
				if tile.Block == data.BlockPlayer2 {
					i = 1
				} else if tile.Block == data.BlockPlayer3 {
					i = 2
				} else if tile.Block == data.BlockPlayer4 {
					i = 3
				}
				tile.Block = data.BlockEmpty
				PlayerCharacter(obj.Pos, i, tile)
			case data.BlockDemon:
				tile.Block = data.BlockEmpty
				DemonCharacter(obj.Pos, tile)
			case data.BlockFly:
				tile.Block = data.BlockEmpty
				FlyCharacter(obj.Pos, tile)
			case data.BlockGemYellow,
				data.BlockGemOrange,
				data.BlockGemGray,
				data.BlockGemCyan,
				data.BlockGemBlue,
				data.BlockGemGreen,
				data.BlockGemPurple,
				data.BlockGemBrown:
				key := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateGem(obj.Pos, key)
			case data.BlockDoorYellow,
				data.BlockDoorOrange,
				data.BlockDoorGray,
				data.BlockDoorCyan,
				data.BlockDoorBlue,
				data.BlockDoorGreen,
				data.BlockDoorPurple,
				data.BlockDoorBrown,
				data.BlockClosedYellow,
				data.BlockClosedOrange,
				data.BlockClosedGray,
				data.BlockClosedCyan,
				data.BlockClosedBlue,
				data.BlockClosedGreen,
				data.BlockClosedPurple,
				data.BlockClosedBrown,
				data.BlockLockYellow,
				data.BlockLockOrange,
				data.BlockLockGray,
				data.BlockLockCyan,
				data.BlockLockBlue,
				data.BlockLockGreen,
				data.BlockLockPurple,
				data.BlockLockBrown:
				doorKey := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateDoor(obj.Pos, doorKey)
			case data.BlockBox:
				tile.Block = data.BlockEmpty
				CreateBox(obj.Pos, tile)
			case data.BlockKeyYellow,
				data.BlockKeyOrange,
				data.BlockKeyGray,
				data.BlockKeyCyan,
				data.BlockKeyBlue,
				data.BlockKeyGreen,
				data.BlockKeyPurple,
				data.BlockKeyBrown:
				keyKey := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateKey(obj.Pos, keyKey)
			case data.BlockBomb:
				key := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateBomb(obj.Pos, key, tile.Metadata, tile.Coords)
			case data.BlockBombLit:
				key := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateLitBomb(obj.Pos, key, tile.Metadata)
			case data.BlockJetpack:
				tile.Block = data.BlockEmpty
				CreateJetpack(obj.Pos, tile.Metadata, tile.Coords)
			case data.BlockDisguise:
				tile.Block = data.BlockEmpty
				CreateDisguise(obj.Pos, tile.Metadata, tile.Coords)
			case data.BlockGear:
				var a *reanimator.Anim
				if (tile.Coords.X+tile.Coords.Y)%2 == 0 {
					a = reanimator.NewBatchAnimationCustom("gear", img.Batchers[constants.TileBatch], "gear", []int{3, 0, 1, 2}, reanimator.Loop).Reverse()
				} else {
					a = reanimator.NewBatchAnimation("gear", img.Batchers[constants.TileBatch], "gear", reanimator.Loop)
				}
				anim := reanimator.NewSimple(a)
				tile.Entity.AddComponent(myecs.Drawable, anim)
				tile.Entity.AddComponent(myecs.Animated, anim)
			}
		}
	}
	CreateFakePlayer()
	PuzzleViewInit()
	for p := 0; p < data.CurrPuzzleSet.Metadata.NumPlayers; p++ {
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
	}
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
		x := random.Level.Intn(constants.PuzzleWidth)
		y := random.Level.Intn(constants.PuzzleHeight)
		tile = data.CurrLevel.Tiles.Get(x, y)
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
		tile := data.CurrLevel.Tiles.Get(c.X, c.Y)
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
