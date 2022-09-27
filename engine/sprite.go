package engine

import (
	"image/color"
	"strings"

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

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type TextSprite struct {
	Hidden     bool
	X, Y       float64
	Text       string
	Align      Align
	Face       font.Face
	Color      color.Color
	Background color.Color
	Margin     float64
	text       string
	lines      []string
	widths     []int
	maxWidth   int
	lineHeight int
	background color.Color
	margin     float64
	img        *ebiten.Image
	deleted    bool
}

func (sprt *TextSprite) Init(mode Mode) {
	if sprt.text != sprt.Text || sprt.background != sprt.Background || sprt.margin != sprt.Margin {
		sprt.text = sprt.Text
		sprt.background = sprt.Background
		sprt.margin = sprt.Margin

		sprt.lines = strings.Split(sprt.text, "\n")
		sprt.widths = nil
		sprt.maxWidth = 0
		for _, line := range sprt.lines {
			rect := text.BoundString(sprt.Face, line)
			w := rect.Max.X - rect.Min.X
			sprt.widths = append(sprt.widths, w)
			if w > sprt.maxWidth {
				sprt.maxWidth = w
			}
		}
		metrics := sprt.Face.Metrics()
		sprt.lineHeight = metrics.Height.Ceil()

		if sprt.background == nil {
			sprt.img = nil
		} else {
			w, h := sprt.Size()
			sprt.img = ebiten.NewImage(int(w), int(h))
			sprt.img.Fill(sprt.background)
		}
	}
}

func (sprt *TextSprite) Update(mode Mode) bool {
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
	return float64(sprt.maxWidth) + sprt.margin*2,
		float64(sprt.lineHeight*len(sprt.lines)) + sprt.margin*2
}

func (sprt *TextSprite) Draw(mode Mode, screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if sprt.Hidden {
		return
	}

	op.GeoM.Translate(sprt.X, sprt.Y)
	if sprt.img != nil {
		screen.DrawImage(sprt.img, op)
	}

	op.GeoM.Translate(sprt.margin, sprt.margin)
	op.ColorM.ScaleWithColor(sprt.Color)
	for cnt, line := range sprt.lines {
		nop := *op
		switch sprt.Align {
		case AlignLeft:
			nop.GeoM.Translate(0, float64((cnt+1)*sprt.lineHeight))
		case AlignCenter:
			nop.GeoM.Translate(float64(sprt.maxWidth-sprt.widths[cnt])/2,
				float64((cnt+1)*sprt.lineHeight))
		case AlignRight:
			nop.GeoM.Translate(float64(sprt.maxWidth-sprt.widths[cnt]),
				float64((cnt+1)*sprt.lineHeight))
		}
		text.DrawWithOptions(screen, line, sprt.Face, &nop)
	}
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
