package controllers

import (
	"gemrunner/internal/data"
)

type BackAndForth struct {
	Ch   *data.Dynamic
	Left bool
}

func NewBackAndForth(dyn *data.Dynamic, left bool) *BackAndForth {
	return &BackAndForth{
		Ch:   dyn,
		Left: left,
	}
}

func (bf *BackAndForth) ClearPrev() {}

func (bf *BackAndForth) GetActions() data.Actions {
	actions := data.NewAction()
	if bf.Ch.Flags.RightWall {
		bf.Left = true
		actions.Direction = data.Left
	} else if bf.Ch.Flags.LeftWall {
		bf.Left = false
		actions.Direction = data.Right
	} else if bf.Left {
		actions.Direction = data.Left
	} else {
		actions.Direction = data.Right
	}
	return actions
}
