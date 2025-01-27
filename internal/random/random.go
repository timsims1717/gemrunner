package random

import (
	"fmt"
	"gemrunner/pkg/debug"
	"math/rand"
	"time"
)

var (
	Global  *rand.Rand
	Level   *rand.Rand
	Effects *rand.Rand
)

func init() {
	seed := time.Now().UnixNano()
	//seed := int64(1627363045028136166)
	Global = rand.New(rand.NewSource(seed))
	PrintSeed(seed, "Global")
	effSeed := Global.Int63()
	Effects = rand.New(rand.NewSource(effSeed))
	PrintSeed(effSeed, "Effects")
	Level = rand.New(rand.NewSource(RandomSeed()))
}

func PrintSeed(seed int64, s string) {
	if debug.Verbose {
		fmt.Printf("%s Seed: %d\n", s, seed)
	}
}

func RandGlobalSeed() {
	seed := time.Now().UnixNano()
	Global.Seed(seed)
	PrintSeed(seed, "Global")
}

func SetGlobalSeed(seed int64) {
	Global.Seed(seed)
	PrintSeed(seed, "Global")
}

func RandomSeed() int64 {
	return Global.Int63()
}

func RandLevelSeed() {
	seed := Global.Int63()
	Level.Seed(seed)
	PrintSeed(seed, "Level")
}

func SetLevelSeed(seed int64) {
	Level.Seed(seed)
	PrintSeed(seed, "Level")
}

func RandEffectsSeed() {
	seed := Global.Int63()
	Effects.Seed(seed)
	PrintSeed(seed, "Effects")
}

func SetEffectsSeed(seed int64) {
	Effects.Seed(seed)
	PrintSeed(seed, "Effects")
}
