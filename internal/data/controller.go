package data

import "github.com/bytearena/ecs"

type Controller interface {
	GetActions() Actions
	ClearPrev()
	GetEntity() *ecs.Entity
}

type Direction int

const (
	NoDirection = iota
	Left
	Right
	Up
	Down
)

func (d Direction) String() string {
	switch d {
	case Left:
		return "Left"
	case Right:
		return "Right"
	case Up:
		return "Up"
	case Down:
		return "Down"
	default:
		return "NoDirection"
	}
}

type Actions struct {
	Direction     Direction `json:"direction,omitempty"`
	PrevDirection Direction `json:"prevDirection,omitempty"`
	PickUp        bool      `json:"pickUp,omitempty"`
	Action        bool      `json:"action,omitempty"`
	DigLeft       bool      `json:"digLeft,omitempty"`
	DigRight      bool      `json:"digRight,omitempty"`
}

func NewAction() Actions {
	return Actions{
		Direction:     NoDirection,
		PrevDirection: NoDirection,
	}
}

func (a Actions) Up() bool {
	return a.Direction == Up || (a.PrevDirection == Up && a.Direction != Down)
}

func (a Actions) Down() bool {
	return a.Direction == Down || (a.PrevDirection == Down && a.Direction != Up)
}

func (a Actions) Left() bool {
	return a.Direction == Left || (a.PrevDirection == Left && a.Direction != Right)
}

func (a Actions) Right() bool {
	return a.Direction == Right || (a.PrevDirection == Right && a.Direction != Left)
}

func (a Actions) Any() bool {
	return a.Direction != NoDirection || a.PickUp || a.Action || a.DigLeft || a.DigRight
}

func (a Actions) Copy() Actions {
	return Actions{
		Direction:     a.Direction,
		PrevDirection: a.PrevDirection,
		PickUp:        a.PickUp,
		Action:        a.Action,
		DigLeft:       a.DigLeft,
		DigRight:      a.DigRight,
	}
}

type LevelReplay struct {
	PuzzleSet  string        `json:"puzzleSet"`
	Filename   string        `json:"filename"`
	ReplayFile string        `json:"replayFile"`
	PuzzleNum  int           `json:"puzzle"`
	Seed       int64         `json:"seed"`
	Frames     []ReplayFrame `json:"replayFrames"`
	FrameIndex int           `json:"-"`
}

type ReplayFrame struct {
	Frame     int      `json:"frame"`
	P1Actions *Actions `json:"p1Actions,omitempty"`
	P2Actions *Actions `json:"p2Actions,omitempty"`
	P3Actions *Actions `json:"p3Actions,omitempty"`
	P4Actions *Actions `json:"p4Actions,omitempty"`
}
