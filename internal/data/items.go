package data

import (
	"bytes"
	"encoding/json"
	"gemrunner/internal/constants"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
)

type Interact struct {
	Fn func(int, *Dynamic, *ecs.Entity) bool
}

func NewInteract(fn func(int, *Dynamic, *ecs.Entity) bool) *Interact {
	return &Interact{Fn: fn}
}

type PickUp struct {
	Cycle     [constants.MaxPlayers]int
	Priority  int
	Inventory int
	Color     ItemColor
}

func NewPickUp(p int, color ItemColor) *PickUp {
	return &PickUp{
		Priority:  p,
		Inventory: -1,
		Color:     color,
	}
}

type BasicItem struct {
	Name     string           `json:"name"`
	Key      string           `json:"key"`
	Object   *object.Object   `json:"-"`
	Entity   *ecs.Entity      `json:"-"`
	Sprite   *img.Sprite      `json:"-"`
	Anim     *reanimator.Tree `json:"-"`
	PickUp   *PickUp          `json:"-"`
	Action   *Interact        `json:"-"`
	Color    ItemColor        `json:"color"`
	Metadata TileMetadata     `json:"metadata"`
	Origin   world.Coords     `json:"origin"`
	Block    Block            `json:"block"`
	Layer    int              `json:"layer"`
	Using    bool             `json:"-"`
	Regen    bool             `json:"-"`
	Waiting  bool             `json:"-"`
	Counter  int              `json:"-"`
	Delay    int              `json:"-"`
	Uses     int              `json:"uses"`
}

type Door struct {
	Item     *BasicItem
	DoorType DoorType
	Unlock   bool
}

type Transporter struct {
	Item  *BasicItem
	Dest  *Tile
	BarE  *ecs.Entity
	BarO  *object.Object
	BarT  *reanimator.Tree
	BarUp bool
}

type ItemColor int

const (
	ColorDefault = iota
	NonPlayerYellow
	NonPlayerBrown
	NonPlayerGray
	NonPlayerCyan
	NonPlayerLime
	NonPlayerPink
	NonPlayerBurnt
	NonPlayerRed
	PlayerBlue
	PlayerGreen
	PlayerPurple
	PlayerOrange
)

func (ic ItemColor) String() string {
	switch ic {
	case ColorDefault:
		return "default"
	case NonPlayerYellow:
		return "yellow"
	case NonPlayerBrown:
		return "brown"
	case NonPlayerGray:
		return "gray"
	case NonPlayerCyan:
		return "cyan"
	case NonPlayerLime:
		return "lime"
	case NonPlayerPink:
		return "pink"
	case NonPlayerBurnt:
		return "burnt"
	case PlayerBlue:
		return "blue"
	case PlayerGreen:
		return "green"
	case PlayerPurple:
		return "purple"
	case PlayerOrange:
		return "orange"
	case NonPlayerRed:
		return "red"
	}
	return ""
}

func (ic ItemColor) SpriteString() string {
	switch ic {
	case NonPlayerYellow:
		return "_yellow"
	case NonPlayerBrown:
		return "_brown"
	case NonPlayerGray:
		return "_gray"
	case NonPlayerCyan:
		return "_cyan"
	case NonPlayerLime:
		return "_lime"
	case NonPlayerPink:
		return "_pink"
	case NonPlayerBurnt:
		return "_burnt"
	case PlayerBlue:
		return "_blue"
	case PlayerGreen:
		return "_green"
	case PlayerPurple:
		return "_purple"
	case PlayerOrange:
		return "_orange"
	case NonPlayerRed:
		return "_red"
	default:
		return "_yellow"
	}
}

var colorToID = map[string]ItemColor{
	"default": ColorDefault,
	"yellow":  NonPlayerYellow,
	"brown":   NonPlayerBrown,
	"gray":    NonPlayerGray,
	"cyan":    NonPlayerCyan,
	"lime":    NonPlayerLime,
	"pink":    NonPlayerPink,
	"burnt":   NonPlayerBurnt,
	"red":     NonPlayerRed,
	"blue":    PlayerBlue,
	"green":   PlayerGreen,
	"purple":  PlayerPurple,
	"orange":  PlayerOrange,
}

func (ic ItemColor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(ic.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (ic *ItemColor) UnmarshalJSON(bts []byte) error {
	var j string
	err := json.Unmarshal(bts, &j)
	if err != nil {
		var ji int
		err = json.Unmarshal(bts, &ji)
		if err != nil {
			return err
		}
		*ic = ItemColor(ji)
		return nil
	}
	*ic = colorToID[j]
	return nil
}

type DoorType int

const (
	Hidden = iota
	Visible
	Locked
	Unlocked
)

type Bomb struct {
	Item   *BasicItem
	Draws  []interface{}
	SymSpr *img.Sprite
	LitKey string
	Prefix string
	Big    bool
}

type Disguise struct {
	Item *BasicItem
	Doff bool
}

type Snare struct {
	Item *BasicItem
	Ch   *Dynamic
}
