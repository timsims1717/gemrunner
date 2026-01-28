package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"golang.org/x/image/colornames"
	"math"
)

func DrawBorder(obj *object.Object, bord *ui.Border, target pixel.Target, level bool) {
	if bord == nil || bord.Hidden {
		return
	}
	switch bord.Style {
	case ui.FancyBorder:
		DrawFancyBorder(bord, obj, level)
	case ui.ThinBorder:
		DrawThinBorder(bord, obj)
	case ui.ThinBorderReverse:
		DrawThinBorderReverse(bord, obj)
	case ui.ThinBorderWhite:
		DrawThinBorderWhite(bord, obj)
	case ui.ThinBorderBlue:
		DrawThinBorderBlue(bord, obj)
	case ui.ThickBorder:
		DrawThickBorder(bord, obj)
	case ui.ThickBorderReverse:
		DrawThickBorderReverse(bord, obj)
	case ui.ThickBorderWhite:
		DrawThickBorderWhite(bord, obj)
	case ui.ThickBorderBlue:
		DrawThickBorderBlue(bord, obj)
	}
	img.Batchers[constants.UIBatch].Draw(target)
	img.Clear()
}

func DrawFancyBorder(bord *ui.Border, obj *object.Object, level bool) {
	for y := 0; y < bord.Height+2; y++ {
		if y == 0 || y == bord.Height+1 {
			for x := 0; x < bord.Width+2; x++ {
				DrawFancyBorderSection(x, y, bord, obj, level)
				//if !bord.Empty && y == bord.Height && x != 0 {
				//	DrawBlackSquare(x, y, bord, obj)
				//}
			}
		} else {
			for x := 0; x < bord.Width+2; x++ {
				if x == 0 || x == bord.Width+1 {
					DrawFancyBorderSection(x, y, bord, obj, level)
					//if !bord.Empty && x == bord.Width {
					//	DrawBlackSquare(x, y, bord, obj)
					//}
				} else if !bord.Empty {
					DrawBlackSquare(x, y, bord, obj)
				}
			}
		}
	}
}

func DrawFancyBorderSection(x, y int, bord *ui.Border, obj *object.Object, level bool) {
	if level && data.CurrLevel != nil {
		if x == 0 && y > 0 && y <= bord.Height {
			if t := data.CurrLevel.Get(x, y-1); !t.IsSolid() {
				if trans, ok := t.Transitions[data.Left]; ok && (trans.Complete || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		} else if x == bord.Width+1 && y > 0 && y <= bord.Height {
			if t := data.CurrLevel.Get(x-2, y-1); !t.IsSolid() {
				if trans, ok := t.Transitions[data.Right]; ok && (trans.Complete || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		} else if y == 0 && x > 0 && x <= bord.Width {
			if t := data.CurrLevel.Get(x-1, y); !t.IsSolid() {
				if trans, ok := t.Transitions[data.Down]; ok && (trans.Complete || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		} else if y == bord.Height+1 && x > 0 && x <= bord.Width {
			if t := data.CurrLevel.Get(x-1, y-2); !t.IsSolid() {
				if trans, ok := t.Transitions[data.Up]; ok && (trans.Complete || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		}
	}
	mat := pixel.IM
	sKey := constants.FancyBorderStraight
	if (x == 0 || x == bord.Width+1) && (y == 0 || y == bord.Height+1) {
		sKey = constants.FancyBorderCorner
	}
	offset := pixel.V(world.TileSize*(float64(x)-float64(bord.Width+1)*0.5), world.TileSize*(float64(y)-float64(bord.Height+1)*0.5))
	if x == 0 {
		offset.X += world.TileSize * 0.25
	} else if x == bord.Width+1 {
		offset.X -= world.TileSize * 0.25
	}
	if y == 0 {
		offset.Y += world.TileSize * 0.25
	} else if y == bord.Height+1 {
		offset.Y -= world.TileSize * 0.25
	}
	if y == 0 {
		if x > 0 && x <= bord.Width {
			mat = mat.Rotated(pixel.ZV, 0.5*math.Pi)
		} else if x == bord.Width+1 {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
		}
	} else if y == bord.Height+1 {
		if x > 0 && x <= bord.Width {
			mat = mat.Rotated(pixel.ZV, -0.5*math.Pi)
		} else if x == 0 {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(1., -1.))
		} else if x == bord.Width+1 {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(-1., -1.))
		}
	} else if x == bord.Width+1 {
		mat = mat.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
	}
	img.Batchers[constants.UIBatch].DrawSpriteColor(sKey, mat.Moved(obj.PostPos).Moved(offset), colornames.White)
}

func DrawThinBorder(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
}

func DrawThinBorderReverse(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
}

func DrawThinBorderWhite(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
}

func DrawThinBorderBlue(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
}

func DrawThickBorder(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	matTB2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+4, 1.))
	matLR2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+4))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+1.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+1.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-1.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-1.5, 0)))
}

func DrawThickBorderReverse(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	matTB2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+4, 1.))
	matLR2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+4))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+1.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+1.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-1.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-1.5, 0)))
}

func DrawThickBorderWhite(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	matTB2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+4, 1.))
	matLR2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+4))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+1.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+1.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-1.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderWhite, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-1.5, 0)))
}

func DrawThickBorderBlue(bord *ui.Border, obj *object.Object) {
	matTB := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+2, 1.))
	matLR := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+2))
	matTB2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(bord.Rect.W()+4, 1.))
	matLR2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., bord.Rect.H()+4))
	// top
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*0.5+1.5)))
	// right
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*0.5+1.5, 0)))
	// bottom
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-0.5)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matTB2.Moved(obj.PostPos).Moved(pixel.V(0, bord.Rect.H()*-0.5-1.5)))
	// left
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-0.5, 0)))
	img.Batchers[constants.UIBatch].DrawSprite(constants.ThinBorderBlue, matLR2.Moved(obj.PostPos).Moved(pixel.V(bord.Rect.W()*-0.5-1.5, 0)))
}

func DrawBlackSquare(x, y int, bord *ui.Border, obj *object.Object) {
	mat := pixel.IM
	offset := pixel.V(world.TileSize*(float64(x)-float64(bord.Width+1)*0.5), world.TileSize*(float64(y)-float64(bord.Height+1)*0.5))
	sKey := "black_square_16"
	img.Batchers[constants.UIBatch].DrawSpriteColor(sKey, mat.Moved(obj.PostPos).Moved(offset), colornames.White)
}

func InitMainBorder() {
	if ui.PuzzleBorder == nil {
		ui.PuzzleBorderObject = object.New()
		ui.PuzzleBorderObject.Pos.X = world.TileSize * 0.5 * constants.PuzzleWidth
		ui.PuzzleBorderObject.Pos.Y = world.TileSize * 0.5 * constants.PuzzleHeight
		ui.PuzzleBorderObject.Layer = 1
		ui.PuzzleBorder = &ui.Border{
			Width:  constants.PuzzleWidth,
			Height: constants.PuzzleHeight,
			Empty:  false,
		}
		ui.MainBorder = myecs.Manager.NewEntity()
		ui.MainBorder.AddComponent(myecs.Object, ui.PuzzleBorderObject).
			AddComponent(myecs.Border, ui.PuzzleBorder)
	}
}

func SetMainBorder(w, h int) {
	if ui.PuzzleBorder != nil {
		ui.PuzzleBorderObject.Pos.X = world.TileSize * 0.5 * float64(w)
		ui.PuzzleBorderObject.Pos.Y = world.TileSize * 0.5 * float64(h)
		ui.PuzzleBorder.Width = w
		ui.PuzzleBorder.Height = h
	}
}
