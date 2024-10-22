package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/ui"
	"gemrunner/pkg/util"
)

func DialogConstructors() {
	for _, key := range constants.DialogKeys {
		path := fmt.Sprintf("assets/ui/%s.json", key)
		dlgCon, err := ui.LoadDialog(path)
		if err != nil {
			fmt.Printf("ERROR: failed to load dialog %s: %s\n", key, err)
		} else {
			ui.DialogConstructors[key] = dlgCon
		}
	}
}

func ReloadDialog(key string) {
	if util.ContainsStr(key, constants.DialogKeys) {
		for k, d := range ui.Dialogs {
			if k == key {
				ui.DisposeDialog(d)
			}
		}
		path := fmt.Sprintf("assets/ui/%s.json", key)
		dlgCon, err := ui.LoadDialog(path)
		if err != nil {
			fmt.Printf("ERROR: failed to reload dialog %s: %s\n", key, err)
		} else {
			ui.DialogConstructors[key] = dlgCon
			ui.NewDialog(ui.DialogConstructors[key])
		}
	}
}
