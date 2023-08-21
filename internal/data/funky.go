package data

import pxginput "github.com/timsims1717/pixel-go-input"

type HoverFunky struct {
	Fn func(in *pxginput.Input)
}

func NewHoverFn(fn func(in *pxginput.Input)) *HoverFunky {
	return &HoverFunky{Fn: fn}
}
