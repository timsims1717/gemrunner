package content

import (
	"fmt"
	"gemrunner/internal/constants"
	"os"
	"os/user"
	"runtime"
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	constants.HomeDir = usr.HomeDir
	constants.ContentDir = constants.HomeDir
	constants.System = runtime.GOOS
	switch constants.System {
	case "windows":
		fmt.Println("OS: Windows")
		constants.ContentDir += constants.WinDir
	case "darwin":
		fmt.Println("OS: Mac")
		constants.ContentDir += constants.MacDir
	case "linux":
		fmt.Println("OS: Linux")
		constants.ContentDir += constants.LinuxDir
	default:
		fmt.Printf("OS Unknown: %s.\n", constants.System)
		constants.ContentDir += constants.LinuxDir
	}
	err = os.MkdirAll(constants.ContentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	constants.SavesDir = constants.ContentDir + constants.SaveDir
	err = os.MkdirAll(constants.SavesDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	constants.PuzzlesDir = constants.ContentDir + constants.PuzzleDir
	err = os.MkdirAll(constants.PuzzlesDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
