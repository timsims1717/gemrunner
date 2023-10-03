// Package gween provides the Tween struct that allows an easing function to be
// animated over time. This can be used in tandem with the ease package to provide
// the easing functions.
package gween

import "gemrunner/pkg/gween64/ease"

type (
	// Tween encapsulates the easing function along with timing data. This allows
	// a ease.TweenFunc to be used to be easily animated.
	Tween struct {
		duration float64
		time     float64
		begin    float64
		end      float64
		change   float64
		Overflow float64
		easing   ease.TweenFunc
		reverse  bool
		End      float64
	}
)

// New will return a new Tween when passed a beginning and end value, the duration
// of the tween and the easing function to animate between the two values. The
// easing function can be one of the provided easing functions from the ease package
// or you can provide one of your own.
func New(begin, end, duration float64, easing ease.TweenFunc) *Tween {
	return &Tween{
		begin:    begin,
		end:      end,
		change:   end - begin,
		duration: duration,
		easing:   easing,
		Overflow: 0,
		reverse:  false,
		End:      end,
	}
}

// Set will set the current time along the duration of the tween. It will then return
// the current value as well as a boolean to determine if the tween is finished.
func (tween *Tween) Set(time float64) (current float64, isFinished bool) {
	switch {
	case time <= 0:
		tween.Overflow = time
		tween.time = 0
		current = tween.begin
	case time >= tween.duration:
		tween.Overflow = time - tween.duration
		tween.time = tween.duration
		current = tween.end
	default:
		tween.Overflow = 0
		tween.time = time
		current = tween.easing(tween.time, tween.begin, tween.change, tween.duration)
	}

	if tween.reverse {
		return current, tween.time <= 0
	}
	return current, tween.time >= tween.duration
}

// Reset will set the Tween to the beginning of the two values.
func (tween *Tween) Reset() {
	if tween.reverse {
		tween.Set(tween.duration)
	} else {
		tween.Set(0)
	}
}

// Update will increment the timer of the Tween and ease the value. It will then
// return the current value as well as a bool to mark if the tween is finished or not.
func (tween *Tween) Update(dt float64) (current float64, isFinished bool) {
	if tween.reverse {
		return tween.Set(tween.time - dt)
	}
	return tween.Set(tween.time + dt)
}
