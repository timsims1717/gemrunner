package util

import (
	"github.com/gopxl/pixel"
	"math"
	"math/rand"
	"reflect"
	"strings"
)

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns the smaller number between a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger number between a and b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Contains(i int, a []int) bool {
	for _, as := range a {
		if as == i {
			return true
		}
	}
	return false
}

func ContainsStr(s string, a []string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

func RemoveStrUO(s string, a []string) []string {
	in := -1
	for i, c := range a {
		if c == s {
			in = i
			break
		}
	}
	if in > -1 {
		a[in] = a[0]
		return a[1:]
	}
	return a
}

// PointInside returns true if the pixel.Vec is inside the pixel.Rect
// when unprojected by the pixel.Matrix
func PointInside(p pixel.Vec, r pixel.Rect, m pixel.Matrix) bool {
	return r.Moved(pixel.V(-(r.W() / 2.0), -(r.H() / 2.0))).Contains(m.Unproject(p))
}

// Normalize takes a pixel.Vec and returns a normalized vector, or
// one with a magnitude of 1.0
func Normalize(p pixel.Vec) pixel.Vec {
	s := p.X*p.X + p.Y*p.Y
	if s == 0 {
		p.Y = 1.
	}
	return p.Scaled(1 / math.Sqrt(p.X*p.X+p.Y*p.Y))
}

// RectToOrigin moves a rectangle's min to (0,0)
func RectToOrigin(r pixel.Rect) pixel.Rect {
	return r.Moved(pixel.V(-r.Min.X, -r.Min.Y))
}

// Magnitude takes a pixel.Vec and returns the magnitude of the vector
func Magnitude(p pixel.Vec) float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func FMod(a, b float64) float64 {
	var mod float64
	if a < 0 {
		mod = -a
	} else {
		mod = a
	}
	if b < 0 {
		b = -b
	}

	for mod >= b {
		mod -= b
	}

	if a < 0 {
		return -mod
	}
	return mod
}

func UBound(a, b float64) float64 {
	if a >= 0. {
		return math.Min(a, math.Abs(b))
	} else {
		return math.Max(a, -math.Abs(b))
	}
}

func LBound(a, b float64) float64 {
	if a >= 0. {
		return math.Max(a, math.Abs(b))
	} else {
		return math.Min(a, -math.Abs(b))
	}
}

// RandomSample returns k unique integers in the specified range [a,b)
func RandomSampleRange(k, a, b int, rando *rand.Rand) []int {
	var res []int
	for i := a; i < b; i++ {
		res = append(res, i)
	}
	for i := len(res) - 1; i > 0; i-- {
		j := rando.Intn(i)
		res[i], res[j] = res[j], res[i]
	}
	if k > len(res) {
		return res
	}
	return res[:k]
}

// RandomSample returns k unique integers from l
func RandomSample(k int, l []int, rando *rand.Rand) []int {
	res := l
	for i := len(res) - 1; i > 0; i-- {
		j := rando.Intn(i)
		res[i], res[j] = res[j], res[i]
	}
	return res[:k]
}

// ConstrainR returns vec1 moved into r2 (at vec2), restrained by r1.
func ConstrainR(vec1, vec2 pixel.Vec, r1, r2 pixel.Rect) pixel.Vec {
	newPos := vec1
	if vec2.X+r2.W()*0.5 < vec1.X+r1.W()*0.5 {
		newPos.X = vec2.X + r2.W()*0.5 - r1.W()*0.5
	} else if vec2.X-r2.W()*0.5 > vec1.X-r1.W()*0.5 {
		newPos.X = vec2.X - r2.W()*0.5 + r1.W()*0.5
	}
	if vec2.Y+r2.H()*0.5 < vec1.Y+r1.H()*0.5 {
		newPos.Y = vec2.Y + r2.H()*0.5 - r1.H()*0.5
	} else if vec2.Y-r2.H()*0.5 > vec1.Y-r1.H()*0.5 {
		newPos.Y = vec2.Y - r2.H()*0.5 + r1.H()*0.5
	}
	return newPos
}

func OnlyNumbers(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		if '0' <= b && b <= '9' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func OnlyAlphaNumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func JustChars(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ' ' <= b && b <= '~' {
			result.WriteByte(b)
		}
	}
	return result.String()
}
