package controllers

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/reanimator"
	"github.com/bytearena/ecs"
)

type ReplayController struct {
	Player int
	Replay *data.LevelReplay
	Entity *ecs.Entity
}

func NewReplayController(replay *data.LevelReplay, p int, e *ecs.Entity) *ReplayController {
	return &ReplayController{
		Player: p,
		Replay: replay,
		Entity: e,
	}
}

func (rc *ReplayController) ClearPrev() {}

func (rc *ReplayController) GetActions() data.Actions {
	var actions data.Actions
	if reanimator.FrameSwitch {
		if rc.Replay.FrameIndex < len(rc.Replay.Frames) {
			frame := rc.Replay.Frames[rc.Replay.FrameIndex]
			if frame.Frame == data.CurrLevel.FrameNumber {
				switch rc.Player {
				case 0:
					actions = *frame.P1Actions
				case 1:
					actions = *frame.P2Actions
				case 2:
					actions = *frame.P3Actions
				case 3:
					actions = *frame.P4Actions
				}
				rc.Replay.FrameIndex++
			}
		}
	}
	return actions
}

func (rc *ReplayController) GetEntity() *ecs.Entity {
	return rc.Entity
}
