package main

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/leftmike/breakout/engine"
	"github.com/leftmike/breakout/fonts"
)

const (
	demoMode engine.Mode = iota
	pauseMode
	playMode
)

const (
	windowWidth  = 500
	windowHeight = 500
	screenWidth  = 400
	screenHeight = 400

	paddleWidth  = 80
	paddleHeight = 10
	paddleY      = screenHeight - (paddleHeight + 2)
	paddleX      = (screenWidth - paddleWidth) / 2

	paddleAccel    = 0.5
	paddleDecel    = 1
	paddleMaxSpeed = 20

	ballWidth  = 10
	ballHeight = 10
	ballX      = (screenWidth - ballWidth) / 2
	ballY      = 0

	blockSize   = 30
	blockMargin = 2
	blockBorder = blockSize * 3 / 2
)

var (
	mode = playMode

	start = true // XXX

	errQuit = errorQuit{}

	face   = NewFace(fonts.RobotoRegular(), 24)
	paddle = PaddleSprite{
		ImageSprite: engine.ImageSprite{
			X:     paddleX,
			Y:     paddleY,
			Image: newImageFill(paddleWidth, paddleHeight, color.RGBA{0, 0xFF, 0, 0xFF}),
		},
	}
	ballImg   = newImageFill(ballWidth, ballHeight, color.RGBA{0, 0, 0, 0xFF})
	blockImg  = newImageFill(blockSize, blockSize, color.RGBA{0, 0, 0xFF, 0xFF})
	gameLayer = engine.Layer{
		Active: []bool{playMode: true},
		Sprites: []engine.Sprite{
			&paddle,
			&engine.RectSprite{ // left
				X: math.MinInt32, Y: 0,
				Width: -(math.MinInt32 + 1), Height: screenHeight,
			},
			&engine.RectSprite{ // right
				X: screenWidth, Y: 0,
				Width: math.MaxInt32 - screenWidth, Height: screenHeight,
			},
			&engine.RectSprite{ // top
				X: 0, Y: math.MinInt32,
				Width: screenWidth, Height: -(math.MinInt32 + 1),
			},
			&engine.RectSprite{ // bottom
				X: 0, Y: screenHeight,
				Width: screenWidth, Height: math.MaxInt32 - screenHeight,
			},
		},
	}
	level = engine.Level{
		Layers: []*engine.Layer{
			&engine.Layer{
				Sprites: []engine.Sprite{
					&engine.RectSprite{
						X: 0, Y: 0, Width: screenWidth, Height: screenHeight,
						Color: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
					},
				},
			},
			&gameLayer,
			&engine.Layer{
				Visible: []bool{pauseMode: true},
				Active:  []bool{pauseMode: true},
				Sprites: []engine.Sprite{
					&TextSprite{
						TextSprite: engine.TextSprite{
							Text:       "Paused\npress C to continue",
							Align:      engine.AlignCenter,
							Face:       face,
							Color:      color.RGBA{0, 0, 0, 0xFF},
							Margin:     10,
							Background: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
						},
						Centered: true,
					},
				},
			},
		},
	}
)

type PaddleSprite struct {
	engine.ImageSprite
}

func (sprt *PaddleSprite) Update(mode engine.Mode) bool {
	if start && ebiten.IsKeyPressed(ebiten.KeySpace) {
		start = false

		ball := BallSprite{
			ImageSprite: engine.ImageSprite{
				X:     sprt.X + (paddleWidth-ballWidth)/2,
				Y:     sprt.Y - ballHeight,
				Image: ballImg,
			},
		}
		ball.speed = 4
		ball.setDXDY((rand.Float64() * ball.speed / 2) - ball.speed)

		gameLayer.Sprites = append(gameLayer.Sprites, &ball)
	}

	cursorX, _ := ebiten.CursorPosition()
	sprt.X = float64(cursorX - paddleWidth/2)

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		sprt.DX -= paddleAccel
		if -sprt.DX > paddleMaxSpeed {
			sprt.DX = -paddleMaxSpeed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		sprt.DX += paddleAccel
		if sprt.DX > paddleMaxSpeed {
			sprt.DX = paddleMaxSpeed
		}
	} else if sprt.DX > paddleDecel {
		sprt.DX -= paddleDecel
	} else if sprt.DX < -paddleDecel {
		sprt.DX += paddleDecel
	} else {
		sprt.DX = 0
	}

	sprt.X += sprt.DX
	sprt.Y += sprt.DY
	return true // Always update to check for collisions.
}

func (sprt *PaddleSprite) Collision(with engine.Sprite) {
	if ball, ok := with.(*BallSprite); ok {
		if ball.DY > 0 {
			//ball.setDXDY(ball.DX + sprt.DX)

			sprtW, _ := sprt.Size()
			sprtX, _ := sprt.Corner()
			sprtX += sprtW / 2

			ballW, _ := ball.Size()
			ballX, _ := ball.Corner()
			ballX += ballW / 2

			angle := math.Min(math.Max((sprtX-ballX)*math.Pi/sprtW, -math.Pi/3), math.Pi/3)
			ball.setDXDY(math.Sin(-angle) * ball.speed)
		}
	} else if sprt.X < 0 {
		sprt.X = 0
		sprt.DX = -sprt.DX / 2
	} else if w, _ := sprt.Size(); sprt.X+w > screenWidth {
		sprt.X = screenWidth - w
		sprt.DX = -sprt.DX / 2
	}
}

