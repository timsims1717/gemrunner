package data

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/random"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"time"
)

var (
	PuzzleView         *viewport.ViewPort
	PuzzleViewNoShader *viewport.ViewPort
	BorderView         *viewport.ViewPort
	IMDraw             *imdraw.IMDraw

	CurrLevelSess *LevelSession
	CurrLevel     *Level
	CurrPuzzleSet *PuzzleSet
	CurrSelect    *Selection
	ClipSelect    *Selection
	EditorDraw    bool

	PuzzleShader string
)

type LevelSession struct {
	PlayerStats [constants.MaxPlayers]*PlayerStats `json:"playerStats"`
	Levels      LevelCompletion                    `json:"levels"`

	LevelStart   time.Time     `json:"-"`
	TimePlayed   time.Duration `json:"-"`
	TotalTime    time.Duration `json:"totalTime"`
	PuzzleIndex  int           `json:"puzzleIndex"`
	PuzzleFile   string        `json:"puzzleFile"`
	Filename     string        `json:"filename"`
	TotalGems    int           `json:"totalGems"`
	TotalRedGems int           `json:"totalRedGems"`

	// who has a red key
	// abilities

	PuzzleSet *PuzzleSet        `json:"-"`
	Metadata  PuzzleSetMetadata `json:"-"`
}

type LevelCompletion struct {
	Index     int            `json:"index"`
	GemScore  int            `json:"gemScore"`
	Completed bool           `json:"completed"`
	Changed   []world.Coords `json:"changed"`
}

type Level struct {
	Tiles     *Tiles
	Enemies   []*Dynamic
	Players   [constants.MaxPlayers]*Dynamic
	PControls [constants.MaxPlayers]Controller
	Start     bool
	Failed    bool
	Complete  bool
	DoorsOpen bool

	FakePlayer        *Dynamic
	FakePlayerDir     Direction
	FakePlayerCounter int

	FrameNumber  int
	FrameCounter int
	FrameCycle   int
	FrameChange  bool

	Puzzle   *Puzzle
	Metadata PuzzleMetadata
}

type PuzzleSet struct {
	Puzzles  []*Puzzle         `json:"puzzles"`
	Metadata PuzzleSetMetadata `json:"puzzleSetMetadata"`

	CurrPuzzle *Puzzle `json:"-"`

	PuzzleIndex int  `json:"-"`
	NeedToSave  bool `json:"-"`
}

type Puzzle struct {
	Tiles *Tiles `json:"tiles"`

	WrenchTiles []*Tile `json:"-"`

	UndoStack  []*Tiles `json:"-"`
	LastChange *Tiles   `json:"-"`
	RedoStack  []*Tiles `json:"-"`

	Click   bool `json:"-"`
	Update  bool `json:"-"`
	Changed bool `json:"-"`

	Metadata PuzzleMetadata `json:"metadata"`
}

type Tiles struct {
	T [constants.PuzzleHeight][constants.PuzzleWidth]*Tile
}

func (t *Tiles) Get(x, y int) *Tile {
	if x < 0 || y < 0 || x >= constants.PuzzleWidth || y >= constants.PuzzleHeight {
		return nil
	}
	return t.T[y][x]
}

func NewTiles() *Tiles {
	t := [constants.PuzzleHeight][constants.PuzzleWidth]*Tile{}
	for _, row := range t {
		for _, tile := range row {
			tile = &Tile{}
			tile.ToEmpty()
		}
	}
	return &Tiles{
		T: t,
	}
}

type Selection struct {
	Tiles  [][]*Tile
	Offset world.Coords
	Origin world.Coords
	Width  int
	Height int
}

func CreatePuzzleSet() *PuzzleSet {
	pzSet := &PuzzleSet{}
	pzSet.SetToFirst()
	return pzSet
}

func (set *PuzzleSet) AppendNew() {
	pzl := CreateBlankPuzzle()
	set.Puzzles = append(set.Puzzles, pzl)
	set.CurrPuzzle = pzl
	set.PuzzleIndex = len(set.Puzzles) - 1
}

func (set *PuzzleSet) Insert(pzl *Puzzle, i int) {
	if i > -1 && len(set.Puzzles) > i {
		if pzl == nil {
			pzl = CreateBlankPuzzle()
		}
		set.Puzzles = append(append(set.Puzzles[:i+1], set.Puzzles[i:]...))
		set.Puzzles[i] = pzl
		set.CurrPuzzle = pzl
		set.PuzzleIndex = i
	} else {
		set.Append(pzl)
	}
}

