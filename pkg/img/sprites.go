package img

import (
	"encoding/json"
	"fmt"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
	"github.com/pkg/errors"
	"image/color"
	"os"
	"path/filepath"
)

var (
	IM        = pixel.IM
	Flip      = pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
	Flop      = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.))
	FlipFlop  = pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.))
	Batchers  = map[string]*Batcher{}
	batchers  []*Batcher
	IMDrawers = map[string]*IMDrawer{}
	imdraws   []*IMDrawer
)

func FullClear() {
	for _, batcher := range batchers {
		batcher.Clear()
	}
	for _, imd := range imdraws {
		imd.Clear()
	}
}

func Clear() {
	for _, batcher := range batchers {
		if batcher.AutoClear {
			batcher.Clear()
		}
	}
}

func Draw(target pixel.Target) {
	for _, batcher := range batchers {
		if batcher.AutoDraw {
			batcher.Draw(target)
		}
	}
	for _, imd := range imdraws {
		if imd.AutoDraw {
			imd.Draw(target)
		}
		if imd.AutoClear {
			imd.Clear()
		}
	}
}

type Sprite struct {
	Key    string
	Mask   string
	Batch  string
	Offset pixel.Vec
	Color  pixel.RGBA
}

func NewSprite(key, batch string) *Sprite {
	return &Sprite{
		Key:   key,
		Batch: batch,
		Color: util.White,
	}
}

func (s *Sprite) WithOffset(offset pixel.Vec) *Sprite {
	s.Offset = offset
	return s
}

func (s *Sprite) WithColor(color color.Color) *Sprite {
	s.Color = pixel.ToRGBA(color)
	return s
}

func (s *Sprite) WithMask(m string) *Sprite {
	s.Mask = m
	return s
}

type SpriteSheet struct {
	Img       pixel.Picture
	Sprites   []pixel.Rect
	SpriteMap map[string]pixel.Rect
	AnimMap   map[string]AnimDef
}

type AnimDef struct {
	Loop    bool
	Hold    bool
	Sprites []pixel.Rect
	dur     float64
}

type spriteFile struct {
	ImgFile   string   `json:"img"`
	Sprites   []sprite `json:"sprites"`
	Width     float64  `json:"width"`
	Height    float64  `json:"height"`
	SingleRow bool     `json:"singleRow"`
}

type sprite struct {
	K string  `json:"key"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`

	Loop   bool    `json:"loop"`
	Hold   bool    `json:"hold"`
	Dur    float64 `json:"dur"`
	Anim   bool    `json:"anim"`
	Frames int     `json:"frames"`
}

func LoadSpriteSheet(path string) (*SpriteSheet, error) {
	errMsg := "load sprite sheet"
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var fileSheet spriteFile
	err = decoder.Decode(&fileSheet)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	img, err := LoadImage(fmt.Sprintf("%s/%s", filepath.Dir(path), fileSheet.ImgFile))
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	sheet := &SpriteSheet{
		Img:       img,
		Sprites:   make([]pixel.Rect, 0),
		SpriteMap: make(map[string]pixel.Rect, 0),
		AnimMap:   make(map[string]AnimDef, 0),
	}
	x := 0.0
	for _, r := range fileSheet.Sprites {
		var rect pixel.Rect
		w := fileSheet.Width
		h := fileSheet.Height
		if r.W > 0.0 {
			w = r.W
		}
		if fileSheet.SingleRow {
			rect = pixel.R(x, 0.0, x+w, h)
			x += w
		} else {
			if r.H > 0.0 {
				h = r.H
			}
			rect = pixel.R(r.X, r.Y, r.X+w, r.Y+h)
		}
		sheet.Sprites = append(sheet.Sprites, rect)
		if r.K != "" {
			if def, ok := sheet.AnimMap[r.K]; ok {
				def.Sprites = append(def.Sprites, rect)
				for i := 1; i < r.Frames; i++ {
					def.Sprites = append(def.Sprites, rect)
				}
				sheet.AnimMap[r.K] = def
			} else {
				if r.Dur != 0.0 || r.Anim {
					spr := []pixel.Rect{rect}
					for i := 1; i < r.Frames; i++ {
						spr = append(spr, rect)
					}
					sheet.AnimMap[r.K] = AnimDef{
						Loop:    r.Loop,
						Hold:    r.Hold,
						Sprites: spr,
						dur:     r.Dur,
					}
				}
				sheet.SpriteMap[r.K] = rect
			}
		}
	}
	return sheet, nil
}
