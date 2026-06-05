package data

type Boss interface {
	GetID() string
	SetState(BossState)
	GetState() string
	Update()
	Reset()
	IsDefeated() bool
	Destroy()
}

var (
	CurrentBoss  Boss
	EditorBoss   Boss
	BossCounter0 int
	BossCounter1 int
)

type BossState int

const (
	BossStart = iota
	BossIntro
	BossWaiting
	BossAction
	BossDying
	BossDefeated
	BossPreview
)

func (bs BossState) String() string {
	switch bs {
	case BossStart:
		return "BossStart"
	case BossIntro:
		return "BossIntro"
	case BossWaiting:
		return "BossWaiting"
	case BossAction:
		return "BossAction"
	case BossDying:
		return "BossDying"
	case BossDefeated:
		return "BossDefeated"
	case BossPreview:
		return "BossPreview"
	default:
		return "Unknown"
	}
}
