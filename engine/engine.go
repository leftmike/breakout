package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Layers []*Layer
}

type Layer struct {
	Hidden  []bool
	Paused  []bool
	Sprites []Sprite
	updated []Sprite
}

type Mode int

type Sprite interface {
	Update(mode Mode) bool
	Visible() bool
	Collision(with Sprite)
	Corners() (float64, float64, float64, float64)
	Center() (float64, float64)
	Draw(mode Mode, screen *ebiten.Image, op *ebiten.DrawImageOptions)
	Deleted() bool
}

func (lvl *Level) Update(mode Mode) {
	for _, lyr := range lvl.Layers {
		lyr.update(mode)
	}
}

func (lvl *Level) Draw(mode Mode, screen *ebiten.Image) {
	for _, lyr := range lvl.Layers {
		lyr.draw(mode, screen)
	}
}

func (lyr *Layer) update(mode Mode) {
	if int(mode) < len(lyr.Paused) && lyr.Paused[mode] {
		return
	}

	if lyr.updated == nil {
		lyr.updated = make([]Sprite, 0, len(lyr.Sprites))
	} else {
		lyr.updated = lyr.updated[:0]
	}

	for _, sprt := range lyr.Sprites {
		if sprt.Deleted() {
			continue
		}

		if sprt.Update(mode) {
			lyr.updated = append(lyr.updated, sprt)
		}
	}

	for _, sprt := range lyr.updated {
		if sprt.Deleted() {
			continue
		}

		if sprt.Visible() {
			sprtMinX, sprtMinY, sprtMaxX, sprtMaxY := sprt.Corners()
			for _, with := range lyr.Sprites {
				if with == sprt || with.Deleted() {
					continue
				}

				if with.Visible() {
					withMinX, withMinY, withMaxX, withMaxY := with.Corners()
					if sprtMinX < withMaxX && withMinX < sprtMaxX && sprtMinY < withMaxY &&
						withMinY < sprtMaxY {

						sprt.Collision(with)
					}
				}
			}
		}
	}
}

func (lyr *Layer) draw(mode Mode, screen *ebiten.Image) {
	if int(mode) < len(lyr.Hidden) && lyr.Hidden[mode] {
		return
	}

	cnt := 0
	for _, sprt := range lyr.Sprites {
		if sprt.Deleted() {
			continue
		}

		var op ebiten.DrawImageOptions
		sprt.Draw(mode, screen, &op)
		lyr.Sprites[cnt] = sprt
		cnt += 1
	}

	lyr.Sprites = lyr.Sprites[:cnt]
}

type CollideDirection int

const (
	CollideXGreater CollideDirection = iota
	CollideXLess
	CollideYGreater
	CollideYLess
)

// Collide returns sprt relative to with.
func Collide(sprt, with Sprite) CollideDirection {
	x1, y1 := sprt.Center()
	x2, y2 := with.Center()
	if math.Abs(x1-x2) > math.Abs(y1-y2) {
		if x1 > x2 {
			return CollideXGreater
		}
		return CollideXLess
	}

	if y1 > y2 {
		return CollideYGreater
	}
	return CollideYLess
	/*
		at2 := math.Atan2(y2-y1, x2-x1)
		if at2 > math.Pi/4 && at2 < math.Pi*3/4 {
			return CollideYLess
		} else if at2 > -math.Pi*3/4 && at2 < -math.Pi/4 {
			return CollideYGreater
		} else if at2 > math.Pi*3/4 || at2 < -math.Pi*3/4 {
			return CollideXGreater
		}
		return CollideXLess
	*/
}
