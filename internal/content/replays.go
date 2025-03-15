package content

import (
	"encoding/json"
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/debug"
	"github.com/pkg/errors"
	"os"
	"sync"
	"time"
)

var replayMu sync.Mutex

func LoadReplay(filename string) error {
	replayMu.Lock()
	defer replayMu.Unlock()
	errMsg := "open replay"
	if filename == "" {
		return errors.Wrap(errors.New("no filename provided"), errMsg)
	}
	replayPath := fmt.Sprintf("%s/%s", constants.ReplaysDir, filename)
	replayFile, err := os.ReadFile(replayPath)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	data.CurrReplay = &data.LevelReplay{}
	err = json.Unmarshal(replayFile, data.CurrReplay)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	if debug.Verbose {
		fmt.Printf("INFO: loaded replay from %s\n", replayPath)
	}
	return nil
}

func SaveReplay(replay *data.LevelReplay) {
	replayMu.Lock()
	defer replayMu.Unlock()
	if replay == nil {
		fmt.Println("ERROR: no replay to save")
		return
	}
	replayPath := fmt.Sprintf("%s/%s", constants.ReplaysDir, replay.ReplayFile)
	replayFile, err := os.Create(replayPath)
	if err != nil {
		fmt.Println("ERROR: could not save replay:", err)
		return
	}
	bts, err := json.Marshal(replay)
	if err != nil {
		fmt.Println("ERROR: could not save replay:", err)
		return
	}
	_, err = replayFile.Write(bts)
	if err != nil {
		fmt.Println("ERROR: could not save replay:", err)
		return
	}
	if debug.Verbose {
		fmt.Printf("INFO: saved replay to %s\n", replayPath)
	}
}

func ReplayFile(puzzleSet string, puzzleIndex int) string {
	return fmt.Sprintf(constants.ReplayPath, puzzleSet, puzzleIndex, time.Now().Format("2006.01.02.15.04.05"))
}
