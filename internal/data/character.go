package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"github.com/bytearena/ecs"
)

type Player int

type Dynamic struct {
	Object  *object.Object
	FauxObj *object.Object
	Anim    *reanimator.Tree
	Entity  *ecs.Entity
	Held    *ecs.Entity
	HeldObj *object.Object

	Actions Actions
	Flags   Flags
	Vars    Vars
	ATimer  *timing.Timer
	//BTimer  *timing.Timer
	LastTile *Tile
	MoveType MoveType
	Player   Player
}

func NewDynamic() *Dynamic {
	return &Dynamic{
		Player: Player(-1),
	}
}

type Actions struct {
	Left   bool
	Right  bool
	Up     bool
	Down   bool
	Jump   bool
	PickUp bool
	Action bool
}

type Vars struct {
	WalkSpeed    float64
	LeapSpeed    float64
	ClimbSpeed   float64
	SlideSpeed   float64
	LeapDelay    float64
	Gravity      float64
	HiJumpVSpeed float64
	HiJumpHSpeed float64
	HiJumpTimer  float64
	LgJumpVSpeed float64
	LgJumpHSpeed float64
	LgJumpTimer  float64
	IdleFreq     int
}

type Flags struct {
	LeftWall   bool
	RightWall  bool
	Ceiling    bool
	Floor      bool
	CanRun     bool
	OnLadder   bool
	GoingUp    bool
	Climbed    bool
	LadderDown bool
	LadderHere bool
	LeapOn     bool
	LeapOff    bool
	LeapTo     bool
	Breath     bool
	HighJump   bool
	LongJump   bool
	JumpR      bool
	JumpL      bool
	DropDown   bool
	Action     bool
	PickUp     bool
	Drop       bool
	HoldUp     bool
	HoldSide   bool
	HeldFlip   bool
	HeldNFlip  bool
	Hit        bool
	Dead       bool
	Attack     bool
}

type Controller interface {
	GetActions() Actions
}

type MoveType int

const (
	Humanoid = iota
	Flyer
)

func PlayerVars() Vars {
	return Vars{
		WalkSpeed:    constants.PlayerWalkSpeed,
		LeapSpeed:    constants.PlayerLeapSpeed,
		ClimbSpeed:   constants.PlayerClimbSpeed,
		SlideSpeed:   constants.PlayerSlideSpeed,
		LeapDelay:    constants.PlayerLeapDelay,
		Gravity:      constants.NormalGravity,
		HiJumpVSpeed: constants.PlayerHighJumpSpeed,
		HiJumpHSpeed: constants.PlayerHighJumpHSpeed,
		HiJumpTimer:  constants.PlayerHighJumpTimer,
		LgJumpVSpeed: constants.PlayerLongJumpVSpeed,
		LgJumpHSpeed: constants.PlayerLongJumpHSpeed,
		LgJumpTimer:  constants.PlayerLongJumpTimer,
		IdleFreq:     constants.IdleFrequency,
	}
}

func DemonVars() Vars {
	return Vars{
		WalkSpeed:    constants.DemonWalkSpeed,
		LeapSpeed:    constants.DemonLeapSpeed,
		ClimbSpeed:   constants.DemonClimbSpeed,
		SlideSpeed:   constants.DemonSlideSpeed,
		LeapDelay:    constants.DemonLeapDelay,
		Gravity:      constants.DemonGravity,
		HiJumpVSpeed: constants.DemonHighJumpSpeed,
		HiJumpHSpeed: constants.DemonHighJumpHSpeed,
		HiJumpTimer:  constants.DemonHighJumpTimer,
		LgJumpVSpeed: constants.DemonLongJumpVSpeed,
		LgJumpHSpeed: constants.DemonLongJumpHSpeed,
		LgJumpTimer:  constants.DemonLongJumpTimer,
		IdleFreq:     constants.IdleFrequency,
	}
}
