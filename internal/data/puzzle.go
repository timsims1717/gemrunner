package data

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/internal/random"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/object"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"time"
)

var (
	CurrentPlayArea *PlayArea
	AllPlayAreas    []*PlayArea

	ScreenView  *viewport.ViewPort
	ScreenShake *util.NoiseShaker

	CurrLevelSess *LevelSession
	CurrLevel     *Level
	CurrPuzzleSet *PuzzleSet
	CurrReplay    *LevelReplay
	LevelTrans    bool

	CurrSelect *Selection
	ClipSelect *Selection
	EditorDraw bool

	ColorShader  string
	PuzzleShader string
	WorldShader  string
	ScreenShader string
	ShaderTime   float32
)

type PlayArea struct {
	PuzzleView         *viewport.ViewPort
	PuzzleViewNoShader *viewport.ViewPort
	WorldView          *viewport.ViewPort
	BorderView         *viewport.ViewPort
	IMDraw             *imdraw.IMDraw

	Level  *Level
	Puzzle *Puzzle

	BorderEntity *ecs.Entity
	Border       *Border
	BorderObject *object.Object
}

type LevelSession struct {
	PlayerStats [constants.MaxPlayers]*PlayerStats `json:"playerStats"`

	LevelArray    []LevelCompletion       `json:"levels"`
	LevelMap      map[int]LevelCompletion `json:"-"`
	GemsCollected []world.Coords          `json:"-"`

	LevelStart   time.Time     `json:"-"`
	TimePlayed   time.Duration `json:"-"`
	TotalTime    time.Duration `json:"totalTime"`
	PuzzleIndex  int           `json:"puzzleIndex"`
	LastPuzzle   int           `json:"lastPuzzle"`
	StartCoords  *world.Coords `json:"startCoords"`
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
	Index         int            `json:"index"`
	GemsCollected []world.Coords `json:"gemsCollected"`
	Completed     bool           `json:"completed"`
	Continuity    bool           `json:"continuity"`
}

