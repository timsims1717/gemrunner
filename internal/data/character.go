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
	Anim    *reanimator.Tree
	Entity  *ecs.Entity
	Held    *ecs.Entity
	HeldObj *object.Object
	Control Controller

	Actions  Actions
	State    CharacterState
	Flags    Flags
	Vars     Vars
	ATimer   *timing.Timer
	ACounter int
	//BTimer  *timing.Timer
	LastTile *Tile
	MoveType MoveType
	Player   Player
	Color    string
}

func NewDynamic() *Dynamic {
	return &Dynamic{
		Player:  Player(-1),
		Actions: NewAction(),
		Flags: Flags{
			Floor: true,
		},
		Vars: Vars{
			Gravity:  constants.NormalGravity,
			IdleFreq: constants.IdleFrequency,
		},
	}
}

type Direction int

const (
	Left = iota
	Right
	Up
	Down
	None
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
		return "None"
	}
}

type Actions struct {
	Direction     Direction
	PrevDirection Direction
	Jump          bool
	PickUp        bool
	Action        bool
}

func NewAction() Actions {
	return Actions{
		Direction:     None,
		PrevDirection: None,
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
	HiJumpCntr   float64
	LgJumpVSpeed float64
	LgJumpHSpeed float64
	LgJumpTimer  float64
	LgJumpCntr   float64
	IdleFreq     int
}

type CharacterState int

const (
	Grounded = iota
	Ladder
	Falling
	Jumping
	Leaping
	Flying
	Attack
	Hit
	Dead
)

func (s CharacterState) String() string {
	switch s {
	case Grounded:
		return "Grounded"
	case Ladder:
		return "Ladder"
	case Falling:
		return "Falling"
	case Jumping:
		return "Jumping"
	case Leaping:
		return "Leaping"
	case Flying:
		return "Flying"
	case Attack:
		return "Attack"
	case Hit:
		return "Hit"
	case Dead:
		return "Dead"
	}
	return "Unknown"
}

type Flags struct {
	LeftWall   bool
	RightWall  bool
	Ceiling    bool
	Floor      bool
	GoingUp    bool
	Climbed    bool
	LeapOn     bool
	LeapOff    bool
	LeapTo     bool
	Breath     bool
	HighJump   bool
	LongJump   bool
	JumpR      bool
	JumpL      bool
	Action     bool
	PickUp     bool
	Drop       bool
	HoldSwitch bool
	HoldUp     bool
	HoldSide   bool
	HeldFlip   bool
	HeldNFlip  bool
	Hit        bool
	Dead       bool
	Attack     bool
	Flying     bool
	Frame      bool
	JumpBuff   int
	PickUpBuff int
	ActionBuff int
}

type Controller interface {
	GetActions() Actions
	ClearPrev()
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
		HiJumpCntr:   constants.PlayerHighJumpCounter,
		LgJumpVSpeed: constants.PlayerLongJumpVSpeed,
		LgJumpHSpeed: constants.PlayerLongJumpHSpeed,
		LgJumpTimer:  constants.PlayerLongJumpTimer,
		LgJumpCntr:   constants.PlayerLongJumpCounter,
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
		HiJumpCntr:   constants.DemonHighJumpCounter,
		LgJumpVSpeed: constants.DemonLongJumpVSpeed,
		LgJumpHSpeed: constants.DemonLongJumpHSpeed,
		LgJumpTimer:  constants.DemonLongJumpTimer,
		LgJumpCntr:   constants.DemonLongJumpCounter,
		IdleFreq:     constants.IdleFrequency,
	}
}

func FlyVars() Vars {
	return Vars{
		WalkSpeed: constants.FlySpeed,
	}
}
