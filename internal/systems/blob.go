package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/gween64/ease"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"math"
)

type Blob struct {
	Object     *object.Object
	FaceSprite *img.Sprite
	State      data.BossState
	UHeadPos   mgl32.Vec2
	HeadPos    pixel.Vec
	Target     pixel.Vec
	StartCount int
	SlugCount  int
	SpawnCount int
	LastSlug   *data.Tile
	UScale     float32
	UTime      float32
	InnerColor pixel.RGBA
	OuterColor pixel.RGBA
	Entity     *ecs.Entity
}

func CreateBlobBoss(count int, inner, outer pixel.RGBA) *Blob {
	obj := object.New().WithFixedID(constants.BossBlob)
	obj.Layer = 8
	spr := img.NewSprite("blob_face", constants.TileBatch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)
	b := &Blob{
		Object:     obj,
		FaceSprite: spr,
		State:      data.BossStart,
		StartCount: count,
		SlugCount:  count,
		UScale:     1.,
		InnerColor: inner,
		OuterColor: outer,
		Entity:     e,
	}
	data.CurrentPlayArea.BackgroundView.Canvas.SetUniform("uTime", &b.UTime)
	data.CurrentPlayArea.BackgroundView.Canvas.SetUniform("uBallCount", int32(5))
	data.CurrentPlayArea.BackgroundView.Canvas.SetUniform("uHeadPos", &b.UHeadPos)
	data.CurrentPlayArea.BackgroundView.Canvas.SetUniform("uColorInner", util.RGBAToVec3(inner))
	data.CurrentPlayArea.BackgroundView.Canvas.SetUniform("uColorOuter", util.RGBAToVec3(outer))
	data.CurrentPlayArea.BackgroundView.Canvas.SetFragmentShader(data.BossBlobShader)
	data.BossCounter0 = 0
	return b
}

func (b *Blob) GetID() string {
	return b.Object.ID
}

func (b *Blob) SetState(state data.BossState) {
	b.Object.Hidden = b.State == data.BossStart
	b.State = state
}

func (b *Blob) GetState() string {
	return b.State.String()
}

