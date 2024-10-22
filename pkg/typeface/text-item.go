package typeface

import (
	"fmt"
	"gemrunner/pkg/object"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
)

var (
	imd *imdraw.IMDraw
)

type Text struct {
	Raw     string
	Text    *text.Text
	Color   pixel.RGBA
	Align   Alignment
	Symbols []symbolHandle
	Debug   bool

	MaxWidth  float64
	MaxHeight float64

	Scalar     float64
	SymbolSize float64
	Obj        *object.Object

	rawLines      []string
	rawWithBreaks string
	lineCutInText []bool
	lineWidths    []float64
	dotPosArray   []pixel.Vec
	fullHeight    float64
}

func New(atlas string) *Text {
	tex := text.New(pixel.ZV, Atlases[atlas])
	tex.AlignedTo(pixel.BottomRight)
	obj := object.New()
	return &Text{
		Text:       tex,
		Color:      pixel.ToRGBA(colornames.White),
		SymbolSize: 1.,
		Scalar:     1.,
		Obj:        obj,
	}
}

func (item *Text) WithAnchor(anchor pixel.Anchor) *Text {
	item.Text.AlignedTo(anchor)
	return item
}

func (item *Text) WithAlign(alignment Alignment) *Text {
	item.Align = alignment
	return item
}

func (item *Text) WithLineHeight(lineHeight float64) *Text {
	item.Text.LineHeight = item.Text.Atlas().LineHeight() * lineHeight
	return item
}

func (item *Text) WithScalar(scalar float64) *Text {
	item.Obj.Sca = pixel.V(scalar, scalar)
	item.Scalar = scalar
	return item
}

func (item *Text) WithWidth(width float64) *Text {
	item.MaxWidth = width
	return item
}

func (item *Text) WithHeight(height float64) *Text {
	item.MaxHeight = height
	return item
}

func (item *Text) Update() {
	item.Obj.Sca = pixel.V(item.Scalar, item.Scalar)
}

func NewOld(atlas string, align Alignment, lineHeight, scalar, maxWidth, maxHeight float64) *Text {
	tex := text.New(pixel.ZV, Atlases[atlas])
	tex.LineHeight *= lineHeight
	obj := object.New()
	obj.Sca = pixel.V(scalar, scalar)
	return &Text{
		Text:       tex,
		Align:      align,
		Color:      pixel.ToRGBA(colornames.White),
		MaxWidth:   maxWidth,
		MaxHeight:  maxHeight,
		Scalar:     scalar,
		SymbolSize: 1.,
		Obj:        obj,
	}
}

func (item *Text) Draw(target pixel.Target) {
	if !item.Obj.Hidden {
		item.Text.Draw(target, item.Obj.Mat)
		if item.Debug {
			if imd == nil {
				imd = imdraw.New(nil)
			}
			imd.Clear()
			for _, d := range item.dotPosArray {
				imd.Color = colornames.Cadetblue
				imd.Push(d, d)
				imd.Line(2)
			}
			imd.EndShape = imdraw.RoundEndShape
			imd.Color = colornames.Indianred
			imd.Push(item.Text.Orig.Add(item.Obj.Pos), item.Text.Orig.Add(item.Obj.Pos))
			imd.Line(2)
			imd.Color = colornames.Lawngreen
			imd.Push(item.Text.Dot.Scaled(item.Scalar).Add(item.Obj.Pos), item.Text.Dot.Scaled(item.Scalar).Add(item.Obj.Pos))
			imd.Line(2)
			imd.Draw(target)
		}
	}
}

func (item *Text) SetWidth(width float64) {
	item.MaxWidth = width
	item.SetText(item.Raw)
}

func (item *Text) SetHeight(height float64) {
	item.MaxHeight = height
	item.SetText(item.Raw)
}

func (item *Text) SetColor(col pixel.RGBA) {
	item.Color = col
	item.updateText()
}

func (item *Text) SetSize(size float64) {
	item.Scalar = size
	item.SetText(item.Raw)
}

func (item *Text) SetPos(pos pixel.Vec) {
	item.Obj.Pos = pos
	item.updateText()
}

func (item *Text) SetOffset(pos pixel.Vec) {
	item.Obj.Offset = pos
	//item.updateText()
}

func (item *Text) UpdateText() {
	item.updateText()
}

func (item *Text) PrintLines() {
	for _, line := range item.rawLines {
		fmt.Println(line)
	}
}

func (item *Text) GetHeight() float64 {
	return item.fullHeight * item.Scalar
}

func (item *Text) Hide() {
	item.Obj.Hidden = true
}

func (item *Text) Show() {
	item.Obj.Hidden = false
}
