package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
)

func EditorInit() {
	// initialize editor panel
	data.NewEditorPane()
	data.EditorPanel.ViewPort = viewport.New(nil)
	data.EditorPanel.ViewPort.SetILock(true)
	data.EditorPanel.ViewPort.SetRect(pixel.R(0, 0, world.TileSize*3., world.TileSize*8.))
	data.EditorPanel.ViewPort.CamPos = pixel.V(world.TileSize*-2.-3., world.TileSize*4.5+3.)

	// border
	borderObj := object.New()
	borderObj.Pos = pixel.V(world.TileSize*-2., world.TileSize*4.5)
	borderObj.SetRect(pixel.R(0, 0, world.TileSize*3, world.TileSize*8))
	borderObj.Layer = 3

	// the viewport for the block selectors
	vp := viewport.New(nil)
	vp.SetRect(pixel.R(0, 0, world.TileSize*6.+2., world.TileSize*3.+2))
	vp.CamPos = pixel.V(world.TileSize*3., -world.TileSize*1.5)
	data.EditorPanel.BlockSelect = vp

	// the block selectors
	b := 0
	for ; b < data.Empty; b++ {
		obj := object.New()
		obj.Pos = data.BlockSelectPlacement(b)
		fmt.Printf("Tile %d: (%d,%d)\n", b, int(obj.Pos.X), int(obj.Pos.Y))
		obj.Layer = 4
		obj.SetRect(pixel.R(0., 0., 16., 16.))
		sprB := img.NewSprite("black_square_big", constants.UIBatch)
		spr := img.NewSprite(data.Block(b).String(), constants.TileBGBatch)
		bId := data.Block(b)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, []interface{}{sprB, spr}).
			AddComponent(myecs.Block, data.Block(b)).
			AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, data.EditorPanel.BlockSelect, func(hvc *data.HoverClick) {
				click := hvc.Input.Get("click")
				if hvc.Hover && data.EditorPanel.SelectVis {
					data.EditorPanel.SelectObj.Pos = obj.Pos
					if click.JustPressed() || click.JustReleased() {
						data.EditorPanel.CurrBlock = bId
						data.EditorPanel.SelectVis = false
						data.EditorPanel.SelectQuick = false
						data.EditorPanel.SelectTimer = nil
						data.EditorPanel.Consume = ""
						click.Consume()
					}
				}
			})).
			AddComponent(myecs.ViewPort, data.EditorPanel.BlockSelect)
	}
	objOutline := object.New()
	objOutline.Pos = data.BlockSelectPlacement(0)
	objOutline.Layer = 5
	sprO := img.NewSprite("white_outline", constants.UIBatch)
	myecs.Manager.NewEntity().
		AddComponent(myecs.Object, objOutline).
		AddComponent(myecs.Drawable, sprO)
	data.EditorPanel.SelectObj = objOutline
	for ; b < 18; b++ {
		obj := object.New()
		obj.Pos = data.BlockSelectPlacement(b)
		obj.Layer = 4
		spr := img.NewSprite("black_square_big", constants.UIBatch)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, spr)
	}

	// block select
	blockObj := object.New()
	blockObj.Pos = borderObj.Pos
	blockObj.Pos.Y -= world.TileSize * 2.5
	blockObj.Layer = 3
	blockObj.Rect = pixel.R(-16., -16., 16., 16.)
	beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
	beFG := img.NewSprite(data.Block(data.RedRock).String(), constants.TileBGBatch)
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, blockObj).
		AddComponent(myecs.Drawable, []interface{}{beBG, beFG}).
		AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, data.EditorPanel.ViewPort, func(hvc *data.HoverClick) {
			beFG.Key = data.EditorPanel.CurrBlock.String()
			data.EditorPanel.Hover = hvc.Hover
			click := hvc.Input.Get("click")
			if hvc.Hover && (data.EditorPanel.Consume == "select" || data.EditorPanel.Consume == "") {
				if data.EditorPanel.Consume == "" {
					data.EditorPanel.SelectQuick = false
					if click.JustPressed() {
						data.EditorPanel.Consume = "select"
						data.EditorPanel.SelectTimer = timing.New(0.2)
					}
				} else if data.EditorPanel.Consume == "select" {
					if click.JustPressed() {
						data.EditorPanel.Consume = ""
						data.EditorPanel.SelectTimer = nil
						data.EditorPanel.SelectQuick = false
					} else if click.JustReleased() {
						if data.EditorPanel.SelectTimer != nil && !data.EditorPanel.SelectTimer.Done() {
							data.EditorPanel.SelectQuick = true
							data.EditorPanel.Consume = "select"
						}
					} else if !click.Pressed() && !data.EditorPanel.SelectQuick {
						data.EditorPanel.Consume = ""
						data.EditorPanel.SelectTimer = nil
					}
				}
				data.EditorPanel.SelectVis = data.EditorPanel.Consume == "select"
			}
		}))

	// border and editor panel movement
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Border, &data.Border{
		Width:  2,
		Height: 7,
		Empty:  false,
	}).
		AddComponent(myecs.Object, borderObj).
		AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, data.EditorPanel.ViewPort, func(hvc *data.HoverClick) {
			data.EditorPanel.Hover = hvc.Hover
			click := hvc.Input.Get("click")
			if hvc.Hover {
				if data.EditorPanel.Consume == "move" || data.EditorPanel.Consume == "" {
					if click.JustPressed() {
						data.EditorPanel.Offset = data.EditorPanel.ViewPort.PostPortPos.Sub(hvc.Input.World)
						data.EditorPanel.Consume = "move"
					}
				}
			}
			if click.JustReleased() && data.EditorPanel.Offset != pixel.ZV {
				data.EditorPanel.ViewPort.PortPos = hvc.Input.World.Add(data.EditorPanel.Offset)
				data.EditorPanel.Offset = pixel.ZV
				data.EditorPanel.Consume = ""
			}
		}))
	data.EditorPanel.Entity = e
	data.EditorPanel.BlockView = &data.BlockView{
		Entity: be,
		Object: blockObj,
	}
}

