package shaders

import (
	"github.com/faiface/pixel/pixelgl"
	"os"
)

func EasyBindUniforms(c *pixelgl.Canvas, unifs ...interface{}) {
	if len(unifs)%2 != 0 {
		panic("needs to be divisable by 2")
	}
	for i := 0; i < len(unifs); i += 2 {

		c.SetUniform(unifs[i+0].(string), unifs[i+1])
	}
}

// LoadFileToString loads the contents of a file into a string
func LoadFileToString(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
