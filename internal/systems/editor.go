package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
	pxginput "github.com/timsims1717/pixel-go-input"
)

func EditorInit() {
	data.NewEditorPane()
	data.EditorPane.ViewPort = viewport.New(nil)
	data.EditorPane.ViewPort.SetILock(true)
	data.EditorPane.ViewPort.SetRect(pixel.R(0, 0, world.TileSize*3., world.TileSize*8.))
	data.EditorPane.ViewPort.CamPos = pixel.V(world.TileSize*-2.-3., world.TileSize*4.5+3.)

	borderObj := object.New()
	borderObj.Pos = pixel.V(world.TileSize*-2., world.TileSize*4.5)
	borderObj.Layer = 3
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Border, &data.Border{
		Width:  2,
		Height: 7,
		Empty:  false,
	}).
		AddComponent(myecs.Object, borderObj)
	data.EditorPane.Entity = e

	blockObj := object.New()
	blockObj.Pos = borderObj.Pos
	blockObj.Pos.Y -= world.TileSize * 2.5
	blockObj.Layer = 3
	blockObj.Rect = pixel.R(-16., -16., 16., 16.)
	beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
	beFG := img.NewSprite(data.Block(data.RedRock).String(), constants.TileBGBatch)
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, blockObj).
		AddComponent(myecs.Drawable, []interface{}{beBG, beFG})
	data.EditorPane.BlockView = &data.BlockView{
		Entity: be,
		Object: blockObj,
	}

	vp := viewport.New(nil)
	vp.SetRect(pixel.R(0, 0, world.TileSize*6.+2., world.TileSize*3.+2))
	vp.CamPos = pixel.V(world.TileSize*3., -world.TileSize*1.5)
	data.EditorPane.BlockSelect = vp

	b := 0
	for ; b < data.Empty; b++ {
		obj := object.New()
		obj.Pos = data.BlockSelectPlacement(b)
		obj.Layer = 4
		obj.Rect = pixel.R(0., 0., 16., 16.)
		sprB := img.NewSprite("black_square_big", constants.UIBatch)
		spr := img.NewSprite(data.Block(b).String(), constants.TileBGBatch)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, []interface{}{sprB, spr}).
			AddComponent(myecs.Block, data.Block(b)).
			AddComponent(myecs.Hover, data.NewHoverFn(BlockSelectHover(b))).
			AddComponent(myecs.ViewPort, data.EditorPane.BlockSelect)
		if b == data.RedRock {
			objOutline := object.New()
			objOutline.Pos = data.BlockSelectPlacement(b)
			objOutline.Layer = 5
			sprO := img.NewSprite("white_outline", constants.UIBatch)
			myecs.Manager.NewEntity().
				AddComponent(myecs.Object, objOutline).
				AddComponent(myecs.Drawable, []interface{}{sprO})
			data.EditorPane.SelectObj = objOutline
		}
	}
	for ; b < 18; b++ {
		obj := object.New()
		obj.Pos = data.BlockSelectPlacement(b)
		obj.Layer = 4
		spr := img.NewSprite("black_square_big", constants.UIBatch)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, spr)
	}
}

func BlockSelectHover(b int) func(*pxginput.Input) {
	return func(in *pxginput.Input) {
		if in.Get("click").JustReleased() || in.Get("click").JustPressed() {

			data.EditorPane.SelectVis = false
		}
	}
}
