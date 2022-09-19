package engine_test

import (
	"testing"

	"github.com/leftmike/breakout/engine"
)

func TestCollide(t *testing.T) {
	cases := []struct {
		x, y float64
		dir  engine.CollideDirection
	}{
		{1, 0, engine.CollideXGreater},
		{-1, 0, engine.CollideXLess},
		{0, -1, engine.CollideYLess},
		{0, 1, engine.CollideYGreater},
	}

	for _, c := range cases {
		dir := engine.Collide(&engine.ImageSprite{X: c.x, Y: c.y}, &engine.ImageSprite{})
		if dir != c.dir {
			t.Errorf("Collide(%v, %v): got %v want %v", c.x, c.y, dir, c.dir)
		}
	}
}
