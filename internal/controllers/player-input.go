package controllers

import (
	"gemrunner/internal/data"
	"github.com/bytearena/ecs"
	pxginput "github.com/timsims1717/pixel-go-input"
)

type PlayerInput struct {
	Input       *pxginput.Input
	LeftKey     string
	RightKey    string
	UpKey       string
	DownKey     string
	JumpKey     string
	PickUpKey   string
	ActionKey   string
	DigLeftKey  string
	DigRightKey string
	Entity      *ecs.Entity

	LastActions data.Actions
}

func (pi *PlayerInput) ClearPrev() {
	pi.LastActions.PrevDirection = data.NoDirection
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
	pickUp := pi.Input.Get(pi.PickUpKey)
	action := pi.Input.Get(pi.ActionKey)
	digLeft := pi.Input.Get(pi.DigLeftKey)
	digRight := pi.Input.Get(pi.DigRightKey)
	if !left.Pressed() && !right.Pressed() && !up.Pressed() && !down.Pressed() {
		actions.Direction = data.NoDirection
		actions.PrevDirection = data.NoDirection
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
			if actions.Direction == data.NoDirection {
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
		}
	}
	if pickUp.JustPressed() {
		actions.PickUp = true
		pickUp.Consume()
	}
	if action.Pressed() {
		actions.Action = true
	}
	if digLeft.Pressed() {
		actions.DigLeft = true
	}
	if digRight.Pressed() {
		actions.DigRight = true
	}
	pi.LastActions = actions
	return actions
}

func (pi *PlayerInput) GetEntity() *ecs.Entity {
	return pi.Entity
}

func NewPlayerInput(in *pxginput.Input, e *ecs.Entity) *PlayerInput {
	return &PlayerInput{
		Input:       in,
		LeftKey:     "left",
		RightKey:    "right",
		UpKey:       "up",
		DownKey:     "down",
		JumpKey:     "jump",
		PickUpKey:   "pickUp",
		ActionKey:   "action",
		DigLeftKey:  "digLeft",
		DigRightKey: "digRight",
		Entity:      e,
	}
}