type Level struct {
	Tiles       *Tiles
	Enemies     []*Dynamic
	AllEntities []*ecs.Entity
	Players     [constants.MaxPlayers]*Dynamic
	PControls   [constants.MaxPlayers]Controller
	PLoc        [constants.MaxPlayers]*mgl32.Vec2
	Start       bool
	Failed      bool
	Complete    bool
	ExitIndex   int
	StartCoords *world.Coords
	DoorsOpen   bool

	FakePlayer        *Dynamic
	FakePlayerDir     Direction
	FakePlayerCounter int

	Continuity  bool
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

type LevelTransition struct {
	Open      bool
	ExitIndex int
	ExitTile  world.Coords
}

type PuzzleSet struct {
	Puzzles  []*Puzzle         `json:"puzzles"`
	Metadata PuzzleSetMetadata `json:"puzzleSetMetadata"`
	HasGrid  bool              `json:"hasGrid"`

	CurrPuzzle *Puzzle              `json:"-"`
	PuzzGrid   map[world.Coords]int `json:"-"`
	GridMin    world.Coords         `json:"-"`
	GridMax    world.Coords         `json:"-"`

	PuzzleIndex int  `json:"-"`
	NeedToSave  bool `json:"-"`
}

type Puzzle struct {
	Tiles *Tiles       `json:"tiles"`
	Grid  world.Coords `json:"grid"`

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
	pzSet.SetUpGrid()
	pzSet.NeedToSave = false
	return pzSet
}

func (set *PuzzleSet) AppendNew() {
	pzl := CreateBlankPuzzle()
	grid := set.GetAvailableGridCoords()
	pzl.Grid = grid
	set.Puzzles = append(set.Puzzles, pzl)
	set.CurrPuzzle = pzl
	set.PuzzleIndex = len(set.Puzzles) - 1
	set.PuzzGrid[grid] = set.PuzzleIndex
	set.UpdateGridMaxMin()
	set.NeedToSave = true
}

func (set *PuzzleSet) Insert(pzl *Puzzle, i int) {
	if i > -1 && len(set.Puzzles) > i {
		if pzl == nil {
			pzl = CreateBlankPuzzle()
		}
		grid := set.GetAvailableGridCoords()
		pzl.Grid = grid
		set.PuzzGrid[grid] = i
		set.Puzzles = append(append(set.Puzzles[:i+1], set.Puzzles[i:]...))
		set.Puzzles[i] = pzl
		set.CurrPuzzle = pzl
		set.PuzzleIndex = i
	} else {
		set.Append(pzl)
	}
	set.UpdateGridMaxMin()
	set.NeedToSave = true
}

func (set *PuzzleSet) InsertGrid(pzl *Puzzle, grid world.Coords) {
	if pzl == nil {
		pzl = CreateBlankPuzzle()
	}
	pzl.Grid = grid
	set.Puzzles = append(set.Puzzles, pzl)
	set.CurrPuzzle = pzl
	set.PuzzleIndex = len(set.Puzzles) - 1
	set.PuzzGrid[grid] = set.PuzzleIndex
	set.UpdateGridMaxMin()
	set.NeedToSave = true
}

func (set *PuzzleSet) Delete(i int) {
	if i > -1 && len(set.Puzzles) > i {
		if len(set.Puzzles) == 1 {
			pzl := CreateBlankPuzzle()
			set.Puzzles = []*Puzzle{pzl}
			set.CurrPuzzle = pzl
			set.PuzzleIndex = 0
			set.HasGrid = false
			set.SetUpGrid()
			set.NeedToSave = true
		} else {
			pzl := set.Puzzles[i]
			if index, ok := set.PuzzGrid[pzl.Grid]; ok && index == i {
				delete(set.PuzzGrid, pzl.Grid)
			}
			set.Puzzles = append(set.Puzzles[:i], set.Puzzles[i+1:]...)
			if set.PuzzleIndex == i {
				set.PuzzleIndex = i - 1
				if set.PuzzleIndex < 0 {
					set.PuzzleIndex = 0
				}
				set.CurrPuzzle = set.Puzzles[set.PuzzleIndex]
			}
			set.SetUpGrid()
			set.NeedToSave = true
		}
	}
}

func (set *PuzzleSet) Append(pzl *Puzzle) {
	if pzl == nil {
		set.AppendNew()
		return
	}
	grid := set.GetAvailableGridCoords()
	pzl.Grid = grid
	set.Puzzles = append(set.Puzzles, pzl)
	set.CurrPuzzle = pzl
	set.PuzzleIndex = len(set.Puzzles) - 1
	set.PuzzGrid[grid] = set.PuzzleIndex
	set.UpdateGridMaxMin()
	set.NeedToSave = true
}

func (set *PuzzleSet) SetToFirst() {
	if len(set.Puzzles) == 0 {
		set.Puzzles = append(set.Puzzles, CreateBlankPuzzle())
		set.UpdateGridMaxMin()
	}
	set.CurrPuzzle = set.Puzzles[0]
	set.PuzzleIndex = 0
}

func (set *PuzzleSet) SetToCoords(c world.Coords) {
	if !set.HasGrid {
		set.SetUpGrid()
	}
	if i, ok := set.PuzzGrid[c]; ok {
		set.SetTo(i)
	}
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

func (set *PuzzleSet) Up() {
	if !set.HasGrid {
		set.SetUpGrid()
	}
	pzl := set.Puzzles[set.PuzzleIndex]
	grid := pzl.Grid
	for {
		grid.Y++
		if i, ok := set.PuzzGrid[grid]; ok {
			set.SetTo(i)
			return
		}
		if grid.Y > set.GridMax.Y {
			grid.Y = set.GridMin.Y - 1
		}
		if grid.Y == pzl.Grid.Y {
			return
		}
	}
}

func (set *PuzzleSet) Down() {
	if !set.HasGrid {
		set.SetUpGrid()
	}
	pzl := set.Puzzles[set.PuzzleIndex]
	grid := pzl.Grid
	for {
		grid.Y--
		if i, ok := set.PuzzGrid[grid]; ok {
			set.SetTo(i)
			return
		}
		if grid.Y < set.GridMin.Y {
			grid.Y = set.GridMax.Y + 1
		}
		if grid.Y == pzl.Grid.Y {
			return
		}
	}
}

func (set *PuzzleSet) Right() {
	if !set.HasGrid {
		set.SetUpGrid()
	}
	pzl := set.Puzzles[set.PuzzleIndex]
	grid := pzl.Grid
	for {
		grid.X++
		if i, ok := set.PuzzGrid[grid]; ok {
			set.SetTo(i)
			return
		}
		if grid.X > set.GridMax.X {
			grid.X = set.GridMin.X - 1
		}
		if grid.X == pzl.Grid.X {
			return
		}
	}
}

func (set *PuzzleSet) Left() {
	if !set.HasGrid {
		set.SetUpGrid()
	}
	pzl := set.Puzzles[set.PuzzleIndex]
	grid := pzl.Grid
	for {
		grid.X--
		if i, ok := set.PuzzGrid[grid]; ok {
			set.SetTo(i)
			return
		}
		if grid.X < set.GridMin.X {
			grid.X = set.GridMax.X + 1
		}
		if grid.X == pzl.Grid.X {
			return
		}
	}
}

func CreateBlankPuzzle() *Puzzle {
	SelectedWorldIndex = 0
	worldNum := constants.WorldMoss
	md := PuzzleMetadata{
		Width:                constants.PuzzleWidth,
		Height:               constants.PuzzleHeight,
		WorldSprite:          constants.WorldSprites[worldNum],
		WorldNumber:          worldNum,
		PrimaryColor:         pixel.ToRGBA(constants.WorldPrimary[worldNum]),
		SecondaryColor:       pixel.ToRGBA(constants.WorldSecondary[worldNum]),
		DoodadColor:          pixel.ToRGBA(constants.WorldDoodad[worldNum]),
		GoopColor:            pixel.ToRGBA(constants.WorldGoopColor[worldNum]),
		WorldLiquid:          constants.WorldLiquids[worldNum],
		LiquidPrimaryColor:   pixel.ToRGBA(constants.WorldLiquidPrimary[worldNum]),
		LiquidSecondaryColor: pixel.ToRGBA(constants.WorldLiquidSecondary[worldNum]),
		MusicTrack:           constants.WorldMusic[worldNum],
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

func (set *PuzzleSet) SetUpUUIDs() {
	for _, pzl := range set.Puzzles {
		if pzl.Metadata.UUID == "" {
			pzl.Metadata.UUID = uuid.New().String()
		}
	}
}

func (set *PuzzleSet) SetUpGrid() {
	if !set.HasGrid {
		if debug.Verbose {
			fmt.Printf("INFO: puzzle set grid init\n")
		}
		set.GridMin = world.Coords{}
		set.GridMax = world.Coords{}
		set.PuzzGrid = make(map[world.Coords]int)
		dx := 1
		dy := 0
		segLen := 1
		segPass := 0
		grid := world.Coords{}
		for i, pzl := range set.Puzzles {
			pzl.Grid = grid
			set.PuzzGrid[grid] = i
			// set max/min
			if grid.X < set.GridMin.X {
				set.GridMin.X = grid.X
			}
			if grid.Y < set.GridMin.Y {
				set.GridMin.Y = grid.Y
			}
			if grid.X > set.GridMax.X {
				set.GridMax.X = grid.X
			}
			if grid.Y > set.GridMax.Y {
				set.GridMax.Y = grid.Y
			}
			// make next step
			grid.X += dx
			grid.Y += dy
			segPass++
			if segPass == segLen { // done w/segment
				segPass = 0
				// rotate
				dx, dy = dy, -dx
				if dy == 0 { // increase segment length
					segLen++
				}
			}
		}
		set.NeedToSave = true
		set.HasGrid = true
	} else {
		set.PuzzGrid = make(map[world.Coords]int)
		set.GridMin = world.Coords{}
		set.GridMax = world.Coords{}
		for i, pzl := range set.Puzzles {
			set.PuzzGrid[pzl.Grid] = i
			// set max/min
			if pzl.Grid.X < set.GridMin.X {
				set.GridMin.X = pzl.Grid.X
			}
			if pzl.Grid.Y < set.GridMin.Y {
				set.GridMin.Y = pzl.Grid.Y
			}
			if pzl.Grid.X > set.GridMax.X {
				set.GridMax.X = pzl.Grid.X
			}
			if pzl.Grid.Y > set.GridMax.Y {
				set.GridMax.Y = pzl.Grid.Y
			}
		}
	}
}

func (set *PuzzleSet) GetAvailableGridCoords() world.Coords {
	if !set.HasGrid {
		set.SetUpGrid()
	}
	set.GridMin = world.Coords{}
	set.GridMax = world.Coords{}
	dx := 1
	dy := 0
	segLen := 1
	segPass := 0
	grid := world.Coords{}
	for {
		_, ok := set.PuzzGrid[grid]
		if !ok {
			return grid
		}
		// set max/min
		if grid.X < set.GridMin.X {
			set.GridMin.X = grid.X
		}
		if grid.Y < set.GridMin.Y {
			set.GridMin.Y = grid.Y
		}
		if grid.X > set.GridMax.X {
			set.GridMax.X = grid.X
		}
		if grid.Y > set.GridMax.Y {
			set.GridMax.Y = grid.Y
		}
		// make next step
		grid.X += dx
		grid.Y += dy
		segPass++
		if segPass == segLen { // done w/segment
			segPass = 0
			// rotate
			dx, dy = dy, -dx
			if dy == 0 { // increase segment length
				segLen++
			}
		}
	}
}

func (set *PuzzleSet) UpdateGridMaxMin() {
	set.GridMin = world.Coords{}
	set.GridMax = world.Coords{}
	for grid := range set.PuzzGrid {
		// set max/min
		if grid.X < set.GridMin.X {
			set.GridMin.X = grid.X
		}
		if grid.Y < set.GridMin.Y {
			set.GridMin.Y = grid.Y
		}
		if grid.X > set.GridMax.X {
			set.GridMax.X = grid.X
		}
		if grid.Y > set.GridMax.Y {
			set.GridMax.Y = grid.Y
		}
	}
}

func (set *PuzzleSet) GetGrid(grid world.Coords) int {
	if i, ok := set.PuzzGrid[grid]; ok {
		return i
	}
	return -1
}

func (set *PuzzleSet) GetGridPuzzle(grid world.Coords) *Puzzle {
	if i, ok := set.PuzzGrid[grid]; ok {
		return set.Puzzles[i]
	}
	return nil
}
