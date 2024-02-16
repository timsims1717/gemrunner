package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
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
			obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
			obj.Layer = 10
			tile.Object = obj
			e := myecs.Manager.NewEntity().
				AddComponent(myecs.Object, obj).
				AddComponent(myecs.Tile, tile)
			tile.Entity = e
			// replace reanimator and items
			switch tile.Block {
			case data.BlockPlayer1:
				tile.Block = data.BlockEmpty
				p1 := PlayerCharacter(obj.Pos, 0)
				data.CurrLevel.Players[0] = p1
				data.CurrLevel.Chars = append(data.CurrLevel.Chars, p1)
				data.CurrLevel.Stats[0] = data.NewStats()
			case data.BlockDemon:
				tile.Block = data.BlockEmpty
				demon := DemonCharacter(obj.Pos)
				//data.CurrLevel.Players[1] = demon
				data.CurrLevel.Chars = append(data.CurrLevel.Chars, demon)
			case data.BlockFly:
				tile.Block = data.BlockEmpty
				fly := FlyCharacter(obj.Pos, tile.Metadata.Flipped)
				data.CurrLevel.Chars = append(data.CurrLevel.Chars, fly)
			case data.BlockGem:
				tile.Block = data.BlockEmpty
				CreateGem(obj.Pos)
			case data.BlockDoorPink, data.BlockDoorBlue, data.BlockLockPink, data.BlockLockBlue:
				doorKey := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateDoor(obj.Pos, doorKey)
			case data.BlockBox:
				tile.Block = data.BlockEmpty
				CreateBox(obj.Pos)
			case data.BlockKeyPink, data.BlockKeyBlue:
				keyKey := tile.Block.String()
				tile.Block = data.BlockEmpty
				CreateKey(obj.Pos, keyKey)
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
