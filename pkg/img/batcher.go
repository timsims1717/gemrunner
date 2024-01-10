package img

import (
	"fmt"
	"github.com/gopxl/pixel"
	"image/color"
)

type Batcher struct {
	Key        string
	Index      int
	Sprites    map[string]*pixel.Sprite
	Animations map[string]*Animation
	batch      *pixel.Batch
	AutoDraw   bool
	AutoClear  bool
}

func AddBatcher(key string, sheet *SpriteSheet, autoDraw, autoClear bool) {
	if _, ok := Batchers[key]; ok {
		Batchers[key].SetSpriteSheet(sheet)
		Batchers[key].AutoDraw = autoDraw
		Batchers[key].AutoClear = autoClear
	} else {
		Batchers[key] = NewBatcher(key, sheet, autoDraw, autoClear)
		batchers = append(batchers, Batchers[key])
	}
}

func NewBatcher(key string, sheet *SpriteSheet, autoDraw, autoClear bool) *Batcher {
	b := &Batcher{
		Key:       key,
		Index:     len(batchers),
		AutoDraw:  autoDraw,
		AutoClear: autoClear,
	}
	b.SetSpriteSheet(sheet)
	return b
}

func (b *Batcher) GetFrame(key string, index int) *pixel.Sprite {
	if a, ok := b.Animations[key]; ok {
		if len(a.S) > index {
			return a.S[index]
		}
	}
	return nil
}

func (b *Batcher) GetSprite(key string) *pixel.Sprite {
	if s, ok := b.Sprites[key]; ok {
		return s
	}
	return nil
}

func (b *Batcher) GetAnimation(key string) *Animation {
	if a, ok := b.Animations[key]; ok {
		return a
	}
	return nil
}

func (b *Batcher) SetSpriteSheet(sheet *SpriteSheet) {
	b.batch = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Img)
	b.Sprites = make(map[string]*pixel.Sprite)
	b.Animations = make(map[string]*Animation)
	for k, r := range sheet.SpriteMap {
		b.Sprites[k] = pixel.NewSprite(sheet.Img, r)
	}
	for k, a := range sheet.AnimMap {
		b.Animations[k] = NewAnimation(sheet, a.Sprites, a.Loop, a.Hold, a.dur)
	}
}

func (b *Batcher) Clear() {
	b.batch.Clear()
}

func (b *Batcher) Batch() *pixel.Batch {
	return b.batch
}

func (b *Batcher) DrawSprite(key string, mat pixel.Matrix) {
	if spr, ok := b.Sprites[key]; ok {
		spr.Draw(b.batch, mat)
	} else {
		fmt.Printf("couldn't draw sprite '%s' with batch %s\n", key, b.Key)
	}
}

func (b *Batcher) DrawSpriteColor(key string, mat pixel.Matrix, mask color.Color) {
	if spr, ok := b.Sprites[key]; ok {
		spr.DrawColorMask(b.batch, mat, mask)
	} else {
		fmt.Printf("couldn't draw sprite '%s' with batch %s\n", key, b.Key)
	}
}

func (b *Batcher) Draw(target pixel.Target) {
	b.batch.Draw(target)
}