func (b *Blob) Update() {
	switch b.State {
	case data.BossStart:
		b.Object.Hidden = true
		b.UTime += float32(timing.DT) * float32(constants.Configuration.Gameplay.FrameRate) / 60.
		b.HeadPos.X = 0.5 + math.Sin(float64(b.UTime))*0.125
		b.HeadPos.Y = math.Cos(float64(b.UTime)*7.) * 0.015
		b.UHeadPos[0] = float32(b.HeadPos.X)
		b.UHeadPos[1] = float32(b.HeadPos.Y)
		b.Object.SetPos(pixel.V(b.HeadPos.X*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().W(), b.HeadPos.Y*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().H()))
		if data.CurrLevel.FrameCycle > 3 {
			b.State = data.BossIntro
			b.Target = BlobHeadPos(b.UTime, pixel.V(0.5, 0.65))
			data.CurrLevel.Paused = true
			b.Object.Hidden = false
			dur := constants.BlobIntroTime * float64(constants.Configuration.Gameplay.FrameRate) / 60.
			interX := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&b.HeadPos.X).
				SetGween(b.HeadPos.X, b.Target.X, dur, ease.InOutCubic)
			interY := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&b.HeadPos.Y).
				SetGween(b.HeadPos.Y, b.Target.Y, dur, ease.InOutBack)
			b.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interX, interY})
			c := random.Effects.Intn(5) + 6
			for i := 0; i < c; i++ {
				pos := pixel.V(b.Object.Pos.X, b.Object.Pos.Y)
				pos.X += random.Effects.Float64()*2. - 0.5
				pos.Y += random.Effects.Float64()
				t := data.CurrLevel.Get(random.Effects.Intn(data.CurrLevel.Metadata.Width), random.Effects.Intn(data.CurrLevel.Metadata.Height/4))
				target := t.Object.Pos
				target.X += (random.Effects.Float64() - 0.5) * world.TileSize
				target.Y += (random.Effects.Float64() - 0.5) * world.TileSize
				SlugSpawner(pos, target, data.TileMetadata{}, false)
			}
		}
	case data.BossIntro:
		b.UHeadPos[0] = float32(b.HeadPos.X)
		b.UHeadPos[1] = float32(b.HeadPos.Y)
		b.Object.SetPos(pixel.V(b.HeadPos.X*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().W(), b.HeadPos.Y*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().H()))
		if !b.Entity.HasComponent(myecs.Interpolation) {
			b.State = data.BossWaiting
			data.CurrLevel.Paused = false
		}
	case data.BossWaiting:
		data.CurrLevel.Paused = false
		b.UTime += float32(timing.DT) * float32(constants.Configuration.Gameplay.FrameRate) / 60.
		b.HeadPos = BlobHeadPos(b.UTime, pixel.V(0.5, 0.65))
		b.UHeadPos[0] = float32(b.HeadPos.X)
		b.UHeadPos[1] = float32(b.HeadPos.Y)
		b.Object.SetPos(pixel.V(b.HeadPos.X*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().W(), b.HeadPos.Y*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().H()))
		if data.CurrLevel.FrameCycle%2 == 0 && data.CurrLevel.FrameChange {
			if (data.BossCounter0 == 0 && b.SlugCount == 0) || data.CurrLevel.DoorsOpen {
				b.State = data.BossDying
				data.CurrLevel.Paused = true
				data.BossCounter1 = 0
			} else if data.BossCounter0 < 5 && data.BossCounter0 < b.SlugCount {
				b.State = data.BossAction
			}
		}
	case data.BossAction:
		if data.CurrLevel.FrameChange {
			// create a slug until up to 5
			// choose a target
			t := GetSlugRegenTile()
			if t == nil { // kill the boss if it can't spawn a slug
				b.State = data.BossDying
				data.CurrLevel.Paused = true
				data.BossCounter1 = 0
				break
			}
			// create slug spawning blob
			SlugSpawner(b.Object.Pos, t.Object.Pos, t.Metadata, true)
			b.SlugCount--
			data.BossCounter0++
			if data.BossCounter0 >= b.SlugCount || data.BossCounter0 >= 5 {
				b.State = data.BossWaiting
			}
		}
	case data.BossDying:
		b.UHeadPos[0] = float32(b.HeadPos.X)
		b.UHeadPos[1] = float32(b.HeadPos.Y)
		b.Object.SetPos(pixel.V(b.HeadPos.X*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().W(), b.HeadPos.Y*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().H()))
		if reanimator.FrameSwitch {
			if data.BossCounter1 < 8 {
				data.BossCounter1++
			}
			if data.BossCounter1 == 8 {
				x := b.HeadPos.X + (random.Effects.Float64()-0.5)*0.3
				y := b.HeadPos.Y + (random.Effects.Float64()-0.5)*0.08
				dur := constants.BlobDyingTime * float64(constants.Configuration.Gameplay.FrameRate) / 60.
				interX := object.NewInterpolation(object.InterpolateCustom).
					SetValue(&b.HeadPos.X).
					SetGween(b.HeadPos.X, x, dur, ease.InElastic)
				interY := object.NewInterpolation(object.InterpolateCustom).
					SetValue(&b.HeadPos.Y).
					SetGween(b.HeadPos.Y, y, dur, ease.InElastic)
				b.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interX, interY})
				data.BossCounter1 = 9
			} else if data.BossCounter1 == 9 && !b.Entity.HasComponent(myecs.Interpolation) {
				data.CurrLevel.Paused = false
				b.Object.Hidden = true
				dur := constants.BlobIntroTime * float64(constants.Configuration.Gameplay.FrameRate) / 60.
				interY := object.NewInterpolation(object.InterpolateCustom).
					SetValue(&b.HeadPos.Y).
					SetGween(b.HeadPos.Y, 0., dur, ease.InOutBack)
				b.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interY})
				c := random.Effects.Intn(5) + 6
				for i := 0; i < c; i++ {
					pos := pixel.V(b.Object.Pos.X, b.Object.Pos.Y)
					pos.X += random.Effects.Float64()*2. - 0.5
					pos.Y += random.Effects.Float64()
					t := data.CurrLevel.Get(random.Effects.Intn(data.CurrLevel.Metadata.Width), random.Effects.Intn(data.CurrLevel.Metadata.Height/4))
					target := t.Object.Pos
					target.X += (random.Effects.Float64() - 0.5) * world.TileSize
					target.Y += (random.Effects.Float64() - 0.5) * world.TileSize
					SlugSpawner(pos, target, data.TileMetadata{}, false)
				}
				data.BossCounter1 = 10
			} else if data.BossCounter1 == 10 {
				if !b.Entity.HasComponent(myecs.Interpolation) {
					dur := constants.BlobFadeTime * float64(constants.Configuration.Gameplay.FrameRate) / 60.
					interY := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&b.HeadPos.Y).
						SetGween(b.HeadPos.Y, -1., dur, ease.Linear)
					b.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interY})
					data.BossCounter1++
				}
			} else if data.BossCounter1 == 11 {
				if !b.Entity.HasComponent(myecs.Interpolation) {
					b.State = data.BossDefeated
				}
			}
		}
	case data.BossDefeated:
		return
	case data.BossPreview:
		b.UTime += float32(timing.DT) * float32(constants.Configuration.Gameplay.FrameRate) / 60.
		b.HeadPos = BlobHeadPos(b.UTime, pixel.V(0.5, 0.65))
		b.UHeadPos[0] = float32(b.HeadPos.X)
		b.UHeadPos[1] = float32(b.HeadPos.Y)
		b.Object.SetPos(pixel.V(b.HeadPos.X*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().W(), b.HeadPos.Y*data.CurrentPlayArea.BackgroundView.Canvas.Bounds().H()))
	}
}

