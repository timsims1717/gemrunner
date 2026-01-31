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

func DrawBorder(obj *object.Object, bord *data.Border, target pixel.Target, level bool) {
	if bord == nil || bord.Hidden {
		return
	}
	switch bord.Style {
	case data.FancyBorder:
		DrawFancyBorder(bord, obj, level)
	case data.ThinBorder:
		DrawThinBorder(bord, obj)
	case data.ThinBorderReverse:
		DrawThinBorderReverse(bord, obj)
	case data.ThinBorderWhite:
		DrawThinBorderWhite(bord, obj)
	case data.ThinBorderBlue:
		DrawThinBorderBlue(bord, obj)
	case data.ThickBorder:
		DrawThickBorder(bord, obj)
	case data.ThickBorderReverse:
		DrawThickBorderReverse(bord, obj)
	case data.ThickBorderWhite:
		DrawThickBorderWhite(bord, obj)
	case data.ThickBorderBlue:
		DrawThickBorderBlue(bord, obj)
	}
	img.Batchers[constants.UIBatch].Draw(target)
	img.Clear()
}

func DrawFancyBorder(bord *data.Border, obj *object.Object, level bool) {
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

func DrawFancyBorderSection(x, y int, bord *data.Border, obj *object.Object, level bool) {
	if level && data.CurrLevel != nil {
		if x == 0 && y > 0 && y <= bord.Height {
			if t := data.CurrLevel.Get(x, y-1); !t.IsSolidLevelTrans(data.CurrLevel.DoorsOpen) {
				if trans, ok := t.Transitions[data.Left]; ok && (trans.Open || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		} else if x == bord.Width+1 && y > 0 && y <= bord.Height {
			if t := data.CurrLevel.Get(x-2, y-1); !t.IsSolidLevelTrans(data.CurrLevel.DoorsOpen) {
				if trans, ok := t.Transitions[data.Right]; ok && (trans.Open || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		} else if y == 0 && x > 0 && x <= bord.Width {
			if t := data.CurrLevel.Get(x-1, y); !t.IsSolidLevelTrans(data.CurrLevel.DoorsOpen) {
				if trans, ok := t.Transitions[data.Down]; ok && (trans.Open || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
					return
				}
			}
		} else if y == bord.Height+1 && x > 0 && x <= bord.Width {
			if t := data.CurrLevel.Get(x-1, y-2); !t.IsSolidLevelTrans(data.CurrLevel.DoorsOpen) {
				if trans, ok := t.Transitions[data.Up]; ok && (trans.Open || (data.CurrLevel.Continuity && data.CurrLevel.DoorsOpen)) {
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

func DrawThinBorder(bord *data.Border, obj *object.Object) {
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

func DrawThinBorderReverse(bord *data.Border, obj *object.Object) {
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

func DrawThinBorderWhite(bord *data.Border, obj *object.Object) {
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

func DrawThinBorderBlue(bord *data.Border, obj *object.Object) {
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

func DrawThickBorder(bord *data.Border, obj *object.Object) {
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

func DrawThickBorderReverse(bord *data.Border, obj *object.Object) {
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

func DrawThickBorderWhite(bord *data.Border, obj *object.Object) {
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

func DrawThickBorderBlue(bord *data.Border, obj *object.Object) {
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

func DrawBlackSquare(x, y int, bord *data.Border, obj *object.Object) {
	mat := pixel.IM
	offset := pixel.V(world.TileSize*(float64(x)-float64(bord.Width+1)*0.5), world.TileSize*(float64(y)-float64(bord.Height+1)*0.5))
	sKey := "black_square_16"
	img.Batchers[constants.UIBatch].DrawSpriteColor(sKey, mat.Moved(obj.PostPos).Moved(offset), colornames.White)
}

func InitMainBorder() {
	if data.PuzzleBorder == nil {
		data.PuzzleBorderObject = object.New()
		data.PuzzleBorderObject.Pos.X = world.TileSize * 0.5 * constants.PuzzleWidth
		data.PuzzleBorderObject.Pos.Y = world.TileSize * 0.5 * constants.PuzzleHeight
		data.PuzzleBorderObject.Layer = 1
		data.PuzzleBorder = &data.Border{
			Width:  constants.PuzzleWidth,
			Height: constants.PuzzleHeight,
			Empty:  false,
		}
		data.MainBorder = myecs.Manager.NewEntity()
		data.MainBorder.AddComponent(myecs.Object, data.PuzzleBorderObject).
			AddComponent(myecs.Border, data.PuzzleBorder)
	}
}

func SetMainBorder(w, h int) {
	if data.PuzzleBorder != nil {
		data.PuzzleBorderObject.Pos.X = world.TileSize * 0.5 * float64(w)
		data.PuzzleBorderObject.Pos.Y = world.TileSize * 0.5 * float64(h)
		data.PuzzleBorder.Width = w
		data.PuzzleBorder.Height = h
	}
}

func InitBorder(fp *data.PlayArea) {
	if fp.Border == nil {
		fp.BorderObject = object.New()
		fp.BorderObject.Pos.X = world.TileSize * 0.5 * constants.PuzzleWidth
		fp.BorderObject.Pos.Y = world.TileSize * 0.5 * constants.PuzzleHeight
		fp.BorderObject.Layer = 1
		fp.Border = &data.Border{
			Width:  constants.PuzzleWidth,
			Height: constants.PuzzleHeight,
			Empty:  false,
		}
		fp.BorderEntity = myecs.Manager.NewEntity()
		fp.BorderEntity.AddComponent(myecs.Object, fp.BorderObject).
			AddComponent(myecs.Border, fp.Border)
	}
}

func SetBorder(fp *data.PlayArea) {
	if fp.Border != nil {
		fp.BorderObject.Pos.X = world.TileSize * 0.5 * float64(fp.Puzzle.Metadata.Width)
		fp.BorderObject.Pos.Y = world.TileSize * 0.5 * float64(fp.Puzzle.Metadata.Height)
		fp.Border.Width = fp.Puzzle.Metadata.Width
		fp.Border.Height = fp.Puzzle.Metadata.Height
	}
}
