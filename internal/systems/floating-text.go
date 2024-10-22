package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

func FloatingTextSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsText) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tf, okTF := result.Components[myecs.Text].(*data.FloatingText)
		if okO && okTF {
			if reanimator.FrameSwitch {
				if int(tf.Text.Obj.Rect.W())%2 != 0 {
					tf.Text.Obj.Offset.X = -0.5
					tf.Shadow.Obj.Offset.X = -0.5
				} else {
					tf.Text.Obj.Offset.X = 0
					tf.Shadow.Obj.Offset.X = 0
				}
				// bob
				if tf.Bob && !tf.Text.Obj.Hidden {
					tf.BobCounter++
					if tf.BobCounter%constants.BobInterval == 0 {
						if tf.Text.Obj.Offset.Y == 0. {
							tf.Text.Obj.Offset.Y = 1.
							tf.Shadow.Obj.Offset.Y = 1.
						} else {
							tf.Text.Obj.Offset.Y = 0.
							tf.Shadow.Obj.Offset.Y = 0.
						}
						tf.BobCounter = 0
					}
				} else {
					tf.Text.Obj.Offset.Y = 0.
					tf.Shadow.Obj.Offset.Y = 0.
				}
				// prox
				if tf.Prox && (!tf.Temp || tf.Text.Obj.Hidden) {
					show := false
					for _, resultP := range myecs.Manager.Query(myecs.IsPlayer) {
						if po, okP := resultP.Components[myecs.Object].(*object.Object); okP {
							if util.Magnitude(tf.Pos.Sub(po.Pos)) < constants.TextProxDist {
								show = true
								break
							}
						}
					}
					if show {
						tf.Show()
						tf.ProxCounter = 0
					} else {
						tf.ProxCounter++
						if tf.ProxCounter > constants.TextProxBuffer {
							tf.Hide()
							tf.ProxCounter = 0
						}
					}
				}
				// temp
				if !tf.Text.Obj.Hidden && tf.Temp {
					tf.TempCounter++
					if tf.TempCounter > tf.Timer*constants.TextTimer {
						if tf.Tile == nil {
							myecs.Manager.DisposeEntity(tf.Entity)
							myecs.Manager.DisposeEntity(tf.ShEntity)
							myecs.Manager.DisposeEntity(result.Entity)
						} else {
							tf.Text.Hide()
							tf.Shadow.Hide()
						}
					}
				}
			}
		}
	}
}

func FloatingTextEditorSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsText) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tf, okTF := result.Components[myecs.Text].(*data.FloatingText)
		if okO && okTF {
			if int(tf.Text.Obj.Rect.W())%2 != 0 {
				tf.Text.Obj.Offset.X = -0.5
				tf.Shadow.Obj.Offset.X = -0.5
			} else {
				tf.Text.Obj.Offset.X = 0
				tf.Shadow.Obj.Offset.X = 0
			}
			tf.ProxCounter = 0
			tf.TempCounter = 0
			tf.Text.Show()
			if tf.HasShadow {
				tf.Shadow.Show()
			} else {
				tf.Shadow.Hide()
			}

			// bob
			if reanimator.FrameSwitch {
				if tf.Bob && !tf.Text.Obj.Hidden {
					tf.BobCounter++
					if tf.BobCounter%constants.BobInterval == 0 {
						if tf.Text.Obj.Offset.Y == 0. {
							tf.Text.Obj.Offset.Y = 1.
							tf.Shadow.Obj.Offset.Y = 1.
						} else {
							tf.Text.Obj.Offset.Y = 0.
							tf.Shadow.Obj.Offset.Y = 0.
						}
						tf.BobCounter = 0
					}
				} else {
					tf.Text.Obj.Offset.Y = 0.
					tf.Shadow.Obj.Offset.Y = 0.
				}
			}

			if data.Editor.Mode == data.ModeText && tf.Tile != nil {
				data.IMDraw.Color = tf.ShadowCol
				data.IMDraw.EndShape = imdraw.SharpEndShape
				data.IMDraw.Push(tf.Tile.Object.Pos.Sub(pixel.V(-world.HalfSize, world.HalfSize)), tf.Tile.Object.Pos.Sub(pixel.V(world.HalfSize, world.HalfSize)))
				data.IMDraw.Push(tf.Tile.Object.Pos.Sub(pixel.V(world.HalfSize, world.HalfSize)), tf.Tile.Object.Pos.Sub(pixel.V(world.HalfSize, -world.HalfSize)))
				data.IMDraw.Push(tf.Tile.Object.Pos.Sub(pixel.V(world.HalfSize, -world.HalfSize)), tf.Tile.Object.Pos.Sub(pixel.V(-world.HalfSize, -world.HalfSize)))
				data.IMDraw.Push(tf.Tile.Object.Pos.Sub(pixel.V(-world.HalfSize, -world.HalfSize)), tf.Tile.Object.Pos.Sub(pixel.V(-world.HalfSize, world.HalfSize)))
				data.IMDraw.Line(2)
			}
		}
	}
}

func FloatingTextStartLevel() {
	for _, result := range myecs.Manager.Query(myecs.IsText) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		tf, okTF := result.Components[myecs.Text].(*data.FloatingText)
		if okO && okTF {
			if tf.Prox {
				tf.Hide()
			}
			tf.BobCounter = random.Effects.Intn(16) - 8
			tf.ProxCounter = 0
			tf.TempCounter = 0
		}
	}
}
