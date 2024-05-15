package typeface

import (
	"bytes"
	"fmt"
	"gemrunner/pkg/object"
	"github.com/gopxl/pixel"
)

const (
	OpenMarker  = '{'
	CloseMarker = '}'
	DivMarker   = ':'
	OpenItem    = "{"
	CloseItem   = "}"
)

var (
	RoundDot = false
)

func (item *Text) SetText(raw string) {
	if item.Raw == raw {
		return
	}
	item.Raw = raw
	item.rawLines = []string{}
	item.lineWidths = []float64{}
	b := 0
	bb := 0
	e := 0
	cut := false
	wSpace := false
	inBrackets := false
	widthMod := 0.
	maxLineWidth := 0.
	mode := ""
	buf := bytes.NewBuffer(nil)
	for i, r := range item.Raw {
		if !inBrackets {
			switch r {
			case '\n':
				cut = true
				e = i
			case ' ', '\t':
				wSpace = true
				e = i
			case OpenMarker:
				inBrackets = true
				bb = i
				continue
			case CloseMarker:
				fmt.Printf("extra closing bracket in text at position %d\n", i)
			}
			lineWidth := item.Text.BoundsOf(item.Raw[b:i]).W() - widthMod
			lineWidthRelative := lineWidth * item.RelativeSize
			if item.MaxWidth > 0. && lineWidthRelative > item.MaxWidth && wSpace {
				cut = true
			}
			if cut {
				if b >= e || e < 0 {
					item.rawLines = append(item.rawLines, "")
					item.lineWidths = append(item.lineWidths, 0.)
				} else {
					item.rawLines = append(item.rawLines, raw[b:e])
					item.lineWidths = append(item.lineWidths, lineWidth)
					if maxLineWidth < lineWidthRelative {
						maxLineWidth = lineWidthRelative
					}
				}
				cut = false
				wSpace = false
				widthMod = 0.
				b = e + 1
			}
		} else {
			switch r {
			case '\n':
				fmt.Printf("new line in bracketed text at position %d\n", i)
			case ' ', '\t':
				continue
			case OpenMarker:
				fmt.Printf("extra opening bracket at position %d\n", i)
				continue
			case CloseMarker:
				switch mode {
				case "symbol":
					if sym, ok := theSymbols[buf.String()]; ok {
						widthMod -= sym.spr.Frame().W() * item.SymbolSize * sym.sca / item.RelativeSize
					}
				}
				widthMod += item.Text.BoundsOf(item.Raw[bb : i+1]).W()
				inBrackets = false
				mode = ""
				buf.Reset()
			case DivMarker:
				mode = buf.String()
				buf.Reset()
			default:
				buf.WriteRune(r)
			}
		}
	}
	item.rawLines = append(item.rawLines, raw[b:])
	lineWidth := item.Text.BoundsOf(item.Raw[b:]).W() - widthMod
	lineWidthRelative := lineWidth * item.RelativeSize
	item.lineWidths = append(item.lineWidths, lineWidth)
	item.fullHeight = float64(len(item.rawLines)) * item.Text.LineHeight
	if maxLineWidth < lineWidthRelative {
		maxLineWidth = lineWidthRelative
	}
	maxX := maxLineWidth
	maxY := item.MaxHeight
	if maxY == 0. {
		maxY = item.fullHeight * item.RelativeSize
	}
	item.Width = maxX
	item.Height = maxY
	item.updateText()
}

func (item *Text) updateText() {
	item.dotPosArray = []pixel.Vec{}
	item.Text.Clear()
	item.Text.Color = item.Color
	if item.Align.H == Center {
		item.Text.Orig.X = -item.Width * 0.25
	} else if item.Align.H == Right {
		item.Text.Orig.X = -item.Width * 0.5
	}
	if item.Align.V == Center {
		item.Text.Orig.Y = -item.fullHeight
	} else if item.Align.V == Top {
		item.Text.Orig.Y = -item.fullHeight
	}
	//var colorStack []color.RGBA
	item.Symbols = []symbolHandle{}
	inBrackets := false
	mode := ""
	buf := bytes.NewBuffer(nil)
	//item.Text.Dot.Y -= item.Text.LineHeight
	//if item.Align.V == Center {
	//	item.Text.Dot.Y += item.fullHeight * 0.5
	//} else if item.Align.V == Bottom {
	//	item.Text.Dot.Y += item.fullHeight
	//}
	for li, line := range item.rawLines {
		if li != 0 {
			item.Text.WriteRune('\n')
		}
		//b := 0
		inBrackets = false
		//if item.Align.H == Center {
		//	item.Text.Dot.X -= item.lineWidths[li] * 0.5
		//} else if item.Align.H == Right {
		//	item.Text.Dot.X -= item.lineWidths[li]
		//}
		for _, r := range line {
			if !inBrackets {
				switch r {
				case OpenMarker:
					inBrackets = true
				default:
					item.dotPosArray = append(item.dotPosArray, item.Text.Dot.Scaled(item.RelativeSize))
					item.Text.WriteRune(r)
				}
			} else {
				switch r {
				case CloseMarker:
					switch mode {
					case "symbol":
						if sym, ok := theSymbols[buf.String()]; ok {
							item.dotPosArray = append(item.dotPosArray, item.Text.Dot.Scaled(item.RelativeSize))
							obj := object.New()
							obj.Sca = pixel.V(item.SymbolSize, item.SymbolSize).Scaled(sym.sca)
							obj.Pos = item.Obj.Pos
							obj.Pos = obj.Pos.Add(item.Text.Dot.Scaled(item.RelativeSize))
							obj.Pos = obj.Pos.Add(pixel.V(sym.spr.Frame().W()*0.5, sym.spr.Frame().H()*0.25).Scaled(item.SymbolSize * sym.sca))
							item.Symbols = append(item.Symbols, symbolHandle{
								symbol: sym,
								trans:  obj,
							})
							item.Text.Dot.X += sym.spr.Frame().W() * item.SymbolSize * sym.sca / item.RelativeSize
						}
					}
					//b = i + 1
					inBrackets = false
					mode = ""
					buf.Reset()
				case DivMarker:
					mode = buf.String()
					buf.Reset()
				default:
					buf.WriteRune(r)
				}
			}
		}
		item.dotPosArray = append(item.dotPosArray, item.Text.Dot.Scaled(item.RelativeSize))
	}
	bounds := item.Text.Bounds()
	item.Obj.SetRect(pixel.R(0, 0, bounds.W()*item.RelativeSize, bounds.H()*item.RelativeSize))
	if item.Align.H == Left {
		item.Obj.Rect = item.Obj.Rect.Moved(pixel.V(item.Obj.Rect.W()*0.5, 0))
	} else if item.Align.H == Right {
		//item.Obj.Rect = item.Obj.Rect.Moved(pixel.V(item.Obj.Rect.W(), 0))
	}
	//if item.Align.V == Bottom {
	//	item.Obj.Rect = item.Obj.Rect.Moved(pixel.V(item.Obj.Rect.H()*-0.5, 0))
	//} else if item.Align.V == Top {
	//	item.Obj.Rect = item.Obj.Rect.Moved(pixel.V(item.Obj.Rect.H(), 0))
	//}
}
