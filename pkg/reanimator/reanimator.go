package reanimator

import (
	"gemrunner/pkg/img"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
	"image/color"
	"time"
)

var (
	Timer       time.Time
	FRate       = 10
	inter       float64
	FrameTime   float64
	FrameSwitch bool
)

type TreeSet struct {
	Set []*Tree
}

func NewSet() *TreeSet {
	return &TreeSet{}
}

func (ts *TreeSet) Add(tree *Tree) *TreeSet {
	ts.Set = append(ts.Set, tree)
	return ts
}

func (ts *TreeSet) Update() {
	var base *Tree
	for i, anim := range ts.Set {
		if anim == nil {
			continue
		}
		if i == 0 {
			base = anim
		}
		anim.Update()
	}
	for _, anim := range ts.Set {
		if anim == nil || base == nil || base.anim == nil {
			continue
		}
		if anim.Dependent {
			anim.SetAnim(base.anim.Key, base.frame)
		}
	}
}

type Tree struct {
	Root      *Switch
	spr       *pixel.Sprite
	anim      *Anim
	animKey   string
	frame     int
	update    bool
	Done      bool
	Default   string
	Dependent bool
}

func SetFrameRate(fRate int) {
	FRate = fRate
	inter = 1. / float64(fRate)
	FrameTime = inter
}

func Reset() {
	Timer = time.Now()
}

func Update() {
	FrameSwitch = time.Since(Timer).Seconds() > inter
	if FrameSwitch {
		Reset()
	}
}

func NewSimple(anim *Anim) *Tree {
	t := &Tree{
		Root: NewSwitch().
			AddAnimation(anim).
			SetChooseFn(func() string {
				return anim.Key
			}),
	}
	t.Update()
	return t
}

func New(root *Switch, def string) *Tree {
	t := &Tree{
		Root:    root,
		update:  true,
		Default: def,
	}
	t.Update()
	return t
}

func (t *Tree) ForceUpdate() {
	t.update = true
}

func (t *Tree) GetCurrentAnim() *Anim {
	return t.anim
}

func (t *Tree) GetCurrentFrame() int {
	return t.frame
}

func (t *Tree) Update() {
	if !t.Done {
		if FrameSwitch || t.update {
			t.update = false
			t.anim = t.Root.choose()
			if t.anim == nil {
				t.spr = nil
				t.animKey = ""
				t.frame = 0
			} else {
				pKey := t.animKey
				pFrame := t.frame
				var trigger int
				if t.anim.Key != t.animKey {
					t.anim.Step = 0
					trigger = 0
				} else if !t.anim.Freeze {
					t.anim.Step++
					trigger = t.anim.Step
					if t.anim.Step%len(t.anim.S) == 0 {
						switch t.anim.Finish {
						case Loop:
							t.anim.Step = 0
							trigger = 0
						case Hold:
							t.anim.Step = len(t.anim.S) - 1
						case Tran:
							t.anim.Step = len(t.anim.S) - 1
							t.update = true
						case Done:
							t.anim.Step = len(t.anim.S) - 1
							t.Done = true
						}
					}
				}
				if t.anim.Triggers != nil {
					if fn, ok := t.anim.Triggers[trigger]; ok {
						fn(t.anim, pKey, pFrame)
					}
				}
				t.spr = t.anim.S[t.anim.Step]
				if t.update {
					t.animKey = t.Default
					t.frame = t.anim.Step
				} else {
					t.animKey = t.anim.Key
					t.frame = t.anim.Step
				}
			}
		}
	}
}

func (t *Tree) SetAnim(key string, frame int) {
	if a, ok := t.Root.Elements[key]; ok {
		if a.Anim != nil {
			t.anim = a.Anim
			t.animKey = key
			t.frame = frame
			t.anim.Step = frame
		}
	}
}

type Result struct {
	Spr   *pixel.Sprite
	Off   pixel.Vec
	Col   pixel.RGBA
	Batch string
}

func (t *Tree) CurrentSprite() *Result {
	if t.spr == nil {
		return nil
	}
	return &Result{
		Spr:   t.spr,
		Off:   t.anim.Offset,
		Col:   t.anim.Color,
		Batch: t.anim.Batch,
	}
}

func (t *Tree) Draw(target pixel.Target, mat pixel.Matrix) {
	if t.spr != nil && !t.Done {
		t.spr.Draw(target, mat)
	}
}

func (t *Tree) DrawColorMask(target pixel.Target, mat pixel.Matrix, col color.RGBA) {
	if t.spr != nil && !t.Done {
		t.spr.DrawColorMask(target, mat, col)
	}
}

type switchEl struct {
	Switch *Switch
	Anim   *Anim
}

type Switch struct {
	Elements map[string]*switchEl
	Choose   func() string
}

func NewSwitch() *Switch {
	return &Switch{
		Elements: map[string]*switchEl{},
		Choose:   func() string { return "" },
	}
}

func (s *Switch) AddNull(key string) *Switch {
	s.Elements[key] = &switchEl{}
	return s
}

