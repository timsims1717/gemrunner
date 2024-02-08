package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
)

func LevelInit() {
	if data.CurrPuzzle == nil {
		panic("no puzzle loaded to start level")
	}
	LevelDispose()
	data.CurrLevel = &data.Level{}
	data.CurrLevel.Tiles = data.CurrPuzzle.CopyTiles()
	data.CurrLevel.Metadata = data.CurrPuzzle.Metadata
	data.CurrLevel.Puzzle = data.CurrPuzzle

	for _, row := range data.CurrLevel.Tiles.T {
		for _, tile := range row {
			obj := object.New()
			obj.Pos = world.MapToWorld(tile.Coords)
			obj.Pos.X += world.TileSize * 0.5
			obj.Pos.Y += world.TileSize * 0.5
			obj.Layer = 10
			tile.Object = obj
			e := myecs.Manager.NewEntity().
				AddComponent(myecs.Object, obj).
				AddComponent(myecs.Tile, tile)
			tile.Entity = e
			// replace characters and items
			switch tile.Block {
			case data.Player1:
				tile.Block = data.Empty
				p1 := PlayerCharacter(obj.Pos, 0)
				data.CurrLevel.Player1 = p1
				data.CurrLevel.Chars = append(data.CurrLevel.Chars, p1)
				data.CurrLevel.Stats1 = data.NewStats()
			case data.Gem:
				tile.Block = data.Empty
				CreateGem(obj.Pos)
			case data.DoorPink, data.DoorBlue:
				doorKey := tile.Block.String()
				tile.Block = data.Empty
				CreateEmptyDoor(obj.Pos, doorKey)
			}
		}
	}

	PuzzleViewInit()
}

func LevelDispose() {
	if data.CurrLevel != nil {
		for _, row := range data.CurrLevel.Tiles.T {
			for _, tile := range row {
				myecs.Manager.DisposeEntity(tile.Entity)
			}
		}
		for _, character := range data.CurrLevel.Chars {
			myecs.Manager.DisposeEntity(character.Entity)
		}
		data.CurrLevel = nil
	}
}
