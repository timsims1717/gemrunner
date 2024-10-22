package ui

import (
	"bytes"
	"encoding/json"
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
	InputElement
	MultiLineInputElement
	ScrollElement
	SpriteElement
	TextElement
	CustomElement
)

var elementTypeStrings = map[ElementType]string{
	ButtonElement:         "button",
	CheckboxElement:       "checkbox",
	ContainerElement:      "container",
	InputElement:          "input",
	MultiLineInputElement: "multiline",
	ScrollElement:         "scroll",
	SpriteElement:         "sprite",
	TextElement:           "text",
	CustomElement:         "custom",
}

var elementTypeIDs = map[string]ElementType{
	"button":    ButtonElement,
	"checkbox":  CheckboxElement,
	"container": ContainerElement,
	"input":     InputElement,
	"multiline": MultiLineInputElement,
	"scroll":    ScrollElement,
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
	Position    pixel.Vec            `json:"pos"`
	CanFocus    bool                 `json:"canFocus,omitempty"`
	Left        string               `json:"leftKey,omitempty"`
	Right       string               `json:"rightKey,omitempty"`
	Up          string               `json:"upKey,omitempty"`
	Down        string               `json:"downKey,omitempty"`
	ElementType ElementType          `json:"type"`
	SubElements []ElementConstructor `json:"elements,omitempty"`
	Anchor      pixel.Anchor         `json:"anchor,omitempty"`
}

type Element struct {
	Key      string
	Sprite   *img.Sprite
	Sprite2  *img.Sprite
	Delay    float64
	HelpText string
	Object   *object.Object
	Entity   *ecs.Entity
	Action   func()
	OnClick  func()
	OnHold   func()
	OnHover  func(bool)
	Left     string
	Right    string
	Up       string
	Down     string

	ElementType ElementType

	Checked    bool
	Value      string
	Focused    bool
	Text       *typeface.Text
	CaretIndex int
	CaretObj   *object.Object
	InputType  InputType
	MultiLine  bool

	Border       *Border
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	ViewPort     *viewport.ViewPort
	Layer        int
	Elements     []*Element

	Bar          *Element
	ScrollUp     *Element
	ScrollDown   *Element
	ButtonHeight float64
	YTop         float64
	YBot         float64
}

type InputType int

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
	return nil
}
