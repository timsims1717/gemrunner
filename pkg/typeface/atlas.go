package typeface

import (
	"github.com/golang/freetype/truetype"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"os"
)

var (
	BasicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	Atlases    = map[string]*text.Atlas{
		"basic": BasicAtlas,
	}
)

func LoadTTF(path string, size float64) (font.Face, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func LoadBytes(bytes []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
