package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/pkg/errors"
)

func PuzzleInit() {
	PuzzleDispose()
	if data.CurrPuzzleSet.CurrPuzzle != nil {
		for y, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
			for x, tile := range row {
				tile.Coords = world.NewCoords(x, y)
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
				if tile.TextData != nil {
					data.CreateFloatingText(tile, tile.TextData)
				}
			}
		}
		num := data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber
		if data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite == "" {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite = constants.WorldSprites[num]
		}
		// todo: check world colors (all black, I assume)
		if data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack == "" {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack = constants.WorldMusic[num]
		}
		data.SelectedWorldIndex = data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber
		if data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber == constants.WorldCustom {
			for n, w := range constants.WorldSprites {
				if data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite == w {
					data.SelectedWorldIndex = n
				}
			}
		}
		data.CurrPuzzleSet.CurrPuzzle.Update = true
	} else {
		panic("no puzzle loaded")
	}
	if data.Editor != nil {
		num := ui.Dialogs[constants.DialogEditorOptionsRight].Get("puzzle_number")
		num.Text.SetText(fmt.Sprintf("%04d", data.CurrPuzzleSet.PuzzleIndex+1))
	}
	PuzzleViewInit()
	UpdateEditorShaders()
	UpdatePuzzleShaders()
	ChangeWorldShader(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode)
}

func PuzzleViewInit() {
	wRatio := viewport.MainCamera.Rect.W() / (constants.PuzzleWidth * world.TileSize)
	hRatio := viewport.MainCamera.Rect.H() / (constants.PuzzleHeight * world.TileSize)
	maxRatio := wRatio
	if hRatio < wRatio {
		maxRatio = hRatio
	}
	maxRatio *= constants.ScreenRatioLimit

	constants.PickedRatio = 1.
	for constants.PickedRatio+1 < maxRatio {
		constants.PickedRatio += 1
	}
	if data.PuzzleView == nil {
		data.PuzzleView = viewport.New(nil)
		data.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
		data.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
		data.PuzzleView.PortPos = viewport.MainCamera.CamPos
		UpdatePuzzleShaders()
		data.PuzzleView.Canvas.SetFragmentShader(data.PuzzleShader)
	}
	if data.PuzzleViewNoShader == nil {
		data.PuzzleViewNoShader = viewport.New(nil)
		data.PuzzleViewNoShader.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
		data.PuzzleViewNoShader.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
		data.PuzzleViewNoShader.PortPos = viewport.MainCamera.CamPos
	}
	if data.BorderView == nil {
		data.BorderView = viewport.New(nil)
		data.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(constants.PuzzleWidth+1), world.TileSize*(constants.PuzzleHeight+1)))
		data.BorderView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
	}
	if data.WorldView == nil {
		data.WorldView = viewport.New(nil)
		xWidth := constants.PuzzleWidth * world.TileSize * constants.PickedRatio
		yHeight := constants.PuzzleHeight * world.TileSize * constants.PickedRatio
		data.WorldView.SetRect(pixel.R(0, 0, xWidth, yHeight))
		data.WorldView.CamPos = viewport.MainCamera.CamPos
		data.WorldView.PortPos = viewport.MainCamera.CamPos
		ChangeWorldShader(0)
		data.WorldView.Canvas.SetFragmentShader(data.WorldShader)
	}
}

func PuzzleDispose() {
	for _, result := range myecs.Manager.Query(myecs.IsText) {
		if ft, ok := result.Components[myecs.Text].(*data.FloatingText); ok {
			myecs.Manager.DisposeEntity(ft.Entity)
			myecs.Manager.DisposeEntity(ft.ShEntity)
			data.RemoveFloatingText(ft.Tile)
		}
	}
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		myecs.Manager.DisposeEntity(result.Entity)
	}
}

