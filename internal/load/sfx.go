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

	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-coin.wav", constants.SFXGem)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-stuff-up.wav", constants.SFXItem)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-bump-wood.wav", constants.SFXDrop)
	sfx.SoundPlayer.RegisterSound("assets/sfx/player-fall1.ogg", constants.SFXFall)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-boxing-punch.wav", constants.SFXLand)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-voice.wav", constants.SFXPlayerCrush)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-fireball.wav", constants.SFXJump1)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-leap-out.wav", constants.SFXJump2)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-bump.wav", constants.SFXBump)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-increase.wav", constants.SFXRegen)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-chip-thud.wav", constants.SFXHit)

	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-moving-block.wav", constants.SFXThrow)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-metal-hit.wav", constants.SFXBoxLand)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-funny-switch.wav", constants.SFXKey)
	sfx.SoundPlayer.RegisterSound("assets/sfx/doors-open.ogg", constants.SFXBombLight)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-explode.wav", constants.SFXBombBlow)
	sfx.SoundPlayer.RegisterSound("assets/sfx/doors-open.ogg", constants.SFXJetpack)

	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-low-wave.wav", constants.SFXCrush)
	sfx.SoundPlayer.RegisterSound("assets/sfx/doors-open.ogg", constants.SFXDoorsOpen)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-bulletin.wav", constants.SFXFanfare1)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-agile.wav", constants.SFXFanfare2)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-sublime.wav", constants.SFXFanfare3)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-rewind.wav", constants.SFXExitLevel)

	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-unknown.wav", constants.SFXDemonAttack)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-car-hit.wav", constants.SFXDemonCrush)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-magic-exploding.wav", constants.SFXFlyBlow)
	sfx.SoundPlayer.RegisterSound("assets/sfx/NFF-drop-02.wav", constants.SFXFlyCrush)
}
