package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"golang.org/x/image/colornames"
	"math"
)

func BorderSystem(layer int) {
	for _, result := range myecs.Manager.Query(myecs.HasBorder) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		bord, okB := result.Components[myecs.Border].(*data.Border)
		if okO && okB && obj.Layer == layer {
			switch bord.Style {
			case data.FancyBorder:
				DrawFancyBorder(bord, obj)
			case data.ThinBorder:
				DrawThinBorder(bord, obj)
			}
		}
	}
}

func DrawFancyBorder(bord *data.Border, obj *object.Object) {
	for y := 0; y < bord.Height+1; y++ {
		if y == 0 || y == bord.Height {
			for x := 0; x < bord.Width+1; x++ {
				DrawFancyBorderSection(x, y, bord, obj)
				if !bord.Empty && y == bord.Height && x != 0 {
					DrawBlackSquare(x, y, bord, obj)
				}
			}
		} else {
			for x := 0; x < bord.Width+1; x++ {
				if x == 0 || x == bord.Width {
					DrawFancyBorderSection(x, y, bord, obj)
					if !bord.Empty && x == bord.Width {
						DrawBlackSquare(x, y, bord, obj)
					}
				} else if !bord.Empty {
					DrawBlackSquare(x, y, bord, obj)
				}
			}
		}
	}
}

func DrawFancyBorderSection(x, y int, bord *data.Border, obj *object.Object) {
	mat := pixel.IM
	offset := pixel.V(world.TileSize*(float64(x)-float64(bord.Width)*0.5), world.TileSize*(float64(y)-float64(bord.Height)*0.5))
	sKey := constants.FancyBorderStraight
	if (x == 0 || x == bord.Width) && (y == 0 || y == bord.Height) {
		sKey = constants.FancyBorderCorner
	}
	if y == 0 {
		if x > 0 && x < bord.Width {
			mat = mat.Rotated(pixel.ZV, 0.5*math.Pi)
		} else if x == bord.Width {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
		}
	} else if y == bord.Height {
		if x > 0 && x < bord.Width {
			mat = mat.Rotated(pixel.ZV, -0.5*math.Pi)
		} else if x == 0 {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
		} else if x == bord.Width {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(-1., -1.))
		}
	} else if x == bord.Width {
		mat = mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
	}
	img.Batchers[constants.UIBatch].DrawSpriteColor(sKey, mat.Moved(obj.PostPos).Moved(offset), colornames.White)
}

func DrawThinBorder(bord *data.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(world.TileSize*float64(bord.Width)+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, world.TileSize*float64(bord.Height)+2))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, world.TileSize*float64(bord.Height)*0.5+1)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(world.TileSize*float64(bord.Width)*0.5+0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, world.TileSize*float64(bord.Height)*-0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(world.TileSize*float64(bord.Width)*-0.5-0.5, 0)))
}

// todo: do this later
func DrawThinBlackSquare(x, y int, bord *data.Border, obj *object.Object) {
	mat := pixel.IM
	offset := pixel.V(world.TileSize*(float64(x)-float64(bord.Width+1)*0.5), world.TileSize*(float64(y)-float64(bord.Height+1)*0.5))
	sKey := "black_square"
	img.Batchers[constants.UIBatch].DrawSpriteColor(sKey, mat.Moved(obj.PostPos).Moved(offset), colornames.White)
}

func DrawBlackSquare(x, y int, bord *data.Border, obj *object.Object) {
	mat := pixel.IM
	offset := pixel.V(world.TileSize*(float64(x)-float64(bord.Width+1)*0.5), world.TileSize*(float64(y)-float64(bord.Height+1)*0.5))
	sKey := "black_square"
	img.Batchers[constants.UIBatch].DrawSpriteColor(sKey, mat.Moved(obj.PostPos).Moved(offset), colornames.White)
}

func InitMainBorder() {
	if data.MainBorder == nil {
		borderObj := object.New()
		borderObj.Pos.X = world.TileSize * 0.5 * constants.PuzzleWidth
		borderObj.Pos.Y = world.TileSize * 0.5 * constants.PuzzleHeight
		borderObj.Layer = 1
		data.MainBorder = myecs.Manager.NewEntity()
		data.MainBorder.AddComponent(myecs.Object, borderObj).
			AddComponent(myecs.Border, &data.Border{
				Width:  constants.PuzzleWidth,
				Height: constants.PuzzleHeight,
				Empty:  false,
			})
	}
}