func UpdatePuzzleShaders() {
	// set puzzle shader uniforms
	data.PuzzleView.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.B))
	data.PuzzleView.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.B))
	data.PuzzleView.Canvas.SetUniform("uRedDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.G))
	data.PuzzleView.Canvas.SetUniform("uBlueDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.B))
	data.PuzzleView.Canvas.SetUniform("uMode", int32(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode))
	data.PuzzleView.Canvas.SetUniform("uTime", &data.CurrPuzzleSet.Elapsed)
	data.PuzzleView.Canvas.SetUniform("uSpeed", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderSpeed)
	data.PuzzleView.Canvas.SetUniform("uXVar", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderX)
	data.PuzzleView.Canvas.SetUniform("uYVar", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderY)
	data.PuzzleView.Canvas.SetUniform("uCustom", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderCustom)
	data.PuzzleView.Canvas.SetUniform("uParticle", int32(data.CurrPuzzleSet.CurrPuzzle.Metadata.ParticleMode))
}

func UpdateWorldShaders() {
	data.WorldView.Canvas.SetUniform("uMode", int32(data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode))
	data.WorldView.Canvas.SetUniform("uTime", &data.CurrPuzzleSet.Elapsed)
	data.WorldView.Canvas.SetUniform("uSpeed", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderSpeed)
	data.WorldView.Canvas.SetUniform("uXVar", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderX)
	data.WorldView.Canvas.SetUniform("uYVar", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderY)
	data.WorldView.Canvas.SetUniform("uCustom", &data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderCustom)
	darkness := int32(0)
	if data.CurrPuzzleSet.CurrPuzzle.Metadata.Darkness {
		darkness = 1
	}
	data.WorldView.Canvas.SetUniform("uDarkness", darkness)
}

func NewPuzzleSet() {
	PuzzleDispose()
	data.CurrPuzzleSet = data.CreatePuzzleSet()
	data.CurrPuzzleSet.SetToFirst()
	PuzzleInit()
}

func AddPuzzle() {
	PuzzleDispose()
	data.CurrPuzzleSet.AppendNew()
	PuzzleInit()
}

func PrevPuzzle() {
	PuzzleDispose()
	data.CurrPuzzleSet.Prev()
	PuzzleInit()
}

func NextPuzzle() {
	PuzzleDispose()
	data.CurrPuzzleSet.Next()
	PuzzleInit()
}

func DeletePuzzle() {
	PuzzleDispose()
	data.CurrPuzzleSet.Delete(data.CurrPuzzleSet.PuzzleIndex)
	PuzzleInit()
}

func SavePuzzleSet() bool {
	if data.CurrPuzzleSet != nil {
		if data.CurrPuzzleSet.Metadata.Name == "" {
			fmt.Println("ERROR: puzzle set has no name")
			return false
		}
		if data.CurrPuzzleSet.Metadata.Filename == "" {
			data.CurrPuzzleSet.Metadata.Filename = fmt.Sprintf("%s.puzzle", data.CurrPuzzleSet.Metadata.Name)
		}
		err := content.SavePuzzleSetToFile()
		if err != nil {
			fmt.Println("ERROR:", err)
			return false
		}
		data.CurrPuzzleSet.NeedToSave = false
		for _, pzl := range data.CurrPuzzleSet.Puzzles {
			pzl.Changed = false
		}
		return true
	} else {
		fmt.Println("ERROR: no puzzle set to save")
		return false
	}
}

func OpenPuzzleSet(filename string) error {
	PuzzleDispose()
	err := content.OpenPuzzleSetFile(filename)
	if err != nil {
		return err
	}
	data.CurrPuzzleSet.SetToFirst()
	if data.CurrPuzzleSet.Metadata.NumPlayers < 1 {
		data.CurrPuzzleSet.Metadata.NumPlayers = data.CurrPuzzleSet.CurrPuzzle.NumPlayers()
	}
	PuzzleInit()
	return nil
}

func CombinePuzzleSet(filename string) error {
	pzlSet, err := content.OpenPuzzleSetFileRt(filename)
	if err != nil {
		return err
	} else if pzlSet == nil {
		return errors.New("no puzzle set to combine")
	}
	oIndex := data.CurrPuzzleSet.PuzzleIndex + 1
	for _, pzl := range pzlSet.Puzzles {
		data.CurrPuzzleSet.Insert(pzl, data.CurrPuzzleSet.PuzzleIndex+1)
	}
	data.CurrPuzzleSet.SetTo(oIndex)
	PuzzleInit()
	return nil
}
