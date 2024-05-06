package data

import "github.com/gopxl/pixel"

var (
	PuzzleSetFileList      []PuzzleSetMetadata
	SelectedPuzzleIndex    int
	CustomWorldSelected    bool
	CustomSelectedBefore   bool
	SelectedWorldIndex     int
	SelectedPrimaryColor   pixel.RGBA
	SelectedSecondaryColor pixel.RGBA
	SelectedDoodadColor    pixel.RGBA
)

type PuzzleMetadata struct {
	Name           string     `json:"title"`
	Filename       string     `json:"filename"`
	WorldSprite    string     `json:"sprite"`
	WorldNumber    int        `json:"world"`
	PrimaryColor   pixel.RGBA `json:"primaryColor"`
	SecondaryColor pixel.RGBA `json:"secondaryColor"`
	DoodadColor    pixel.RGBA `json:"doodadColor"`
	MusicTrack     string     `json:"musicTrack"`
	Completed      bool       `json:"completed"`
}

type PuzzleSetMetadata struct {
	Name     string `json:"title"`
	Filename string `json:"filename"`
	Author   string `json:"author"`
}
