package enum

import (
	"enum/otherpkg"
)

type T2 int

type T1 = T2
type T3 = otherpkg.T4
type T5 = otherpkg.T5
