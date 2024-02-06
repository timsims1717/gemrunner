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
				obj.Pos.X += world.TileSize * 0.5
				obj.Pos.Y += world.TileSize * 0.5
				obj.Layer = 2
				myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
				tile.Object = obj
			}
		}
		data.CurrPuzzle.Update = true
	} else {
		panic("no puzzle loaded")
	}

	if data.PuzzleView == nil {
		data.PuzzleView = viewport.New(nil)
		data.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
		data.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
		data.PuzzleView.PortPos = viewport.MainCamera.CamPos
	}
	data.PuzzleView.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzle.PrimaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzle.PrimaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzle.PrimaryColor.B))
	data.PuzzleView.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzle.SecondaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzle.SecondaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzle.SecondaryColor.B))
	data.PuzzleView.Canvas.SetFragmentShader(data.PuzzleShader)
	if data.BorderView == nil {
		data.BorderView = viewport.New(nil)
		data.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(constants.PuzzleWidth+1), world.TileSize*(constants.PuzzleHeight+1)))
		data.BorderView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
	}
}

func PuzzleDispose(full bool) {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		myecs.Manager.DisposeEntity(result)
	}
	if full {
		data.CurrPuzzle = nil
	}
}

func SavePuzzle() error {
	errMsg := "save puzzle"
	if data.CurrPuzzle == nil {
		return errors.Wrap(errors.New("no puzzle to save"), errMsg)
	}
	var svgFName = "test.puzzle"
	if data.CurrPuzzle.PuzzleInfo.Filename != "" {
		svgFName = data.CurrPuzzle.PuzzleInfo.Filename
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
	err = json.Unmarshal(pzlFile, data.CurrPuzzle)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	PuzzleInit()
	fmt.Printf("INFO: loaded puzzle from %s\n", filename)
	return nil
}
