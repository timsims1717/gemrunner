package data

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
)

type Dynamic struct {
	Object *object.Object
	Anim   *reanimator.Tree
	Entity *ecs.Entity
	//Held      *ecs.Entity
	//HeldObj   *object.Object
	Inventory *ecs.Entity
	Control   Controller

	Actions  Actions
	State    CharacterState
	Flags    Flags
	Vars     Vars
	Options  CharacterOptions
	ACounter int
	LastTile *Tile
	MoveType MoveType
	Player   int
	Enemy    int
	Color    string
}

func NewDynamic() *Dynamic {
	return &Dynamic{
		Player:  -1,
		Enemy:   -1,
		Actions: NewAction(),
		Flags: Flags{
			Floor: true,
		},
		Vars: Vars{
			Gravity:  constants.NormalGravity,
			IdleFreq: constants.IdleFrequency,
		},
		Options: CharacterOptions{
			Regen: false,
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
	Gravity      float64
	HiJumpVSpeed float64
	HiJumpHSpeed float64
	HiJumpCntr   float64
	LgJumpVSpeed float64
	LgJumpHSpeed float64
	LgJumpCntr   float64
	IdleFreq     int
}

type CharacterState int

const (
	Grounded = iota
	OnLadder
	Falling
	Jumping
	Leaping
	Flying
	DoingAction
	Attack
	Hit
	Dying
	Dead
	Regen
)

type ItemAction int

const (
	NoItemAction = iota
	ThrowBox
)

func (s CharacterState) String() string {
	switch s {
	case Grounded:
		return "Grounded"
	case OnLadder:
		return "Ladder"
	case Falling:
		return "Falling"
	case Jumping:
		return "Jumping"
	case Leaping:
		return "Leaping"
	case Flying:
		return "Flying"
	case DoingAction:
		return "DoingAction"
	case Attack:
		return "Attack"
	case Hit:
		return "Hit"
	case Dying:
		return "Dying"
	case Dead:
		return "Dead"
	}
	return "Unknown"
}

type Flags struct {
	LeftWall  bool
	RightWall bool
	Ceiling   bool
	Floor     bool
	EnemyL    bool
	EnemyR    bool
	EnemyU    bool
	EnemyD    bool
	NoLadders bool
	GoingUp   bool
	Climbed   bool
	LeapOn    bool
	LeapOff   bool
	LeapTo    bool
	Breath    bool
	HighJump  bool
	LongJump  bool
	JumpR     bool
	JumpL     bool
	Action    bool
	Thrown    bool
	//Drop      bool
	//HoldSwitch bool
	HeldFlip   bool
	HeldNFlip  bool
	Hit        bool
	Crush      bool
	Attack     bool
	Regen      bool
	Flying     bool
	Frame      bool
	JumpBuff   int
	PickUpBuff int
	ActionBuff int
	ItemAction ItemAction
}

type CharacterOptions struct {
	Regen       bool
	RegenFlip   bool
	Flying      bool
	LinkedTiles []world.Coords
}

type Controller interface {
	GetActions() Actions
	ClearPrev()
	GetEntity() *ecs.Entity
}

type MoveType int

const (
	Humanoid = iota
	Flyer
)

func PlayerVars() Vars {
	return Vars{
		WalkSpeed:    constants.PlayerWalkSpeed - constants.SpeedMod,
		LeapSpeed:    constants.PlayerLeapSpeed - constants.SpeedMod,
		ClimbSpeed:   constants.PlayerClimbSpeed - constants.SpeedMod,
		SlideSpeed:   constants.PlayerSlideSpeed - constants.SpeedMod,
		Gravity:      constants.PlayerGravity - constants.SpeedMod,
		HiJumpVSpeed: constants.PlayerHighJumpSpeed - constants.SpeedMod,
		HiJumpHSpeed: constants.PlayerHighJumpHSpeed - constants.SpeedMod,
		HiJumpCntr:   constants.PlayerHighJumpCounter,
		LgJumpVSpeed: constants.PlayerLongJumpVSpeed - constants.SpeedMod,
		LgJumpHSpeed: constants.PlayerLongJumpHSpeed - constants.SpeedMod,
		LgJumpCntr:   constants.PlayerLongJumpCounter,
		IdleFreq:     constants.IdleFrequency,
	}
}

func DemonVars() Vars {
	return Vars{
		WalkSpeed:    constants.DemonWalkSpeed - constants.SpeedMod,
		LeapSpeed:    constants.DemonLeapSpeed - constants.SpeedMod,
		ClimbSpeed:   constants.DemonClimbSpeed - constants.SpeedMod,
		SlideSpeed:   constants.DemonSlideSpeed - constants.SpeedMod,
		Gravity:      constants.DemonGravity - constants.SpeedMod,
		HiJumpVSpeed: constants.DemonHighJumpSpeed - constants.SpeedMod,
		HiJumpHSpeed: constants.DemonHighJumpHSpeed - constants.SpeedMod,
		HiJumpCntr:   constants.DemonHighJumpCounter,
		LgJumpVSpeed: constants.DemonLongJumpVSpeed - constants.SpeedMod,
		LgJumpHSpeed: constants.DemonLongJumpHSpeed - constants.SpeedMod,
		LgJumpCntr:   constants.DemonLongJumpCounter,
		IdleFreq:     constants.IdleFrequency,
	}
}

func FlyVars() Vars {
	return Vars{
		WalkSpeed: constants.FlySpeed,
	}
}
