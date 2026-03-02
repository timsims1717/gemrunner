package animations

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/data/death"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func SlugAnimation(ch *data.Dynamic) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, "slug_idle", reanimator.Hold)
	move := reanimator.NewBatchAnimation("move", batch, "slug_move", reanimator.Loop)
	corner := reanimator.NewBatchAnimation("corner", batch, "slug_corner", reanimator.Tran)
	edge := reanimator.NewBatchAnimation("edge", batch, "slug_edge", reanimator.Tran)
	regen := reanimator.NewBatchAnimation("regen", batch, "slug_regen", reanimator.Tran)
	crush := reanimator.NewBatchAnimation("crush", batch, "slug_crush", reanimator.Tran)
	explode := reanimator.NewBatchAnimation("explode", batch, "slug_explode", reanimator.Tran)

	corner.SetEndTrigger(func() {
		ch.Flags.NextStep = true
		switch ch.Flags.Orientation {
		case data.Up:
			if ch.Object.Flip {
				ch.Flags.Orientation = data.Right
			} else {
				ch.Flags.Orientation = data.Left
			}
		case data.Left:
			if ch.Object.Flip {
				ch.Flags.Orientation = data.Up
			} else {
				ch.Flags.Orientation = data.Down
			}
		case data.Right:
			if ch.Object.Flip {
				ch.Flags.Orientation = data.Down
			} else {
				ch.Flags.Orientation = data.Up
			}
		default:
			if ch.Object.Flip {
				ch.Flags.Orientation = data.Left
			} else {
				ch.Flags.Orientation = data.Right
			}
		}
	})
	edge.SetTriggerCAll(func(a *reanimator.Anim, _ string, _ int) {
		var dx, dy float64
		switch a.Step {
		case 0:
			switch ch.Flags.Orientation {
			case data.Left:
				if ch.Object.Flip {
					dy = 10
				} else {
					dy = -10
				}
				dx = 15
			case data.Right:
				if ch.Object.Flip {
					dy = -10
				} else {
					dy = 10
				}
				dx = -15
			case data.Up:
				if ch.Object.Flip {
					dx = 10
				} else {
					dx = -10
				}
				dy = -15
			default:
				if ch.Object.Flip {
					dx = -10
				} else {
					dx = 10
				}
				dy = 15
			}
		case 1:
			switch ch.Flags.Orientation {
			case data.Left:
				if ch.Object.Flip {
					dy = 11
				} else {
					dy = -11
				}
				dx = 12
			case data.Right:
				if ch.Object.Flip {
					dy = -11
				} else {
					dy = 11
				}
				dx = -12
			case data.Up:
				if ch.Object.Flip {
					dx = 11
				} else {
					dx = -11
				}
				dy = -12
			default:
				if ch.Object.Flip {
					dx = -11
				} else {
					dx = 11
				}
				dy = 12
			}
		case 2:
			switch ch.Flags.Orientation {
			case data.Left:
				if ch.Object.Flip {
					dy = 12
				} else {
					dy = -12
				}
				dx = 9
			case data.Right:
				if ch.Object.Flip {
					dy = -12
				} else {
					dy = 12
				}
				dx = -9
			case data.Up:
				if ch.Object.Flip {
					dx = 12
				} else {
					dx = -12
				}
				dy = -9
			default:
				if ch.Object.Flip {
					dx = -12
				} else {
					dx = 12
				}
				dy = 9
			}
		case 3:
			switch ch.Flags.Orientation {
			case data.Left:
				if ch.Object.Flip {
					dy = 14
				} else {
					dy = -14
				}
				dx = 7
			case data.Right:
				if ch.Object.Flip {
					dy = -14
				} else {
					dy = 14
				}
				dx = -7
			case data.Up:
				if ch.Object.Flip {
					dx = 14
				} else {
					dx = -14
				}
				dy = -7
			default:
				if ch.Object.Flip {
					dx = -14
				} else {
					dx = 14
				}
				dy = 7
			}
			switch ch.Flags.Orientation {
			case data.Up:
				if ch.Object.Flip {
					ch.Flags.Orientation = data.Left
				} else {
					ch.Flags.Orientation = data.Right
				}
			case data.Left:
				if ch.Object.Flip {
					ch.Flags.Orientation = data.Down
				} else {
					ch.Flags.Orientation = data.Up
				}
			case data.Right:
				if ch.Object.Flip {
					ch.Flags.Orientation = data.Up
				} else {
					ch.Flags.Orientation = data.Down
				}
			default:
				if ch.Object.Flip {
					ch.Flags.Orientation = data.Right
				} else {
					ch.Flags.Orientation = data.Left
				}
			}
			ch.Flags.NextStep = true
		}
		if ch.BelowTile != nil {
			ch.Object.SetPos(ch.BelowTile.Object.Pos.Add(pixel.V(dx, dy)))
		}
	})
	crush.SetEndTrigger(func() {
		ch.Flags.NextStep = true
	})
	explode.SetEndTrigger(func() {
		ch.Flags.NextStep = true
	})
	regen.SetEndTrigger(func() {
		ch.Flags.Regen = false
	})
	return reanimator.New().
		AddAnimation(idle).
		AddAnimation(move).
		AddAnimation(corner).
		AddAnimation(edge).
		AddAnimation(regen).
		AddAnimation(crush).
		AddAnimation(explode).
		AddNull("none").
		SetChooseFn(func() string {
			switch ch.State {
			case data.Hit:
				switch ch.Flags.Death {
				case death.Crushed:
					return "crush"
				case death.Exploded, death.Drowned:
					return "explode"
				default:
					return "none"
				}
			case data.Regen:
				return "regen"
			case data.Dead:
				return "none"
			case data.AroundCorner:
				return "edge"
			case data.Grounded:
				if !data.CurrLevel.Start {
					return "idle"
				}
				if ch.Flags.NextStep {
					ch.Flags.NextStep = false
					return "move"
				}
				cx, cy := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
				switch ch.Flags.Orientation {
				case data.Up:
					r := data.CurrLevel.Get(cx+1, cy)
					l := data.CurrLevel.Get(cx-1, cy)
					if (ch.Flags.RightWall && r.IsSolid() && ch.Object.Flip) ||
						(ch.Flags.LeftWall && l.IsSolid() && !ch.Object.Flip) {
						return "corner"
					} else if (ch.Flags.RightWall && ch.Object.Flip) ||
						(ch.Flags.LeftWall && !ch.Object.Flip) {
						return "idle"
					}
				case data.Left:
					u := data.CurrLevel.Get(cx, cy+1)
					d := data.CurrLevel.Get(cx, cy-1)
					if (ch.Flags.Ceiling && u.IsSolid() && ch.Object.Flip) ||
						(ch.Flags.Floor && d.IsSolid() && !ch.Object.Flip) {
						return "corner"
					} else if (ch.Flags.Ceiling && ch.Object.Flip) ||
						(ch.Flags.Floor && !ch.Object.Flip) {
						return "idle"
					}
				case data.Right:
					u := data.CurrLevel.Get(cx, cy+1)
					d := data.CurrLevel.Get(cx, cy-1)
					if (ch.Flags.Floor && d.IsSolid() && ch.Object.Flip) ||
						(ch.Flags.Ceiling && u.IsSolid() && !ch.Object.Flip) {
						return "corner"
					} else if (ch.Flags.Floor && ch.Object.Flip) ||
						(ch.Flags.Ceiling && !ch.Object.Flip) {
						return "idle"
					}
				default:
					r := data.CurrLevel.Get(cx+1, cy)
					l := data.CurrLevel.Get(cx-1, cy)
					if (ch.Flags.LeftWall && l.IsSolid() && ch.Object.Flip) ||
						(ch.Flags.RightWall && r.IsSolid() && !ch.Object.Flip) {
						return "corner"
					} else if (ch.Flags.LeftWall && ch.Object.Flip) ||
						(ch.Flags.RightWall && !ch.Object.Flip) {
						return "idle"
					}
				}
				if ch.Actions.Direction != data.NoDirection {
					return "move"
				}
				return "idle"
			default:
				return "idle"
			}
		})
}
