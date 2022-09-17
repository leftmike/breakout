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
	Sprites []Sprite
	updated []Sprite
}

type Sprite interface {
	Update() bool
	Visible() bool
	Collision(with Sprite)
	Corners() (float64, float64, float64, float64)
	Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions)
	Deleted() bool
}

func (lvl *Level) Update() {
	for _, lyr := range lvl.Layers {
		lyr.update()
	}
}

func (lvl *Level) Draw(screen *ebiten.Image) {
	for _, lyr := range lvl.Layers {
		lyr.draw(screen)
	}
}

func (lyr *Layer) update() {
	if lyr.updated == nil {
		lyr.updated = make([]Sprite, 0, len(lyr.Sprites))
	} else {
		lyr.updated = lyr.updated[:0]
	}

	for _, sprt := range lyr.Sprites {
		if sprt.Deleted() {
			continue
		}

		if sprt.Update() {
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

func (lyr *Layer) draw(screen *ebiten.Image) {
	if lyr.Hidden {
		return
	}

	// XXX: cleanup deleted sprites
	for _, sprt := range lyr.Sprites {
		if sprt.Deleted() {
			continue
		}

		var op ebiten.DrawImageOptions
		sprt.Draw(screen, &op)
	}
}

type ImageSprite struct {
	Hidden        bool
	X, Y          float64
	DX, DY        float64
	Width, Height float64
	Image         *ebiten.Image
	deleted       bool
}

func (sprt *ImageSprite) Update() bool {
	if sprt.DX == 0 && sprt.DY == 0 {
		return false
	}

	sprt.X += sprt.DX
	sprt.Y += sprt.DY
	return true
}

func (sprt *ImageSprite) Visible() bool {
	return !sprt.Hidden
}

func (sprt *ImageSprite) Collision(with Sprite) {
	// Nothing
}

func (sprt *ImageSprite) Corners() (float64, float64, float64, float64) {
	return sprt.X, sprt.Y, sprt.X + sprt.Width, sprt.Y + sprt.Height
}

func (sprt *ImageSprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if sprt.Hidden || sprt.Image == nil {
		return
	}

	op.GeoM.Translate(sprt.X, sprt.Y)
	screen.DrawImage(sprt.Image, op)
}

func (sprt *ImageSprite) Deleted() bool {
	return sprt.deleted
}

func (sprt *ImageSprite) Delete() {
	sprt.deleted = true
}

func NewImageSprite(x, y float64, img *ebiten.Image) ImageSprite {
	w, h := img.Size()

	return ImageSprite{
		X:      x,
		Y:      y,
		Width:  float64(w),
		Height: float64(h),
		Image:  img,
	}
}

func NewImageFill(w, h int, clr color.Color) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(clr)

	return img
}
