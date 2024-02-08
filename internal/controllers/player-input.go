package controllers

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/timing"
)

type PlayerInput struct {
	LeftKey     string
	RightKey    string
	UpKey       string
	DownKey     string
	JumpKey     string
	ActionKey   string
	Stack       []string
	Buffer      string
	BufferTimer *timing.Timer
}

func (pi *PlayerInput) GetActions() data.Actions {
	actions := data.Actions{}
	// get all inputs
	left := data.GameInput.Get(pi.LeftKey)
	right := data.GameInput.Get(pi.RightKey)
	up := data.GameInput.Get(pi.UpKey)
	down := data.GameInput.Get(pi.DownKey)
	jump := data.GameInput.Get(pi.JumpKey)
	action := data.GameInput.Get(pi.ActionKey)
	// remove any missing from the stack
	for i := len(pi.Stack) - 1; i >= 0; i-- {
		if !data.GameInput.Get(pi.Stack[i]).Pressed() {
			if len(pi.Stack) > 1 {
				pi.Stack = append(pi.Stack[:i], pi.Stack[i+1:]...)
			} else {
				pi.Stack = []string{}
			}
		}
	}
	// add any newly pressed inputs to the stack, put top input into buffer
	if left.JustPressed() {
		pi.Stack = append(pi.Stack, pi.LeftKey)
	} else if right.JustPressed() {
		pi.Stack = append(pi.Stack, pi.RightKey)
	} else if up.JustPressed() {
		pi.Stack = append(pi.Stack, pi.UpKey)
	} else if down.JustPressed() {
		pi.Stack = append(pi.Stack, pi.DownKey)
	}
	if len(pi.Stack) > 0 {
		switch pi.Stack[len(pi.Stack)-1] {
		case pi.LeftKey:
			actions.Left = true
		case pi.RightKey:
			actions.Right = true
		case pi.UpKey:
			actions.Up = true
		case pi.DownKey:
			actions.Down = true
		}
		//if pi.Buffer != "" {
		//	if pi.BufferTimer == nil {
		//		pi.BufferTimer = timing.New(0.25)
		//	}
		//	if pi.BufferTimer.UpdateDone() {
		//		pi.Buffer = ""
		//	} else {
		//		switch pi.Buffer {
		//		case pi.LeftKey:
		//			actions.Left = true
		//		case pi.RightKey:
		//			actions.Right = true
		//		case pi.UpKey:
		//			actions.Up = true
		//		case pi.DownKey:
		//			actions.Down = true
		//		}
		//	}
		//}
	} else {
		pi.Buffer = ""
		pi.BufferTimer = nil
	}
	// Basic Controller
	//if left.Pressed() {
	//	actions.Left = true
	//	if left.JustPressed() {
	//		right.Consume()
	//		up.Consume()
	//		down.Consume()
	//	}
	//}
	//if right.Pressed() {
	//	actions.Right = true
	//	if right.JustPressed() {
	//		left.Consume()
	//		up.Consume()
	//		down.Consume()
	//	}
	//}
	//if up.Pressed() {
	//	actions.Up = true
	//	if up.JustPressed() {
	//		left.Consume()
	//		right.Consume()
	//		down.Consume()
	//	}
	//}
	//if down.Pressed() {
	//	actions.Down = true
	//	if down.JustPressed() {
	//		left.Consume()
	//		right.Consume()
	//		up.Consume()
	//	}
	//}
	if jump.JustPressed() {
		actions.Jump = true
		jump.Consume()
	}
	if action.JustPressed() {
		actions.Action = true
		action.Consume()
	}
	return actions
}

func NewPlayerInput(l, r, u, d, j, a string) *PlayerInput {
	return &PlayerInput{
		LeftKey:   l,
		RightKey:  r,
		UpKey:     u,
		DownKey:   d,
		JumpKey:   j,
		ActionKey: a,
	}
}
