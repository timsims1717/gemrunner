package data

var (
	PuzzleInfos         []PuzzleInfo
	SelectedPuzzleIndex int
)

type PuzzleInfo struct {
	Name     string `json:"title"`
	Filename string `json:"-"`
}
