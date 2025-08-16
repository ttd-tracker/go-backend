package server

import (
	"testing"
)

func TestRuble(t *testing.T) {
	t.Run("test using decimal", func(t *testing.T) {
		a := NewRuble(.1)
		b := NewRuble(.2)

		sum := a.Add(b)
		if sum.Float64() != .3 {
			t.Errorf("got %v want 0.3", sum)
		}
	})
}
