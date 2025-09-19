package debug

import (
	"fmt"
	"gemrunner/pkg/typeface"
	"github.com/gopxl/pixel"
	"strings"
)

var (
	debugText *typeface.Text
	lines     []string
)

func InitializeText() {
	debugText = typeface.New("basic").
		WithAnchor(pixel.BottomRight).
		WithScalar(2.)
}

func DrawText(target pixel.Target, dim pixel.Vec) {
	if ShowText {
		var sb strings.Builder
		for i, line := range lines {
			if i != 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(line)
		}
		debugText.SetText(sb.String())
		debugText.Obj.Pos = winV.Add(pixel.V(dim.X*-0.5+2., dim.Y*0.5-2))
		debugText.Obj.Update()
		debugText.Draw(target)
	}
}

func AddText(s string) {
	lines = append(lines, s)
}

func AddIntCoords(label string, x, y int) {
	lines = append(lines, fmt.Sprintf("%s: (%d,%d)", label, x, y))
}

func AddTruthText(label string, b bool) {
	lines = append(lines, fmt.Sprintf("%s: %t", label, b))
}

func InsertText(s string, i int) {
	if i < 0 || len(lines) <= i || len(lines) == 0 {
		AddText(s)
	} else {
		tmp := append(lines[:i], s)
		tmp = append(tmp, lines[i:]...)
		lines = tmp
	}
}

func SetTextColor(col pixel.RGBA) {
	fpsText.SetColor(col)
	versionText.SetColor(col)
	debugText.SetColor(col)
}
