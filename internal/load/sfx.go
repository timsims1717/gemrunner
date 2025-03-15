package load

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/sfx"
)

func SoundEffects() {
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-confirm-02.wav", constants.SFXMainButton)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-boxing-punch.wav", constants.SFXPlaceTile)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-cancel-03.wav", constants.SFXConfirm)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-no-go.wav", constants.SFXCancel)
	//sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-alert-03.wav", constants.SFXCancel)
	//sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-cancel-02.wav", constants.SFXCancel)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-select.wav", constants.SFXScroll)

	sfx.SoundPlayer.RegisterSound("assets/sfx/gem.ogg", constants.SFXGem)
	sfx.SoundPlayer.RegisterSound("assets/sfx/pick-up-item.ogg", constants.SFXItem)
	sfx.SoundPlayer.RegisterSound("assets/sfx/drop-item.ogg", constants.SFXDrop)
	sfx.SoundPlayer.RegisterSound("assets/sfx/player-fall1.ogg", constants.SFXFall)
	sfx.SoundPlayer.RegisterSound("assets/sfx/player-land.ogg", constants.SFXLand)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-voice.wav", constants.SFXPlayerCrush)
	sfx.SoundPlayer.RegisterSound("assets/sfx/jump.ogg", constants.SFXJump)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-bump.wav", constants.SFXBump)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-increase.wav", constants.SFXRegen)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-chip-thud.wav", constants.SFXHit)

	sfx.SoundPlayer.RegisterSound("assets/sfx/throw-box.ogg", constants.SFXThrow)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-metal-hit.wav", constants.SFXBoxLand)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-funny-switch.wav", constants.SFXKey)
	sfx.SoundPlayer.RegisterSound("assets/sfx/doors-open.ogg", constants.SFXBombLight)
	sfx.SoundPlayer.RegisterSound("assets/sfx/explode.ogg", constants.SFXBombBlow)
	sfx.SoundPlayer.RegisterSound("assets/sfx/jetpack-b.ogg", constants.SFXJetpackStart)
	sfx.SoundPlayer.RegisterSound("assets/sfx/jetpack-e.ogg", constants.SFXJetpackEnd)
	sfx.SoundPlayer.RegisterSound("assets/sfx/flamethrower.ogg", constants.SFXFlamethrower)
	sfx.SoundPlayer.RegisterSound("assets/sfx/trans_in.ogg", constants.SFXTransIn)
	sfx.SoundPlayer.RegisterSound("assets/sfx/trans_out.ogg", constants.SFXTransOut)

	sfx.SoundPlayer.RegisterSound("assets/sfx/crushed.ogg", constants.SFXCrush)
	sfx.SoundPlayer.RegisterSound("assets/sfx/doors-open2.ogg", constants.SFXDoorsOpen)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-bulletin.wav", constants.SFXFanfare1)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-agile.wav", constants.SFXFanfare2)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-sublime.wav", constants.SFXFanfare3)
	sfx.SoundPlayer.RegisterSound("assets/sfx/exit-level.ogg", constants.SFXExitLevel)

	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-unknown.wav", constants.SFXDemonAttack)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-magic-exploding.wav", constants.SFXFlyBlow)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-drop-02.wav", constants.SFXFlyCrush)
}
