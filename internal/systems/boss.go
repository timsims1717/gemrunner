package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"github.com/gopxl/pixel"
)

func BossSystem() {
	if data.CurrentBoss != nil {
		data.CurrentBoss.Update()
		if data.CurrentBoss.IsDefeated() {
			data.CurrentBoss.Destroy()
			data.CurrentBoss = nil
		}
	}
}

func InitBoss(level *data.Level) {
	switch level.Puzzle.Metadata.Boss {
	case constants.BossBlob:
		boss := CreateBlobBoss(12, pixel.ToRGBA(constants.ColorGreen), pixel.ToRGBA(constants.ColorDarkGreen))
		data.CurrentBoss = boss
	}
}

func BossEditorSystem() {
	if data.EditorBoss != nil {
		data.EditorBoss.Update()
	}
}

func RemoveEditorBoss() {
	if data.EditorBoss != nil {
		data.EditorBoss.Destroy()
		data.EditorBoss = nil
	}
}

func SetEditorBoss(id string) {
	if id == "" {
		return
	}
	var boss data.Boss
	switch id {
	case constants.BossBlob:
		boss = CreateBlobBoss(12, pixel.ToRGBA(constants.ColorGreen), pixel.ToRGBA(constants.ColorDarkGreen))
		data.CurrentBoss = boss
	}
	boss.SetState(data.BossPreview)
	data.CurrPuzzleSet.CurrPuzzle.Metadata.Boss = boss.GetID()
	data.EditorBoss = boss
}
