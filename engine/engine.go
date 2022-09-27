package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Layers []*Layer
}

type Layer struct {
	Visible []bool
	Active  []bool
	Sprites []Sprite
	updated []Sprite
}

type Mode int

type Sprite interface {
	Init(mode Mode)
	Update(mode Mode) bool
	Visible() bool
	Collision(with Sprite)
	Corner() (float64, float64)
	Size() (float64, float64)
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
	for _, sprt := range lyr.Sprites {
		if sprt.Deleted() {
			continue
		}

		sprt.Init(mode)
	}

	if lyr.Active != nil && (int(mode) >= len(lyr.Active) || !lyr.Active[mode]) {
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
			sprtX, sprtY := sprt.Corner()
			sprtW, sprtH := sprt.Size()
			for _, with := range lyr.Sprites {
				if with == sprt || with.Deleted() {
					continue
				}

				if with.Visible() {
					withX, withY := with.Corner()
					withW, withH := with.Size()
					if sprtX < withX+withW && withX < sprtX+sprtW &&
						sprtY < withY+withH && withY < sprtY+sprtH {

						sprt.Collision(with)
					}
				}
			}
		}
	}
}

func (lyr *Layer) draw(mode Mode, screen *ebiten.Image) {
	if lyr.Visible != nil && (int(mode) >= len(lyr.Visible) || !lyr.Visible[mode]) {
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
	x1, y1 := sprt.Corner()
	w, h := sprt.Size()
	x1 += w / 2
	y1 += h / 2

	x2, y2 := with.Corner()
	w, h = with.Size()
	x2 += w / 2
	y2 += h / 2

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
