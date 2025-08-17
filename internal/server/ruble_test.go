package server

import (
	"testing"
)

func TestRuble(t *testing.T) {
	t.Run("test add", func(t *testing.T) {
		a := NewRuble(.1)
		b := NewRuble(.2)

		sum := a.Add(b)
		if sum.Float64() != .3 {
			t.Errorf("got %v want 0.3", sum)
		}
	})

	t.Run("test subtract", func(t *testing.T) {
		a := NewRuble(10)
		b := NewRuble(5.5)

		sub := a.Sub(b)
		if sub.Float64() != 4.5 {
			t.Errorf("got %v want 4.5", sub.Float64())
		}
	})
}
