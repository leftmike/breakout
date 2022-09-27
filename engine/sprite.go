package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type ImageSprite struct {
	Hidden        bool
	X, Y          float64
	DX, DY        float64
	width, height float64
	Image         *ebiten.Image
	img           *ebiten.Image
	deleted       bool
}

func (sprt *ImageSprite) Init(mode Mode) {
	if sprt.img != sprt.Image {
		sprt.img = sprt.Image
		w, h := sprt.img.Size()
		sprt.width = float64(w)
		sprt.height = float64(h)
	}
}

func (sprt *ImageSprite) Update(mode Mode) bool {
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

func (sprt *ImageSprite) Corner() (float64, float64) {
	return sprt.X, sprt.Y
}

func (sprt *ImageSprite) Size() (float64, float64) {
	return sprt.width, sprt.height
}

func (sprt *ImageSprite) Draw(mode Mode, screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if sprt.Hidden {
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

type TextSprite struct {
	Hidden        bool
	X, Y          float64
	width, height float64
	Text          string
	Face          font.Face
	Color         color.Color
	deleted       bool
}

func (sprt *TextSprite) Init(mode Mode) {
	// Nothing
}

func (sprt *TextSprite) Update(mode Mode) bool {
	rect := text.BoundString(sprt.Face, sprt.Text)
	sprt.width = float64(rect.Max.X - rect.Min.X)
	sprt.height = float64(rect.Max.Y - rect.Min.Y)
	return false
}

func (sprt *TextSprite) Visible() bool {
	return !sprt.Hidden
}

func (sprt *TextSprite) Collision(with Sprite) {
	// Nothing
}

func (sprt *TextSprite) Corner() (float64, float64) {
	return sprt.X, sprt.Y
}

func (sprt *TextSprite) Size() (float64, float64) {
	return sprt.width, sprt.height
}

func (sprt *TextSprite) Draw(mode Mode, screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if sprt.Hidden || sprt.Text == "" {
		return
	}

	op.ColorM.ScaleWithColor(sprt.Color)
	op.GeoM.Translate(sprt.X, sprt.Y+sprt.height)
	text.DrawWithOptions(screen, sprt.Text, sprt.Face, op)
}

func (sprt *TextSprite) Deleted() bool {
	return sprt.deleted
}

func (sprt *TextSprite) Delete() {
	sprt.deleted = true
}

type RectSprite struct {
	Hidden        bool
	X, Y          float64
	Width, Height float64
	Color         color.RGBA
	width, height float64
	clr           color.RGBA
	img           *ebiten.Image
	deleted       bool
}

func (sprt *RectSprite) Init(mode Mode) {
	if sprt.width != sprt.Width || sprt.height != sprt.Height || sprt.clr != sprt.Color {
		sprt.width = sprt.Width
		sprt.height = sprt.Height
		sprt.clr = sprt.Color
		if sprt.width > 0 && sprt.height > 0 {
			sprt.img = ebiten.NewImage(int(sprt.width), int(sprt.height))
			sprt.img.Fill(sprt.clr)
		} else {
			sprt.img = nil
		}
	}
}

func (sprt *RectSprite) Update(mode Mode) bool {
	return false
}

func (sprt *RectSprite) Visible() bool {
	return !sprt.Hidden
}

func (sprt *RectSprite) Collision(with Sprite) {
	// Nothing
}

func (sprt *RectSprite) Corner() (float64, float64) {
	return sprt.X, sprt.Y
}

func (sprt *RectSprite) Size() (float64, float64) {
	return sprt.width, sprt.height
}

func (sprt *RectSprite) Draw(mode Mode, screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if sprt.Hidden || sprt.img == nil {
		return
	}

	op.GeoM.Translate(sprt.X, sprt.Y)
	screen.DrawImage(sprt.img, op)
}

func (sprt *RectSprite) Deleted() bool {
	return sprt.deleted
}

func (sprt *RectSprite) Delete() {
	sprt.deleted = true
}
