package content

import (
	"encoding/json"
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/util"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var localPzlMu sync.Mutex

func LoadLocalPuzzleList() {
	localPzlMu.Lock()
	defer localPzlMu.Unlock()
	list, err := os.ReadDir(constants.PuzzlesDir)
	if err != nil {
		fmt.Printf("ERROR: couldn't read local puzzle list: %s\n", err)
		return
	}
	data.PuzzleSetFileList = []data.PuzzleSetMetadata{}
	for _, d := range list {
		if !d.IsDir() && filepath.Ext(d.Name()) == constants.PuzzleExt {
			filename := fmt.Sprintf("%s/%s", constants.PuzzlesDir, d.Name())
			pzlFile, err := os.ReadFile(filename)
			if err != nil {
				continue
			}
			pzs := &data.PuzzleSet{}
			err = json.Unmarshal(pzlFile, pzs)
			if err != nil {
				continue
			}
			pi := pzs.Metadata
			if pi.Filename == "" {
				pi.Filename = d.Name()
			}
			if pi.Name == "" {
				pi.Name = "No Name"
			}
			if pi.Author == "" {
				pi.Author = "Unknown"
			}
			if pi.NumPlayers == 0 {
				pi.NumPlayers = 1
			}
			if pi.NumPuzzles == 0 {
				pi.NumPuzzles = len(pzs.Puzzles)
			}
			if pi.UUID != nil {
				pi.Favorite = util.ContainsStr(pi.UUID.String(), data.FavoritesList)
			}
			data.PuzzleSetFileList = append(data.PuzzleSetFileList, pi)
		}
	}
}

type PuzzleOrder int

const (
	MostRecentPlay = iota
	LeastRecentPlay
	MostRecentAdd
	LeastRecentAdd
	MostPuzzles
	LeastPuzzles
	AtoZ
	ZtoA
	MostPopular
	LeastPopular
	MostFavCount
	LeastFavCount
)

type PuzzleFilters struct {
	Local     bool
	Online    bool
	Favorite  bool
	PlayerCnt int
	Search    string
	Ordering  PuzzleOrder
}

var DefaultFilters = PuzzleFilters{
	Local:     false,
	Online:    false,
	Favorite:  false,
	PlayerCnt: 0,
	Search:    "",
	Ordering:  AtoZ,
}

func OrganizeLocalPuzzles(pzFilters PuzzleFilters) {
	localPzlMu.Lock()
	defer localPzlMu.Unlock()
	data.PuzzleSetSortedList = []data.PuzzleSetMetadata{}
	for _, pzl := range data.PuzzleSetFileList {
		inList := true
		if pzFilters.Local && !pzl.Local {
			inList = false
		} else if pzFilters.Online && !pzl.Online {
			inList = false
		} else if pzFilters.Favorite && !pzl.Favorite {
			inList = false
		} else if pzFilters.PlayerCnt > 0 && pzFilters.PlayerCnt != pzl.NumPlayers {
			inList = false
		} else if pzFilters.Search != "" {
			found := false
			if strings.Contains(pzl.Name, pzFilters.Search) ||
				strings.Contains(pzl.Author, pzFilters.Search) ||
				strings.Contains(pzl.Desc, pzFilters.Search) {
				found = true
			}
			if !found {
				inList = false
			}
		}
		if inList {
			data.PuzzleSetSortedList = append(data.PuzzleSetSortedList, pzl)
		}
	}
	sortOrderedList(pzFilters.Ordering)
}

func sortOrderedList(o PuzzleOrder) {
	i := 1
	for i < len(data.PuzzleSetSortedList) {
		j := 1
		var less bool
		for j >= 1 {
			switch o {
			case MostRecentPlay:
				if data.PuzzleSetSortedList[j].RecentPlay == nil {
					less = false
				} else if data.PuzzleSetSortedList[j-1].RecentPlay == nil {
					less = true
				} else {
					less = data.PuzzleSetSortedList[j].RecentPlay.After(*data.PuzzleSetSortedList[j-1].RecentPlay)
				}
			case LeastRecentPlay:
				if data.PuzzleSetSortedList[j].RecentPlay == nil {
					less = true
				} else if data.PuzzleSetSortedList[j-1].RecentPlay == nil {
					less = false
				} else {
					less = data.PuzzleSetSortedList[j].RecentPlay.Before(*data.PuzzleSetSortedList[j-1].RecentPlay)
				}
			case MostRecentAdd:
				if data.PuzzleSetSortedList[j].RecordDate == nil {
					less = false
				} else if data.PuzzleSetSortedList[j-1].RecordDate == nil {
					less = true
				} else {
					less = data.PuzzleSetSortedList[j].RecordDate.After(*data.PuzzleSetSortedList[j-1].RecordDate)
				}
			case LeastRecentAdd:
				if data.PuzzleSetSortedList[j].RecordDate == nil {
					less = true
				} else if data.PuzzleSetSortedList[j-1].RecordDate == nil {
					less = false
				} else {
					less = data.PuzzleSetSortedList[j].RecordDate.Before(*data.PuzzleSetSortedList[j-1].RecordDate)
				}
			case MostPuzzles:
				less = data.PuzzleSetSortedList[j].NumPuzzles > data.PuzzleSetSortedList[j-1].NumPuzzles
			case LeastPuzzles:
				less = data.PuzzleSetSortedList[j].NumPuzzles < data.PuzzleSetSortedList[j-1].NumPuzzles
			case AtoZ:
				less = strings.ToLower(data.PuzzleSetSortedList[j].Name) < strings.ToLower(data.PuzzleSetSortedList[j-1].Name)
			case ZtoA:
				less = strings.ToLower(data.PuzzleSetSortedList[j].Name) > strings.ToLower(data.PuzzleSetSortedList[j-1].Name)
			case MostPopular:
				less = data.PuzzleSetSortedList[j].Downloads > data.PuzzleSetSortedList[j-1].Downloads
			case LeastPopular:
				less = data.PuzzleSetSortedList[j].Downloads < data.PuzzleSetSortedList[j-1].Downloads
			case MostFavCount:
				less = data.PuzzleSetSortedList[j].Favorites > data.PuzzleSetSortedList[j-1].Favorites
			case LeastFavCount:
				less = data.PuzzleSetSortedList[j].Favorites < data.PuzzleSetSortedList[j-1].Favorites
			}
			if less {
				data.PuzzleSetSortedList[j], data.PuzzleSetSortedList[j-1] = data.PuzzleSetSortedList[j-1], data.PuzzleSetSortedList[j]
			} else {
				break
			}
			j--
		}
		i++
	}
}

func SavePuzzleSetToFile() error {
	localPzlMu.Lock()
	defer localPzlMu.Unlock()
	errMsg := "save puzzle set"
	if data.CurrPuzzleSet == nil {
		return errors.Wrap(errors.New("no puzzle set to save"), errMsg)
	}
	var pzlFName = "test.puzzle"
	if data.CurrPuzzleSet.Metadata.Filename != "" {
		pzlFName = data.CurrPuzzleSet.Metadata.Filename
	}
	if data.CurrPuzzleSet.Metadata.UUID == nil {
		id := uuid.New()
		data.CurrPuzzleSet.Metadata.UUID = &id
	}
	data.CurrPuzzleSet.Metadata.NumPuzzles = len(data.CurrPuzzleSet.Puzzles)
	pzlPath := fmt.Sprintf("%s/%s", constants.PuzzlesDir, pzlFName)
	pzlFile, err := os.Create(pzlPath)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	bts, err := json.Marshal(data.CurrPuzzleSet)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	_, err = pzlFile.Write(bts)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	fmt.Printf("INFO: saved puzzle set to %s\n", pzlPath)
	return nil
}

func OpenPuzzleSetFile(filename string) error {
	errMsg := "open puzzle set"
	if filename == "" {
		return errors.Wrap(errors.New("no filename provided"), errMsg)
	}
	pzlFName := fmt.Sprintf("%s/%s", constants.PuzzlesDir, filename)
	pzlFile, err := os.ReadFile(pzlFName)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	data.CurrPuzzleSet = &data.PuzzleSet{}
	err = json.Unmarshal(pzlFile, data.CurrPuzzleSet)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	fmt.Printf("INFO: loaded puzzle set from %s\n", pzlFName)
	return nil
}

func OpenPuzzleSetFileRt(filename string) (*data.PuzzleSet, error) {
	errMsg := "open puzzle set"
	if filename == "" {
		return nil, errors.Wrap(errors.New("no filename provided"), errMsg)
	}
	pzlFName := fmt.Sprintf("%s/%s", constants.PuzzlesDir, filename)
	pzlFile, err := os.ReadFile(pzlFName)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	rtPzlSet := &data.PuzzleSet{}
	err = json.Unmarshal(pzlFile, rtPzlSet)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	fmt.Printf("INFO: loaded puzzle set from %s\n", pzlFName)
	return rtPzlSet, nil
}

func OpenPuzzleFile(filename string) error {
	localPzlMu.Lock()
	defer localPzlMu.Unlock()
	errMsg := "open puzzle"
	if filename == "" {
		return errors.Wrap(errors.New("no filename provided"), errMsg)
	}
	pzlFile, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	data.CurrPuzzleSet = data.CreatePuzzleSet()
	err = json.Unmarshal(pzlFile, data.CurrPuzzleSet.CurrPuzzle)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	if data.CurrPuzzleSet.CurrPuzzle.Metadata.Name == "" {
		return errors.Wrap(errors.New("not a puzzle"), errMsg)
	}
	data.CurrPuzzleSet.Metadata.Name = data.CurrPuzzleSet.CurrPuzzle.Metadata.Name
	data.CurrPuzzleSet.Metadata.Filename = data.CurrPuzzleSet.CurrPuzzle.Metadata.Filename
	fmt.Printf("INFO: loaded puzzle from %s\n", filename)
	return nil
}
