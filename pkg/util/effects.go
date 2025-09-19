package util

import (
	gween "gemrunner/pkg/gween64"
	"gemrunner/pkg/gween64/ease"
	"github.com/KEINOS/go-noise"
	"github.com/gopxl/pixel"
)

type NoiseShaker struct {
	NoiseI    float64
	Speed     float64
	Seed      int64
	Type      noise.Algo
	Generator noise.Generator
	Gween     *gween.Tween
}

func NewShaker(speed, strength, decay float64, seed int64) *NoiseShaker {
	ns := &NoiseShaker{
		Speed: speed,
		Seed:  seed,
		Type:  noise.OpenSimplex,
		Gween: gween.New(strength, 0, decay, ease.Linear),
	}
	ns.Generator, _ = noise.New(ns.Type, ns.Seed)
	return ns
}

func (ns *NoiseShaker) Reset(seed int64) {
	ns.Seed = seed
	ns.NoiseI = 0
	ns.Gween.Reset()
}

func (ns *NoiseShaker) Shake(delta float64) (pixel.Vec, bool) {
	strength, f := ns.Gween.Update(delta)
	v := ns.ShakeOffset(delta, strength)
	return v, f
}

func (ns *NoiseShaker) ShakeOffset(delta, strength float64) pixel.Vec {
	ns.NoiseI += delta * ns.Speed
	x := ns.Generator.Eval64(1, ns.NoiseI) * strength
	y := ns.Generator.Eval64(100, ns.NoiseI) * strength
	return pixel.V(x, y)
}