func (sprt *PaddleSprite) Draw(mode engine.Mode, screen *ebiten.Image,
	op *ebiten.DrawImageOptions) {

	op.GeoM.Translate(sprt.X, sprt.Y)
	screen.DrawImage(sprt.Image, op)

	if start {
		op.GeoM.Translate((paddleWidth-ballWidth)/2, -ballHeight)
		screen.DrawImage(ballImg, op)
	}
}

type BallSprite struct {
	engine.ImageSprite
	speed float64
}

func (sprt *BallSprite) setDXDY(dx float64) {
	var dy float64
	if math.Abs(dx) < sprt.speed {
		dy = math.Sqrt(sprt.speed*sprt.speed - dx*dx)
	}
	if dy < sprt.speed/4 {
		dy = sprt.speed / 4
		if dx < 0 {
			dx = -math.Sqrt(sprt.speed*sprt.speed - dy*dy)
		} else {
			dx = math.Sqrt(sprt.speed*sprt.speed - dy*dy)
		}
	}
	sprt.DX = dx
	sprt.DY = -dy
}

func (sprt *BallSprite) Update(mode engine.Mode) bool {
	if sprt.DX == 0 && sprt.DY == 0 {
		return false
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		sprt.X += sprt.DX / 4
		sprt.Y += sprt.DY / 4
	} else {
		sprt.X += sprt.DX
		sprt.Y += sprt.DY
	}
	return true
}

func (sprt *BallSprite) Collision(with engine.Sprite) {
	if block, ok := with.(*BlockSprite); ok {
		block.collision(sprt)
	} else {
		if sprt.X < 0 {
			sprt.X = 0
			sprt.DX = -sprt.DX
		} else if sprt.Y < 0 {
			sprt.Y = 0
			sprt.DY = -sprt.DY
		} else {
			w, h := sprt.Size()
			if sprt.X+w > screenWidth {
				sprt.X = screenWidth - w
				sprt.DX = -sprt.DX
			} else if sprt.Y+h > screenHeight {
				start = true
				sprt.Delete()
			}
		}
	}
}

type BlockSprite struct {
	engine.ImageSprite
}

func (sprt *BlockSprite) collision(ball *BallSprite) {
	switch engine.Collide(ball, sprt) {
	case engine.CollideXGreater:
		ball.DX = math.Abs(ball.DX)
	case engine.CollideXLess:
		ball.DX = -math.Abs(ball.DX)
	case engine.CollideYGreater:
		ball.DY = math.Abs(ball.DY)
	case engine.CollideYLess:
		ball.DY = -math.Abs(ball.DY)
	}

	sprt.Delete()
}

type TextSprite struct {
	engine.TextSprite
	Centered bool
}

func (sprt *TextSprite) Update(mode engine.Mode) bool {
	sprt.TextSprite.Update(mode)
	if sprt.Centered {
		w, h := sprt.TextSprite.Size()
		sprt.X = (screenWidth - w) / 2
		sprt.Y = (screenHeight - h) / 2
	}

	return false
}

type breakout struct{}

type errorQuit struct{}

func (_ errorQuit) Error() string {
	return "error quit"
}

func (bo *breakout) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errQuit
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		mode = pauseMode
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		mode = playMode
	}

	level.Update(mode)
	return nil
}

func (bo *breakout) Draw(screen *ebiten.Image) {
	level.Draw(mode, screen)
}

func (bo *breakout) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func NewFace(fnt *opentype.Font, sz float64) font.Face {
	face, err := opentype.NewFace(fnt,
		&opentype.FaceOptions{
			Size:    sz,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "font face: %s\n", err)
		os.Exit(1)
	}
	return face
}

func newImageFill(w, h int, clr color.Color) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(clr)

	return img
}

func main() {
	cols := screenWidth / blockSize
	for cols*blockSize+(cols-1)*blockMargin+blockBorder*2 > screenWidth {
		cols -= 1
	}

	leftBorder := (screenWidth - (cols*blockSize + (cols-1)*blockMargin)) / 2
	for col := 0; col < cols; col += 1 {
		for row := 0; row*(blockSize+blockMargin) < screenHeight/2; row += 1 {
			gameLayer.Sprites = append(gameLayer.Sprites,
				&BlockSprite{
					ImageSprite: engine.ImageSprite{
						X:     float64(leftBorder + col*blockSize + (col-1)*blockMargin),
						Y:     float64(blockBorder + row*blockSize + (row-1)*blockMargin),
						Image: blockImg,
					},
				})
		}
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Breakout")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetCursorMode(ebiten.CursorModeHidden) // XXX: Show it in menu mode?

	err := ebiten.RunGame(&breakout{})
	if err != nil && !errors.Is(err, errQuit) {
		fmt.Fprintf(os.Stderr, "breakout failed: %s\n", err)
		os.Exit(1)
	}
}
