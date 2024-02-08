package data

import "github.com/gopxl/pixel"

var (
	PuzzleInfos         []PuzzleMetadata
	SelectedPuzzleIndex int
)

type PuzzleMetadata struct {
	Name           string     `json:"title"`
	Filename       string     `json:"-"`
	WorldSprite    string     `json:"sprite"`
	WorldNumber    int        `json:"world"`
	PrimaryColor   pixel.RGBA `json:"primaryColor"`
	SecondaryColor pixel.RGBA `json:"secondaryColor"`
	Completed      bool       `json:"completed"`
}
