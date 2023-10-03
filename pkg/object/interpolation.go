package object

import (
	gween "gemrunner/pkg/gween64"
	"gemrunner/pkg/gween64/ease"
)

type InterpolationTarget int

const (
	InterpolateX = iota
	InterpolateY
	InterpolateOffX
	InterpolateOffY
	InterpolateRot
	InterpolateSX
	InterpolateSY
	InterpolateR
	InterpolateG
	InterpolateB
	InterpolateA
	InterpolateCol
	InterpolateCustom
)

type Interpolation struct {
	Sequence *gween.Sequence
	Target   InterpolationTarget
	Value    *float64

	OnComplete func()
}

func NewInterpolation(target InterpolationTarget) *Interpolation {
	return &Interpolation{
		Target: target,
	}
}

func (i *Interpolation) SetValue(v *float64) *Interpolation {
	i.Value = v
	return i
}

func (i *Interpolation) SetGween(begin, end, duration float64, easing ease.TweenFunc) *Interpolation {
	i.Sequence = gween.NewSequence(gween.New(begin, end, duration, easing))
	return i
}

func (i *Interpolation) AddGween(begin, end, duration float64, easing ease.TweenFunc) *Interpolation {
	if i.Sequence == nil {
		i.Sequence = gween.NewSequence(gween.New(begin, end, duration, easing))
	} else {
		i.Sequence.Add(gween.New(begin, end, duration, easing))
	}
	return i
}

func (i *Interpolation) SetOnComplete(fn func()) *Interpolation {
	i.OnComplete = fn
	return i
}
