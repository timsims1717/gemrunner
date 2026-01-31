package ui

import (
	"bytes"
	"encoding/json"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/viewport"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

type ElementType int

const (
	ButtonElement = iota
	CheckboxElement
	ContainerElement
	DropdownElement
	InputElement
	MultiLineInputElement
	ScrollElement
	SliderElement
	SpriteElement
	TextElement
	CustomElement
)

var elementTypeStrings = map[ElementType]string{
	ButtonElement:         "button",
	CheckboxElement:       "checkbox",
	ContainerElement:      "container",
	DropdownElement:       "dropdown",
	InputElement:          "input",
	MultiLineInputElement: "multiline",
	ScrollElement:         "scroll",
	SliderElement:         "slider",
	SpriteElement:         "sprite",
	TextElement:           "text",
	CustomElement:         "custom",
}

var elementTypeIDs = map[string]ElementType{
	"button":    ButtonElement,
	"checkbox":  CheckboxElement,
	"container": ContainerElement,
	"dropdown":  DropdownElement,
	"input":     InputElement,
	"multiline": MultiLineInputElement,
	"scroll":    ScrollElement,
	"slider":    SliderElement,
	"sprite":    SpriteElement,
	"text":      TextElement,
	"custom":    CustomElement,
}

func (et ElementType) String() string {
	return elementTypeStrings[et]
}

func (et ElementType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(et.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (et *ElementType) UnmarshalJSON(bts []byte) error {
	var j string
	err := json.Unmarshal(bts, &j)
	if err != nil {
		return err
	}
	*et = elementTypeIDs[j]
	return nil
}

type ElementConstructor struct {
	Key         string               `json:"key"`
	Width       float64              `json:"width,omitempty"`
	Height      float64              `json:"height,omitempty"`
	SprKey      string               `json:"spriteKey,omitempty"`
	SprKey2     string               `json:"spriteKey2,omitempty"`
	Batch       string               `json:"batchKey,omitempty"`
	Text        string               `json:"text,omitempty"`
	HelpText    string               `json:"helpText,omitempty"`
	Color       pixel.RGBA           `json:"color,omitempty"`
	Background  pixel.RGBA           `json:"background,omitempty"`
	Position    pixel.Vec            `json:"pos"`
	CanFocus    bool                 `json:"canFocus,omitempty"`
	Left        string               `json:"left,omitempty"`
	Right       string               `json:"right,omitempty"`
	Up          string               `json:"up,omitempty"`
	Down        string               `json:"down,omitempty"`
	ElementType ElementType          `json:"type"`
	SubElements []ElementConstructor `json:"elements,omitempty"`
	Anchor      pixel.Anchor         `json:"anchor,omitempty"`
	Min         int                  `json:"min"`
	Max         int                  `json:"max"`
	Interval    int                  `json:"interval"`
}

type Element struct {
	Key      string
	Sprite   *img.Sprite
	Sprite2  *img.Sprite
	SpriteP  *pixel.Sprite
	Delay    float64
	HelpText string
	Object   *object.Object
	Entity   *ecs.Entity
	Action   func()
	OnClick  func()
	OnHold   func()
	OnHover  func(bool)
	OnFocus  func(bool)
	Left     string
	Right    string
	Up       string
	Down     string
	Disabled bool

	ElementType ElementType

	Checked  bool
	Value    string
	IntValue int
	Min      int
	Max      int
	Interval int
	//FloatValue float64
	InFocus    bool
	Text       *typeface.Text
	TextEntity *ecs.Entity
	CaretIndex int
	CaretObj   *object.Object
	InputType  ValueType
	MultiLine  bool

	Border     *data.Border
	ViewPort   *viewport.ViewPort
	Layer      int
	Elements   []*Element
	Focused    string
	Background pixel.RGBA

	Bar          *Element
	ScrollUp     *Element
	ScrollDown   *Element
	ButtonHeight float64
	Top          float64
	Bot          float64
}

type ValueType int

const (
	AlphaNumeric = iota
	Numeric
	Special
	Any
)

func (e *Element) Get(key string) *Element {
	for _, e1 := range e.Elements {
		if e1.Key == key {
			return e1
		}
	}
	for _, e2 := range e.Elements {
		if se2 := e2.Get(key); se2 != nil {
			return se2
		}
	}
	return nil
}
