package data

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data/death"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/google/uuid"
)

type Dynamic struct {
	Object       *object.Object
	Anims        *reanimator.TreeSet
	Entity       *ecs.Entity
	Inventory    *BasicItem
	StoredBlocks []*Tile
	Control      Controller

	Actions  Actions
	State    CharacterState
	Flags    Flags
	Vars     Vars
	Options  CharacterOptions
	ACounter int
	AnInt    int
	LastTile *Tile
	NextTile *Tile
	MoveType MoveType
	Player   int
	Enemy    int
	Color    ItemColor
	Layer    int

	SFX *uuid.UUID
}

func NewDynamic(tile *Tile) *Dynamic {
	return &Dynamic{
		Anims:   reanimator.NewSet(),
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
		AnInt:    -1,
		LastTile: tile,
	}
}

type Vars struct {
	WalkSpeed    float64
	BarSpeed     float64
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
	OnBar
	Falling
	Jumping
	Leaping
	Flying
	InHiding
	DoingAction
	Attack
	Hit
	Dying
	Dead
	Waiting
	Regen
)

type ItemAction int

const (
	NoItemAction = iota
	MagicDig
	MagicPlace
	ThrowBox
	DonDisguise
	DrillStart
	Drilling
	Hiding
	FireFlamethrower
	TransportIn
	TransportExit
)

func (s CharacterState) String() string {
	switch s {
	case Grounded:
		return "Grounded"
	case OnLadder:
		return "OnLadder"
	case OnBar:
		return "OnBar"
	case Falling:
		return "Falling"
	case Jumping:
		return "Jumping"
	case Leaping:
		return "Leaping"
	case Flying:
		return "Flying"
	case InHiding:
		return "InHiding"
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
	case Waiting:
		return "Waiting"
	}
	return "Unknown"
}

type Flags struct {
	LeftWall     bool
	RightWall    bool
	Ceiling      bool
	Floor        bool
	EnemyL       bool
	EnemyR       bool
	EnemyU       bool
	EnemyD       bool
	NoLadders    bool
	GoingUp      bool
	Climbed      bool
	LeapOn       bool
	LeapOff      bool
	LeapTo       bool
	CanActLeap   bool
	Breath       bool
	HighJump     bool
	LongJump     bool
	Landing      bool
	JumpR        bool
	JumpL        bool
	Thrown       bool
	Death        death.Type
	Attack       bool
	Regen        bool
	Transport    bool
	Flying       bool
	Disguised    bool
	CheckAction  bool
	Frame        bool
	PickUpBuff   int
	ActionBuff   int
	DigLeftBuff  int
	DigRightBuff int
	ItemAction   ItemAction
}

type CharacterOptions struct {
	Regen       bool
	RegenFlip   bool
	Flying      bool
	LinkedTiles []world.Coords
	StoredCount int
}

type MoveType int

const (
	Humanoid = iota
	Flyer
)

func PlayerVars() Vars {
	return Vars{
		WalkSpeed:    constants.PlayerWalkSpeed - constants.SpeedMod,
		BarSpeed:     constants.PlayerBarSpeed - constants.SpeedMod,
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
		BarSpeed:     constants.DemonBarSpeed - constants.SpeedMod,
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
