package content

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/options"
	"gemrunner/pkg/sfx"
	"github.com/BurntSushi/toml"
	"os"
	"sync"
)

var (
	configMu sync.Mutex
)

func CreateConfig() {
	configMu.Lock()
	defer configMu.Unlock()
	os.Remove(constants.ConfigFile)
	file, err := os.Create(constants.ConfigFile)
	if err != nil {
		panic(err)
	}
	encode := toml.NewEncoder(file)
	err = encode.Encode(constants.DefaultConfiguration)
	if err != nil {
		panic(err)
	} else if debug.Verbose {
		fmt.Printf("INFO: created configuration at %s\n", constants.ConfigFile)
	}
	UpdateConfiguration()
}

func LoadConfig() {
	configMu.Lock()
	defer configMu.Unlock()
	if _, err := toml.DecodeFile(constants.ConfigFile, &constants.Configuration); err != nil {
		fmt.Printf("ERROR: couldn't decode configuration file: %s\n", err)
		CreateConfig()
	} else if debug.Verbose {
		fmt.Printf("INFO: saved configuration to %s\n", constants.ConfigFile)
	}
	UpdateConfiguration()
}

func SaveConfig() {
	configMu.Lock()
	defer configMu.Unlock()
	os.Remove(constants.ConfigFile)
	file, err := os.Create(constants.ConfigFile)
	if err != nil {
		panic(err)
	}

	encode := toml.NewEncoder(file)
	err = encode.Encode(constants.Configuration)
	if err != nil {
		fmt.Printf("ERROR: couldn't save configuration: %s\n", err)
	} else if debug.Verbose {
		fmt.Printf("INFO: saved configuration to %s\n", constants.ConfigFile)
	}
}

func UpdateConfiguration() {
	options.VSync = constants.Configuration.Graphics.VSync
	options.FullScreen = constants.Configuration.Graphics.Fullscreen
	options.BilinearFilter = constants.Configuration.Graphics.BilinearFilter
	options.ResolutionIndex = constants.Configuration.Graphics.Resolution

	uOn := int32(0)
	if constants.Configuration.Graphics.Scanlines {
		uOn = int32(1)
	}
	constants.Scanlines = uOn

	// change the shader settings here

	sfx.SetMusicVolume(constants.Configuration.Audio.MusicVolume)
	sfx.MuteMusic(!constants.Configuration.Audio.MusicOn)
	sfx.SetSoundVolume(constants.Configuration.Audio.SfxVolume)
	sfx.MuteSound(!constants.Configuration.Audio.SfxOn)
	sfx.SetMasterVolume(constants.Configuration.Audio.MasterVolume)
	sfx.MuteMaster(!constants.Configuration.Audio.MasterOn)
}
