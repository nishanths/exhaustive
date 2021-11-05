package someone

import "goplay/pkg"

// type OwnY uint -->
type SomeonesOwnY = pkg.OwnY

const (
	SomeonesMM SomeonesOwnY = pkg.MemM
	SomeonesMN SomeonesOwnY = pkg.MemN
	SomeonesMO SomeonesOwnY = 444 // new extra here
)

func P() SomeonesOwnY {
	return SomeonesOwnY(100)
}