func BlobHeadPos(u float32, offset pixel.Vec) pixel.Vec {
	return pixel.V(offset.X+math.Sin(float64(u))*0.125, offset.Y+math.Cos(float64(u))*0.15)
}

func (b *Blob) Reset() {
	b.State = data.BossWaiting
	b.SlugCount = b.StartCount
	data.BossCounter0 = 0
}

func (b *Blob) IsDefeated() bool {
	return b.State == data.BossDefeated
}

func (b *Blob) Destroy() {
	if data.CurrLevel != nil {
		data.CurrLevel.Paused = false
	}
	myecs.Manager.DisposeEntity(b.Entity)
	data.CurrentPlayArea.BackgroundView.Canvas = pixelgl.NewCanvas(pixel.R(0, 0, 0, 0))
	UpdateViews()
}

func GetSlugRegenTile() *data.Tile {
	var tiles []*data.Tile
	for _, row := range data.CurrLevel.Tiles.T {
		for _, tile := range row {
			if tile.IsEmpty() && tile.Block == data.BlockSlugRegen {
				tiles = append(tiles, tile)
			}
		}
	}
	return GetBestRegenTile(tiles)
}

func SlugSpawner(pos pixel.Vec, target pixel.Vec, md data.TileMetadata, spawn bool) {
	obj := object.New().WithID("slug_spawner").SetPos(pos)
	obj.SetRect(pixel.R(0, 0, 16, 16))
	//obj.Flip = pos.X > target.X
	obj.Flop = random.Effects.Intn(2) == 0
	obj.Layer = 8
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	e.AddComponent(myecs.Drawable, img.NewSprite("blob_glob", constants.TileBatch))
	e.AddComponent(myecs.Update, data.NewFn(func() {
		if reanimator.FrameSwitch {
			lastPos := obj.LastPos
			newPos := obj.Pos
			if lastPos != newPos {
				angle := newPos.Sub(lastPos).Angle()
				obj.Rot = angle
			}
			if util.Magnitude(target.Sub(obj.Pos)) < 1. || !e.HasComponent(myecs.Interpolation) {
				if spawn {
					slugMD := data.TileMetadata{
						Flipped:     md.Flipped,
						Regenerate:  false,
						Orientation: md.Orientation,
					}
					slug := SlugCharacter(target, slugMD)
					slug.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if slug.State == data.Dead {
							data.BossCounter0--
							myecs.Manager.DisposeEntity(slug.Entity)
						}
					}))
				}
				myecs.Manager.DisposeEntity(e)
			}
		}
	}))
	interX := object.NewInterpolation(object.InterpolateX).
		SetGween(obj.Pos.X, target.X, constants.BlobSpawnSpeed, ease.Linear)
	diff := obj.Pos.Y - target.Y
	var midY float64
	if diff >= 0. {
		midY = obj.Pos.Y + world.TileSize + (diff / 4.)
	} else {
		midY = target.Y + world.TileSize + (-diff / 4.)
	}
	interY := object.NewInterpolation(object.InterpolateY).
		SetGween(obj.Pos.Y, midY, constants.BlobSpawnSpeed/2, ease.OutQuad).
		AddGween(midY, target.Y, constants.BlobSpawnSpeed/2, ease.InQuad)
	e.AddComponent(myecs.Interpolation, []*object.Interpolation{interX, interY})
}
