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
	Hidden  bool
	Debug   bool

	Increment bool
	CurrPos   int
	Width     float64
	Height    float64
	MaxWidth  float64
	MaxHeight float64
	MaxLines  int

	RelativeSize float64
	SymbolSize   float64
	Obj          *object.Object

	rawLines    []string
	lineWidths  []float64
	dotPosArray []pixel.Vec
	fullHeight  float64
}

func New(atlas string, align Alignment, lineHeight, relativeSize, maxWidth, maxHeight float64) *Text {
	tex := text.New(pixel.ZV, Atlases[atlas])
	tex.LineHeight *= lineHeight
	obj := object.New()
	obj.Sca = pixel.V(relativeSize, relativeSize)
	return &Text{
		Text:         tex,
		Align:        align,
		Color:        pixel.ToRGBA(colornames.White),
		Width:        maxWidth,
		Height:       maxHeight,
		MaxWidth:     maxWidth,
		MaxHeight:    maxHeight,
		MaxLines:     int(maxHeight / (tex.LineHeight * relativeSize)),
		RelativeSize: relativeSize,
		SymbolSize:   1.,
		Obj:          obj,
	}
}

func (item *Text) Draw(target pixel.Target) {
	if !item.Hidden {
		item.Text.Draw(target, item.Obj.Mat)
		if item.Debug {
			if imd == nil {
				imd = imdraw.New(nil)
			}
			imd.Clear()
			for _, d := range item.dotPosArray {
				imd.Color = colornames.Cadetblue
				imd.Push(d.Add(item.Obj.Pos), d.Add(item.Obj.Pos))
				imd.Line(2)
			}
			imd.EndShape = imdraw.RoundEndShape
			imd.Color = colornames.Indianred
			imd.Push(item.Text.Orig.Add(item.Obj.Pos), item.Text.Orig.Add(item.Obj.Pos))
			imd.Line(2)
			imd.Color = colornames.Lawngreen
			imd.Push(item.Text.Dot.Scaled(item.RelativeSize).Add(item.Obj.Pos), item.Text.Dot.Scaled(item.RelativeSize).Add(item.Obj.Pos))
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
	item.RelativeSize = size
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

func (item *Text) IncrementTextPos() {
	if item.Increment {

	}
}

func (item *Text) SkipIncrement() {
	if item.Increment {

	}
}

func (item *Text) PrintLines() {
	for _, line := range item.rawLines {
		fmt.Println(line)
	}
}
