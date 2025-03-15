package data

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"time"
)

var (
	PuzzleView         *viewport.ViewPort
	PuzzleViewNoShader *viewport.ViewPort
	WorldView          *viewport.ViewPort
	BorderView         *viewport.ViewPort
	IMDraw             *imdraw.IMDraw

	CurrLevelSess *LevelSession
	CurrLevel     *Level
	CurrPuzzleSet *PuzzleSet
	CurrReplay    *LevelReplay

	CurrSelect *Selection
	ClipSelect *Selection
	EditorDraw bool

	ColorShader  string
	PuzzleShader string
	WorldShader  string
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
	PLoc      [constants.MaxPlayers]*mgl32.Vec2
	Start     bool
	Failed    bool
	Complete  bool
	ExitIndex int
	DoorsOpen bool

	FakePlayer        *Dynamic
	FakePlayerDir     Direction
	FakePlayerCounter int

	Recording   bool
	SaveRecord  bool
	LevelReplay *LevelReplay
	ReplayFrame ReplayFrame

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

	Elapsed float32 `json:"-"`

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
	T [][]*Tile
}

func (l *Level) Get(x, y int) *Tile {
	if x < 0 || y < 0 || x >= l.Metadata.Width || y >= l.Metadata.Height {
		return nil
	}
	return l.Tiles.T[y][x]
}

func (p *Puzzle) Get(x, y int) *Tile {
	if x < 0 || y < 0 || x >= p.Metadata.Width || y >= p.Metadata.Height {
		return nil
	}
	return p.Tiles.T[y][x]
}

func (p *Puzzle) CoordsLegal(coords world.Coords) bool {
	return coords.X >= 0 && coords.Y >= 0 && coords.X < p.Metadata.Width && coords.Y < p.Metadata.Height
}

func (p *Puzzle) GetClosestLegal(coords world.Coords) world.Coords {
	if p.CoordsLegal(coords) {
		return coords
	}
	if coords.X < 0 {
		coords.X = 0
	} else if coords.X >= p.Metadata.Width {
		coords.X = p.Metadata.Width - 1
	}
	if coords.Y < 0 {
		coords.Y = 0
	} else if coords.Y >= p.Metadata.Height {
		coords.Y = p.Metadata.Height - 1
	}
	return coords
}

func (p *Puzzle) SetWidth(w int) {
	if w < 6 {
		w = 6
	} else if w > constants.PuzzleMaxWidth {
		w = constants.PuzzleMaxWidth
	}
	if p.Metadata.Width == w {
		return
	} else if p.Metadata.Width < w { // add width
		for y, row := range p.Tiles.T {
			for x := len(row); x < w; x++ {
				tile := &Tile{}
				tile.Coords = world.NewCoords(x, y)
				obj := object.New()
				obj.Pos = world.MapToWorld(tile.Coords)
				obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
				obj.Layer = 2
				e := myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
				tile.Object = obj
				tile.Entity = e
				tile.ToEmpty()
				tile.Update = true
				row = append(row, tile)
			}
		}
	} else { // remove width
		for _, row := range p.Tiles.T {
			for x := p.Metadata.Width - 1; x >= w; x-- {
				myecs.Manager.DisposeEntity(row[x].Entity)
				row = row[:x]
			}
		}
	}
	p.Metadata.Width = w
	p.Update = true
}

func (p *Puzzle) SetHeight(h int) {
	if h < 6 {
		h = 6
	} else if h > constants.PuzzleMaxHeight {
		h = constants.PuzzleMaxHeight
	}
	if p.Metadata.Height == h {
		return
	} else if p.Metadata.Height < h { // add height
		for y := p.Metadata.Height; y < h; y++ {
			p.Tiles.T = append(p.Tiles.T, []*Tile{})
			for x := 0; x < p.Metadata.Width; x++ {
				tile := &Tile{}
				tile.Coords = world.NewCoords(x, y)
				obj := object.New()
				obj.Pos = world.MapToWorld(tile.Coords)
				obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
				obj.Layer = 2
				e := myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
				tile.Object = obj
				tile.Entity = e
				tile.ToEmpty()
				tile.Update = true
				p.Tiles.T[y] = append(p.Tiles.T[y], tile)
			}
		}
	} else { // remove height
		for y := p.Metadata.Height - 1; y >= h; y-- {
			for _, tile := range p.Tiles.T[y] {
				myecs.Manager.DisposeEntity(tile.Entity)
			}
		}
		p.Tiles.T = p.Tiles.T[:h]
	}
	p.Metadata.Height = h
	p.Update = true
}

func NewTiles(w, h int) *Tiles {
	var t [][]*Tile
	for y := 0; y < h; y++ {
		t = append(t, []*Tile{})
		for x := 0; x < w; x++ {
			tile := &Tile{}
			tile.ToEmpty()
			t[y] = append(t[y], tile)
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
		Width:          constants.PuzzleWidth,
		Height:         constants.PuzzleHeight,
		WorldSprite:    constants.WorldSprites[worldNum],
		WorldNumber:    worldNum,
		PrimaryColor:   pixel.ToRGBA(constants.WorldPrimary[worldNum]),
		SecondaryColor: pixel.ToRGBA(constants.WorldSecondary[worldNum]),
		DoodadColor:    pixel.ToRGBA(constants.WorldDoodad[worldNum]),
		MusicTrack:     constants.WorldMusic[worldNum],
	}
	puz := &Puzzle{
		Tiles:    NewTiles(md.Width, md.Height),
		Metadata: md,
	}
	for y := 0; y < md.Height; y++ {
		for x := 0; x < md.Width; x++ {
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
	tiles := NewTiles(p.Metadata.Width, p.Metadata.Height)
	for y := 0; y < p.Metadata.Height; y++ {
		for x := 0; x < p.Metadata.Width; x++ {
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
