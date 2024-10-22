package ui

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

func LoadDialog(path string) (*DialogConstructor, error) {
	errMsg := "load dialog"
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer file.Close()
	var dc DialogConstructor
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dc)
	return &dc, err
}
