package shape

import (
	"testing"
)

func TestRectangle(t *testing.T) {
	r1 := NewRectangle(0, 0, 1, 1)
	r2 := NewRectangle(1, 0, 2, 1)
	if !r1.Overlaps(r1) {
		t.Error("r1 doesnt not overlap r1")
	}
	if r1.Overlaps(r2) {
		t.Error("r1 != r2")
	}
	if !r1.Contains(0, 0) {
		t.Error("r1 doesnt contain (0,0)")
	}
}
