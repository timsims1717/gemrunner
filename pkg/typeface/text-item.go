package typeface

import (
	"fmt"
	"gemrunner/pkg/object"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Text struct {
	Raw     string
	Text    *text.Text
	Color   pixel.RGBA
	Align   Alignment
	Symbols []symbolHandle
	NoShow  bool

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

	rawLines   []string
	lineWidths []float64
	fullHeight float64
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
	if !item.NoShow {
		item.Text.DrawColorMask(target, item.Obj.Mat, item.Obj.Mask)
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
