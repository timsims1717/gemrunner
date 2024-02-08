package data

import (
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"github.com/bytearena/ecs"
)

type Character struct {
	Object  *object.Object
	FauxObj *object.Object
	Anim    *reanimator.Tree
	Entity  *ecs.Entity

	Actions Actions
	Flags   Flags

	PlayerIndex int
}

type Actions struct {
	Left   bool
	Right  bool
	Up     bool
	Down   bool
	Jump   bool
	Action bool
}

type Flags struct {
	LeftWall   bool
	RightWall  bool
	Ceiling    bool
	Floor      bool
	CanRun     bool
	OnLadder   bool
	GoingUp    bool
	WentUp     bool
	LadderDown bool
	LadderHere bool
}

type Controller interface {
	GetActions() Actions
}