func (set *PuzzleSet) Delete(i int) {
	if i > -1 && len(set.Puzzles) > i {
		if len(set.Puzzles) == 1 {
			pzl := CreateBlankPuzzle()
			set.Puzzles = []*Puzzle{pzl}
			set.CurrPuzzle = pzl
			set.PuzzleIndex = 0
		} else {
			set.Puzzles = append(set.Puzzles[:i], set.Puzzles[i+1:]...)
			if set.PuzzleIndex == i {
				set.PuzzleIndex = i - 1
				if set.PuzzleIndex < 0 {
					set.PuzzleIndex = 0
				}
				set.CurrPuzzle = set.Puzzles[set.PuzzleIndex]
			}
		}
	}
}

func (set *PuzzleSet) Append(pzl *Puzzle) {
	if pzl == nil {
		set.AppendNew()
		return
	}
	set.Puzzles = append(set.Puzzles, pzl)
	set.CurrPuzzle = pzl
	set.PuzzleIndex = len(set.Puzzles) - 1
}

func (set *PuzzleSet) Replace(i int, pzl *Puzzle) {
	if i > -1 && len(set.Puzzles) > i {
		set.Puzzles[i] = pzl
		set.CurrPuzzle = pzl
		set.PuzzleIndex = i
	}
}

func (set *PuzzleSet) SetToFirst() {
	if len(set.Puzzles) == 0 {
		set.Puzzles = append(set.Puzzles, CreateBlankPuzzle())
	}
	set.CurrPuzzle = set.Puzzles[0]
	set.PuzzleIndex = 0
}

func (set *PuzzleSet) SetTo(i int) {
	if i > -1 && len(set.Puzzles) > i {
		set.CurrPuzzle = set.Puzzles[i]
		set.PuzzleIndex = i
	}
}

func (set *PuzzleSet) Next() {
	set.PuzzleIndex++
	if set.PuzzleIndex == len(set.Puzzles) {
		set.PuzzleIndex = 0
	}
	set.CurrPuzzle = set.Puzzles[set.PuzzleIndex]
}

func (set *PuzzleSet) Prev() {
	set.PuzzleIndex--
	if set.PuzzleIndex == -1 {
		set.PuzzleIndex = len(set.Puzzles) - 1
	}
	set.CurrPuzzle = set.Puzzles[set.PuzzleIndex]
}

func CreateBlankPuzzle() *Puzzle {
	SelectedWorldIndex = 0
	worldNum := constants.WorldMoss
	md := PuzzleMetadata{
		WorldSprite:    constants.WorldSprites[worldNum],
		WorldNumber:    worldNum,
		PrimaryColor:   pixel.ToRGBA(constants.WorldPrimary[worldNum]),
		SecondaryColor: pixel.ToRGBA(constants.WorldSecondary[worldNum]),
		DoodadColor:    pixel.ToRGBA(constants.WorldDoodad[worldNum]),
		MusicTrack:     constants.WorldMusic[worldNum],
	}
	puz := &Puzzle{
		Tiles:    NewTiles(),
		Metadata: md,
	}
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			c := world.Coords{X: x, Y: y}
			alt := 0
			if random.Effects.Intn(80) == 0 {
				alt = 1
			}
			puz.Tiles.T[y][x] = &Tile{
				Block:    Block(BlockEmpty),
				Coords:   c,
				AltBlock: alt,
			}
		}
	}
	return puz
}

func (p *Puzzle) CopyTiles() *Tiles {
	tiles := NewTiles()
	for y := 0; y < constants.PuzzleHeight; y++ {
		for x := 0; x < constants.PuzzleWidth; x++ {
			tiles.T[y][x] = p.Tiles.T[y][x].Copy()
		}
	}
	return tiles
}

func (p *Puzzle) HasPlayers() bool {
	hasPlayers := false
playerCheck:
	for _, row := range p.Tiles.T {
		for _, t := range row {
			if t.Block == BlockPlayer1 ||
				t.Block == BlockPlayer2 ||
				t.Block == BlockPlayer3 ||
				t.Block == BlockPlayer4 {
				hasPlayers = true
				break playerCheck
			}
		}
	}
	return hasPlayers
}

func (p *Puzzle) NumPlayers() int {
	numPlayers := 0
	for _, row := range p.Tiles.T {
		for _, t := range row {
			if t.Block == BlockPlayer1 ||
				t.Block == BlockPlayer2 ||
				t.Block == BlockPlayer3 ||
				t.Block == BlockPlayer4 {
				numPlayers++
			}
		}
	}
	return numPlayers
}
