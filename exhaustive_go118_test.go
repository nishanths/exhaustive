//go:build go1.18
// +build go1.18

package exhaustive

import (
	"testing"
)

func TestExhaustiveGo118(t *testing.T) {
	runTest(t, "typeparam/...")
}
