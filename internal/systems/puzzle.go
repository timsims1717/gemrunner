package systems

import (
	"encoding/json"
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/pkg/errors"
	"os"
)

func PuzzleInit() {
	PuzzleDispose(false)
	if data.CurrPuzzle != nil {
		for _, row := range data.CurrPuzzle.Tiles.T {
			for _, tile := range row {
				obj := object.New()
				obj.Pos = world.MapToWorld(tile.Coords)
				obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
				obj.Flip = tile.Metadata.Flipped
				obj.Layer = 2
				e := myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
				tile.Object = obj
				tile.Entity = e
			}
		}
		data.CurrPuzzle.Update = true
	} else {
		panic("no puzzle loaded")
	}
	PuzzleViewInit()
}

func PuzzleViewInit() {
	if data.PuzzleView == nil {
		data.PuzzleView = viewport.New(nil)
		data.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
		data.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
		data.PuzzleView.PortPos = viewport.MainCamera.CamPos
	}
	data.PuzzleView.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzle.Metadata.PrimaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzle.Metadata.PrimaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzle.Metadata.PrimaryColor.B))
	data.PuzzleView.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzle.Metadata.SecondaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzle.Metadata.SecondaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzle.Metadata.SecondaryColor.B))
	data.PuzzleView.Canvas.SetFragmentShader(data.PuzzleShader)
	if data.BorderView == nil {
		data.BorderView = viewport.New(nil)
		data.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(constants.PuzzleWidth+1), world.TileSize*(constants.PuzzleHeight+1)))
		data.BorderView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
	}
}

func PuzzleDispose(full bool) {
	if data.CurrPuzzle != nil {
		for _, row := range data.CurrPuzzle.Tiles.T {
			for _, tile := range row {
				myecs.Manager.DisposeEntity(tile.Entity)
			}
		}
		if full {
			data.CurrPuzzle = nil
		}
	}
}

func NewPuzzle() {
	if data.CurrPuzzle != nil {
		PuzzleDispose(true)
	}
	data.CurrPuzzle = data.CreateBlankPuzzle()
	PuzzleInit()
	UpdateWorldShaders()
}

func SavePuzzle() {
	if data.Editor != nil {
		if data.CurrPuzzle.Changed {
			if data.CurrPuzzle.Metadata == nil {
				data.CurrPuzzle.Metadata = &data.PuzzleMetadata{}
			}
			if data.CurrPuzzle.Metadata.Name == "" {
				data.CurrPuzzle.Metadata.Name = "test"
				data.CurrPuzzle.Metadata.Filename = "test.puzzle"
			}
			data.CurrPuzzle.Metadata.Filename = fmt.Sprintf("%s.puzzle", data.CurrPuzzle.Metadata.Name)
			err := SavePuzzleToFile()
			if err != nil {
				fmt.Println("ERROR:", err)
			} else {
				data.CurrPuzzle.Changed = false
			}
		}
	}
}

func SavePuzzleToFile() error {
	errMsg := "save puzzle"
	if data.CurrPuzzle == nil {
		return errors.Wrap(errors.New("no puzzle to save"), errMsg)
	}
	var svgFName = "test.puzzle"
	if data.CurrPuzzle.Metadata.Filename != "" {
		svgFName = data.CurrPuzzle.Metadata.Filename
	}
	svgPath := fmt.Sprintf("%s/%s", constants.PuzzlesDir, svgFName)
	saveFile, err := os.Create(svgPath)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	bts, err := json.Marshal(data.CurrPuzzle)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	_, err = saveFile.Write(bts)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	fmt.Printf("INFO: saved puzzle to %s\n", svgPath)
	return nil
}

func OpenPuzzle(filename string) error {
	errMsg := "open puzzle"
	if filename == "" {
		filename = "assets/save.savegame"
	}
	pzlFile, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	if data.CurrPuzzle == nil {
		data.CurrPuzzle = data.CreateBlankPuzzle()
	}
	err = json.Unmarshal(pzlFile, data.CurrPuzzle)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	PuzzleInit()
	UpdateWorldShaders()
	fmt.Printf("INFO: loaded puzzle from %s\n", filename)
	return nil
}
