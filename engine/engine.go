package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Layers []*Layer
}

type Layer struct {
	Hidden  bool
	Sprites []*Sprite
}

type Sprite struct {
	Hidden    bool
	X, Y      float64
	DX, DY    float64
	W, H      float64
	Image     *ebiten.Image
	Collision func(sprt, with *Sprite)
}

func (lvl *Level) Update() {
	for _, lyr := range lvl.Layers {
		lyr.update()
	}
}

func (lvl *Level) Draw(screen *ebiten.Image) {
	for _, lyr := range lvl.Layers {
		if lyr.Hidden {
			continue
		}

		lyr.draw(screen)
	}
}

func (lyr *Layer) update() {
	for _, sprt := range lyr.Sprites {
		sprt.X += sprt.DX
		sprt.Y += sprt.DY
	}

	for _, sprt := range lyr.Sprites {
		if sprt.Hidden {
			continue
		}

		if sprt.Collision != nil {
			for _, with := range lyr.Sprites {
				if with == sprt || with.Hidden {
					continue
				}
				if sprt.overlaps(with) {
					sprt.Collision(sprt, with)
				}
			}
		}
	}
}

func (lyr *Layer) draw(screen *ebiten.Image) {
	for _, sprt := range lyr.Sprites {
		if sprt.Hidden || sprt.Image == nil {
			continue
		}

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(sprt.X, sprt.Y)
		screen.DrawImage(sprt.Image, &op)
	}
}

func (sprt *Sprite) NewImageFill(w, h int, clr color.Color) *Sprite {
	sprt.Image = ebiten.NewImage(w, h)
	sprt.Image.Fill(clr)
	sprt.W = float64(w)
	sprt.H = float64(h)
	return sprt
}

func (sprt *Sprite) overlaps(with *Sprite) bool {
	sprtMaxX := sprt.X + sprt.W
	sprtMaxY := sprt.Y + sprt.H
	withMaxX := with.X + with.W
	withMaxY := with.Y + with.H

	return sprt.X < withMaxX && with.X < sprtMaxX && sprt.Y < withMaxY && with.Y < sprtMaxY
}
