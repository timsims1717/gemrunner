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
)

var savegameMu sync.Mutex

func LoadSaveGame(filename string) error {
	savegameMu.Lock()
	defer savegameMu.Unlock()
	errMsg := "open save session"
	if filename == "" {
		return errors.Wrap(errors.New("no filename provided"), errMsg)
	}
	svgPath := fmt.Sprintf("%s/%s", constants.SavesDir, filename)
	svgFile, err := os.ReadFile(svgPath)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	data.CurrLevelSess = &data.LevelSession{}
	err = json.Unmarshal(svgFile, data.CurrLevelSess)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	if debug.Verbose {
		fmt.Printf("INFO: loaded session from %s\n", svgPath)
	}
	return nil
}

func SaveSaveGame() {
	savegameMu.Lock()
	defer savegameMu.Unlock()
	if data.CurrLevelSess == nil {
		fmt.Println("ERROR: no session to save")
		return
	}
	svgFName := data.CurrLevelSess.Filename
	svgPath := fmt.Sprintf("%s/%s", constants.SavesDir, svgFName)
	svgFile, err := os.Create(svgPath)
	if err != nil {
		fmt.Println("ERROR: could not save session:", err)
		return
	}
	bts, err := json.Marshal(data.CurrLevelSess)
	if err != nil {
		fmt.Println("ERROR: could not save session:", err)
		return
	}
	_, err = svgFile.Write(bts)
	if err != nil {
		fmt.Println("ERROR: could not save session:", err)
		return
	}
	if debug.Verbose {
		fmt.Printf("INFO: saved session to %s\n", svgPath)
	}
}
