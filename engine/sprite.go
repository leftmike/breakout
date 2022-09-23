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
	Width, Height float64
	Image         *ebiten.Image
	deleted       bool
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

func (sprt *ImageSprite) Corners() (float64, float64, float64, float64) {
	return sprt.X, sprt.Y, sprt.X + sprt.Width, sprt.Y + sprt.Height
}

func (sprt *ImageSprite) Center() (float64, float64) {
	return sprt.X + sprt.Width/2, sprt.Y + sprt.Height/2
}

func (sprt *ImageSprite) Draw(mode Mode, screen *ebiten.Image, op *ebiten.DrawImageOptions) {
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

type TextSprite struct {
	Hidden        bool
	X, Y          float64
	width, height float64
	Text          string
	Face          font.Face
	Color         color.Color
	deleted       bool
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

func (sprt *TextSprite) Corners() (float64, float64, float64, float64) {
	return sprt.X, sprt.Y, sprt.X + sprt.width, sprt.Y + sprt.height
}

func (sprt *TextSprite) Center() (float64, float64) {
	return sprt.X + sprt.width/2, sprt.Y + sprt.height/2
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
