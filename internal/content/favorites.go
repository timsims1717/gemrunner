package content

import (
	"bufio"
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"github.com/pkg/errors"
	"os"
	"sync"
)

var favMu sync.Mutex

// LoadFavoritesFile populates the favorites list.
// As this doesn't use the Mutex, don't run this in
// go routine.
func LoadFavoritesFile() {
	filename := fmt.Sprintf("%s/%s", constants.ContentDir, constants.Favorites)
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// create file
		_, err = os.Create(filename)
		if err != nil {
			fmt.Println("ERROR: Couldn't create favorites file")
			return
		}
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR: Couldn't open favorites file")
		return
	}
	s := bufio.NewScanner(file)
	for s.Scan() {
		data.FavoritesList = append(data.FavoritesList, s.Text())
	}
}

// SaveFavoritesFile saves the current list of favorites.
// This only is safe for the file access, but could be unsafe
// if another goroutine wrote to the favorites list.
// Currently not a big deal.
func SaveFavoritesFile() {
	favMu.Lock()
	defer favMu.Unlock()
	filename := fmt.Sprintf("%s/%s", constants.ContentDir, constants.Favorites)
	favFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("ERROR: couldn't save favorites file")
		return
	}
	for _, line := range data.FavoritesList {
		favFile.WriteString(line)
		favFile.WriteString("\n")
	}
}
