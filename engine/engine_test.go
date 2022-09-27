package engine_test

import (
	"reflect"
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

func gotSprites(sprites []engine.Sprite) []int {
	var got []int
	for _, sprite := range sprites {
		got = append(got, int(sprite.(*engine.RectSprite).X))
	}
	return got
}

func TestDelete(t *testing.T) {
	sprites := []engine.Sprite{
		&engine.RectSprite{X: 0},
		&engine.RectSprite{X: 1},
		&engine.RectSprite{X: 2},
		&engine.RectSprite{X: 3},
		&engine.RectSprite{X: 4},
		&engine.RectSprite{X: 5},
		&engine.RectSprite{X: 6},
		&engine.RectSprite{X: 7},
	}

	level := engine.Level{
		Layers: []*engine.Layer{
			{Sprites: append([]engine.Sprite{}, sprites...)},
		},
	}

	level.Draw(0, nil)

	got := gotSprites(level.Layers[0].Sprites)
	want := []int{0, 1, 2, 3, 4, 5, 6, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sprites: got %v want %v", got, want)
	}

	sprites[3].(*engine.RectSprite).Delete()
	level.Draw(0, nil)
	got = gotSprites(level.Layers[0].Sprites)
	want = []int{0, 1, 2, 4, 5, 6, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete(3): got %v want %v", got, want)
	}

	sprites[7].(*engine.RectSprite).Delete()
	level.Draw(0, nil)
	got = gotSprites(level.Layers[0].Sprites)
	want = []int{0, 1, 2, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete(7): got %v want %v", got, want)
	}

	sprites[0].(*engine.RectSprite).Delete()
	level.Draw(0, nil)
	got = gotSprites(level.Layers[0].Sprites)
	want = []int{1, 2, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete(0): got %v want %v", got, want)
	}

	level.Layers[0].Sprites = append(level.Layers[0].Sprites, &engine.RectSprite{X: 8})
	level.Draw(0, nil)
	got = gotSprites(level.Layers[0].Sprites)
	want = []int{1, 2, 4, 5, 6, 8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Add(8): got %v want %v", got, want)
	}

	sprites[4].(*engine.RectSprite).Delete()
	level.Draw(0, nil)
	got = gotSprites(level.Layers[0].Sprites)
	want = []int{1, 2, 5, 6, 8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete(4): got %v want %v", got, want)
	}
}
