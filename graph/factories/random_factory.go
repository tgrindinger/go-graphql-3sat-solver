package factories

import (
	"math/rand"
	"time"
)

type RandomFactory interface {
	Build() *rand.Rand
}

type TimeRandomFactory struct {
}

func (f *TimeRandomFactory) Build() *rand.Rand {
	return rand.New(rand.NewSource(int64(time.Now().UnixNano())))
}

type ZeroRandomFactory struct {
	staticRandom *rand.Rand
}

func (f *ZeroRandomFactory) Build() *rand.Rand {
	if f.staticRandom == nil {
		f.staticRandom = rand.New(rand.NewSource(0))
	}
	return f.staticRandom
}