func PuzzleInit() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		myecs.Manager.DisposeEntity(result)
	}
	if data.CurrPuzzle != nil {
		for _, row := range data.CurrPuzzle.Tiles {
			for _, tile := range row {
				obj := object.New()
				obj.Pos = world.MapToWorld(tile.Coords)
				obj.Pos.X += world.TileSize * 0.5
				obj.Pos.Y += world.TileSize * 0.5
				obj.Layer = 2
				myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
			}
		}
		TileUpdate = true
	}

	data.PuzzleView = viewport.New(nil)
	data.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
	data.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
	data.PuzzleView.PortPos = viewport.MainCamera.CamPos
	data.BorderView = viewport.New(nil)
	data.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(constants.PuzzleWidth+1), world.TileSize*(constants.PuzzleHeight+1)))
	data.BorderView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
}

func PuzzleEditSystem() {
	if data.EditorPanel.SelectTimer != nil {
		data.EditorPanel.SelectTimer.Update()
	}
	if data.EditorPanel.Consume != "select" {
		data.EditorPanel.SelectObj.Pos = data.BlockSelectPlacement(int(data.EditorPanel.CurrBlock))
	}
	if data.EditorPanel.Consume != "" {
		switch data.EditorPanel.Consume {
		case "move":
			if data.EditorPanel.Offset.X != 0 || data.EditorPanel.Offset.Y != 0 {
				data.EditorPanel.ViewPort.PortPos = data.EditorInput.World.Add(data.EditorPanel.Offset)
			}
		}
	} else if !data.EditorPanel.Hover {
		//data.EditorPanel.SelectVis = false
		projPos := data.PuzzleView.ProjectWorld(data.EditorInput.World)
		if data.EditorInput.Get("rightClick").Pressed() {
			x, y := world.WorldToMap(projPos.X, projPos.Y)
			coords := world.Coords{X: x, Y: y}
			DeleteBlock(coords)
		} else if data.EditorInput.Get("click").Pressed() {
			x, y := world.WorldToMap(projPos.X, projPos.Y)
			coords := world.Coords{X: x, Y: y}
			ChangeBlock(coords, data.EditorPanel.CurrBlock)
		}
	}
}
