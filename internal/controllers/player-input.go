package controllers

import (
	"gemrunner/internal/data"
	pxginput "github.com/timsims1717/pixel-go-input"
)

type PlayerInput struct {
	Input     *pxginput.Input
	LeftKey   string
	RightKey  string
	UpKey     string
	DownKey   string
	JumpKey   string
	PickUpKey string
	LiftKey   string
	ActionKey string

	LastActions data.Actions
}

func (pi *PlayerInput) ClearPrev() {
	pi.LastActions.PrevDirection = data.None
}

func (pi *PlayerInput) GetActions() data.Actions {
	actions := data.NewAction()
	actions.Direction = pi.LastActions.Direction
	actions.PrevDirection = pi.LastActions.PrevDirection
	// get all inputs
	left := pi.Input.Get(pi.LeftKey)
	right := pi.Input.Get(pi.RightKey)
	up := pi.Input.Get(pi.UpKey)
	down := pi.Input.Get(pi.DownKey)
	jump := pi.Input.Get(pi.JumpKey)
	pickUp := pi.Input.Get(pi.PickUpKey)
	lift := pi.Input.Get(pi.LiftKey)
	action := pi.Input.Get(pi.ActionKey)
	if !left.Pressed() && !right.Pressed() && !up.Pressed() && !down.Pressed() {
		actions.Direction = data.None
		actions.PrevDirection = data.None
	} else {
		if left.JustPressed() {
			actions.PrevDirection = actions.Direction
			actions.Direction = data.Left
		} else if right.JustPressed() {
			actions.PrevDirection = actions.Direction
			actions.Direction = data.Right
		} else if up.JustPressed() {
			actions.PrevDirection = actions.Direction
			actions.Direction = data.Up
		} else if down.JustPressed() {
			actions.PrevDirection = actions.Direction
			actions.Direction = data.Down
		} else {
			if (!left.Pressed() && actions.Direction == data.Left) ||
				(!right.Pressed() && actions.Direction == data.Right) ||
				(!up.Pressed() && actions.Direction == data.Up) ||
				(!down.Pressed() && actions.Direction == data.Down) {
				actions.Direction = actions.PrevDirection
			}
			if actions.Direction == data.None {
				if up.Pressed() {
					actions.Direction = data.Up
				} else if down.Pressed() {
					actions.Direction = data.Down
				} else if left.Pressed() {
					actions.Direction = data.Left
				} else if right.Pressed() {
					actions.Direction = data.Right
				}
			}
			//if actions.PrevDirection == data.None {
			//	if up.Pressed() && actions.Direction != data.Up {
			//		actions.PrevDirection = data.Up
			//	} else if down.Pressed() && actions.Direction != data.Down {
			//		actions.PrevDirection = data.Down
			//	} else if left.Pressed() && actions.Direction != data.Left {
			//		actions.PrevDirection = data.Left
			//	} else if right.Pressed() && actions.Direction != data.Right {
			//		actions.PrevDirection = data.Right
			//	}
			//}
		}
	}
	// set any missing directional priority to zero
	// find the highest directional priority
	//high := 0
	//anyPressed := false
	//if !left.Pressed() {
	//	pi.LeftPriority = 0
	//} else {
	//	high = pi.LeftPriority
	//	anyPressed = true
	//}
	//if !right.Pressed() {
	//	pi.RightPriority = 0
	//} else if high < pi.RightPriority {
	//	high = pi.RightPriority
	//	anyPressed = true
	//}
	//if !up.Pressed() {
	//	pi.UpPriority = 0
	//} else if high < pi.UpPriority {
	//	high = pi.UpPriority
	//	anyPressed = true
	//}
	//if !down.Pressed() {
	//	pi.DownPriority = 0
	//} else if high < pi.DownPriority {
	//	high = pi.DownPriority
	//	anyPressed = true
	//}
	//// if any direction just pressed, set that priority to high+1
	//if left.JustPressed() {
	//	pi.LeftPriority = high + 1
	//} else if right.JustPressed() {
	//	pi.RightPriority = high + 1
	//} else if up.JustPressed() {
	//	pi.UpPriority = high + 1
	//} else if down.JustPressed() {
	//	pi.DownPriority = high + 1
	//}
	//// if any directions are released, put one into the LastPressed string
	//if anyPressed {
	//	if left.JustReleased() {
	//		pi.LastPressed = pi.LeftKey
	//	} else if right.JustReleased() {
	//		pi.LastPressed = pi.RightKey
	//	} else if up.JustReleased() {
	//		pi.LastPressed = pi.UpKey
	//	} else if down.JustReleased() {
	//		pi.LastPressed = pi.DownKey
	//	}
	//} else {
	//	pi.LastPressed = ""
	//}
	//// Assign up to two directions
	//assigned := 0
	//if pi.LeftPriority > 0 && pi.LeftPriority > pi.RightPriority {
	//	actions.Left = true
	//	assigned++
	//} else if pi.RightPriority > 0 && pi.RightPriority > pi.LeftPriority {
	//	actions.Right = true
	//	assigned++
	//}
	//if pi.UpPriority > 0 && pi.UpPriority > pi.DownPriority {
	//	actions.Up = true
	//	assigned++
	//} else if pi.DownPriority > 0 && pi.DownPriority > pi.UpPriority {
	//	actions.Down = true
	//	assigned++
	//}
	//if assigned < 2 && pi.LastPressed != "" {
	//	switch pi.LastPressed {
	//	case pi.LeftKey:
	//		actions.Left = true
	//	case pi.RightKey:
	//		actions.Right = true
	//	case pi.UpKey:
	//		actions.Up = true
	//	case pi.DownKey:
	//		actions.Down = true
	//	}
	//}
	// add any newly pressed inputs to the stack, put top input into buffer
	//if left.JustPressed() {
	//	pi.Stack = append(pi.Stack, pi.LeftKey)
	//} else if right.JustPressed() {
	//	pi.Stack = append(pi.Stack, pi.RightKey)
	//} else if up.JustPressed() {
	//	pi.Stack = append(pi.Stack, pi.UpKey)
	//} else if down.JustPressed() {
	//	pi.Stack = append(pi.Stack, pi.DownKey)
	//}
	//if len(pi.Stack) > 0 {
	//	switch pi.Stack[len(pi.Stack)-1] {
	//	case pi.LeftKey:
	//		actions.Left = true
	//	case pi.RightKey:
	//		actions.Right = true
	//	case pi.UpKey:
	//		actions.Up = true
	//	case pi.DownKey:
	//		actions.Down = true
	//	}
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
	//} else {
	//	pi.Buffer = ""
	//	pi.BufferTimer = nil
	//}
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
	if pickUp.JustPressed() {
		actions.PickUp = true
		pickUp.Consume()
	}
	if action.JustPressed() {
		actions.Action = true
		action.Consume()
	}
	if lift.JustPressed() {
		actions.Lift = true
		lift.Consume()
	}
	pi.LastActions = actions
	return actions
}

func NewPlayerInput(in *pxginput.Input) *PlayerInput {
	return &PlayerInput{
		Input:     in,
		LeftKey:   "left",
		RightKey:  "right",
		UpKey:     "up",
		DownKey:   "down",
		JumpKey:   "jump",
		PickUpKey: "pickUp",
		LiftKey:   "lift",
		ActionKey: "action",
	}
}