func (s *Switch) AddAnimation(anim *Anim) *Switch {
	s.Elements[anim.Key] = &switchEl{
		Anim: anim,
	}
	return s
}

func (s *Switch) AddSubSwitch(ss *Switch, key string) *Switch {
	s.Elements[key] = &switchEl{
		Switch: ss,
	}
	return s
}

func (s *Switch) SetChooseFn(fn func() string) *Switch {
	s.Choose = fn
	return s
}

func (s *Switch) choose() *Anim {
	el := s.Elements[s.Choose()]
	if el.Switch != nil {
		return el.Switch.choose()
	} else if el.Anim != nil {
		return el.Anim
	} else {
		return nil
	}
}

type Anim struct {
	Key      string
	S        []*pixel.Sprite
	Step     int
	Finish   Finish
	Freeze   bool
	Triggers map[int]func(*Anim, string, int)

	Offset pixel.Vec
	Color  pixel.RGBA
	Batch  string
}

type Finish int

const (
	Loop = iota
	Hold
	Tran
	Done
)

func (anim *Anim) WithColor(col pixel.RGBA) *Anim {
	anim.Color = col
	return anim
}

func (anim *Anim) WithBatch(batch string) *Anim {
	anim.Batch = batch
	return anim
}

func (anim *Anim) WithOffset(offset pixel.Vec) *Anim {
	anim.Offset = offset
	return anim
}

func NewAnimFromSprite(key string, spr *pixel.Sprite, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      []*pixel.Sprite{spr},
		Finish: f,
		Color:  util.White,
	}
}

func NewAnimFromSprites(key string, spr []*pixel.Sprite, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      spr,
		Finish: f,
		Color:  util.White,
	}
}

func NewBatchSprite(key string, batch *img.Batcher, spr string, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      []*pixel.Sprite{batch.GetSprite(spr)},
		Finish: f,
		Batch:  batch.Key,
		Color:  util.White,
	}
}

func NewBatchAnimation(key string, batch *img.Batcher, anim string, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      batch.GetAnimation(anim).S,
		Finish: f,
		Batch:  batch.Key,
		Color:  util.White,
	}
}

func NewBatchAnimationFrame(key string, batch *img.Batcher, anim string, frame int, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      []*pixel.Sprite{batch.GetAnimation(anim).S[frame]},
		Finish: f,
		Batch:  batch.Key,
		Color:  util.White,
	}
}

func NewBatchAnimationCustom(key string, batch *img.Batcher, anim string, frames []int, f Finish) *Anim {
	spr := batch.GetAnimation(anim).S
	var nSpr []*pixel.Sprite
	for i := 0; i < len(frames); i++ {
		nSpr = append(nSpr, spr[frames[i]])
	}
	return &Anim{
		Key:    key,
		S:      nSpr,
		Finish: f,
		Batch:  batch.Key,
		Color:  util.White,
	}
}

func NewAnimFromSheet(key string, spriteSheet *img.SpriteSheet, rs []int, f Finish) *Anim {
	var spr []*pixel.Sprite
	if len(rs) > 0 {
		for _, r := range rs {
			spr = append(spr, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[r]))
		}
	} else {
		for _, s := range spriteSheet.Sprites {
			spr = append(spr, pixel.NewSprite(spriteSheet.Img, s))
		}
	}
	return &Anim{
		Key:    key,
		S:      spr,
		Step:   0,
		Finish: f,
		Color:  util.White,
	}
}

func (anim *Anim) SetEndTrigger(fn func()) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	anim.Triggers[len(anim.S)] = func(*Anim, string, int) {
		fn()
	}
	return anim
}

func (anim *Anim) SetTriggerC(i int, fn func(*Anim, string, int)) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	anim.Triggers[i] = fn
	return anim
}

func (anim *Anim) SetTrigger(i int, fn func()) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	anim.Triggers[i] = func(*Anim, string, int) {
		fn()
	}
	return anim
}

func (anim *Anim) SetTriggerCAll(fn func(*Anim, string, int)) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	for i := range anim.S {
		anim.SetTriggerC(i, fn)
	}
	//anim.SetTriggerC(len(anim.S), fn)
	return anim
}

func (anim *Anim) SetTriggerAll(fn func()) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	for i := range anim.S {
		anim.SetTrigger(i, fn)
	}
	//anim.SetTrigger(len(anim.S), fn)
	return anim
}

func (anim *Anim) Reverse() *Anim {
	var r []*pixel.Sprite
	for i := len(anim.S) - 1; i >= 0; i-- {
		r = append(r, anim.S[i])
	}
	anim.S = r
	return anim
}

func (anim *Anim) Copy() *Anim {
	return &Anim{
		Key:      anim.Key,
		S:        anim.S,
		Step:     anim.Step,
		Finish:   anim.Finish,
		Triggers: anim.Triggers,
		Color:    anim.Color,
		Batch:    anim.Batch,
		Offset:   anim.Offset,
	}
}
