package toolkit

import "testing"

func TestWalkExt(t *testing.T) {
	t.Log(WalkExt("../", ".go"))
}
