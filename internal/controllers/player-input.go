package controllers

import (
	"gemrunner/internal/data"
	"github.com/bytearena/ecs"
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
	StowKey   string
	ActionKey string
	Entity    *ecs.Entity

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
	lift := pi.Input.Get(pi.StowKey)
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
		}
	}
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
		actions.Stow = true
		lift.Consume()
	}
	pi.LastActions = actions
	return actions
}

func (pi *PlayerInput) GetEntity() *ecs.Entity {
	return pi.Entity
}

func NewPlayerInput(in *pxginput.Input, e *ecs.Entity) *PlayerInput {
	return &PlayerInput{
		Input:     in,
		LeftKey:   "left",
		RightKey:  "right",
		UpKey:     "up",
		DownKey:   "down",
		JumpKey:   "jump",
		PickUpKey: "pickUp",
		StowKey:   "stow",
		ActionKey: "action",
		Entity:    e,
	}
}
