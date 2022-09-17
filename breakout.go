package main

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leftmike/breakout/engine"
)

const (
	windowWidth  = 500
	windowHeight = 500
	screenWidth  = 400
	screenHeight = 400

	paddleWidth    = 80
	paddleHeight   = 10
	paddleY        = screenHeight - (paddleHeight + 2)
	paddleAccel    = 0.5
	paddleDecel    = 1
	paddleMaxSpeed = 20

	ballWidth  = 10
	ballHeight = 10
)

var (
	paddleImg   *ebiten.Image
	paddleX     float64 = (screenWidth - paddleWidth) / 2
	paddleSpeed float64 = 0

	ballImg *ebiten.Image
	ballX   float64 = (screenWidth - ballWidth) / 2
	ballY   float64 = 0

	errQuit = errorQuit{}

	background = engine.NewImageSpriteFill(0, 0, screenWidth, screenHeight,
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	paddle = PaddleSprite{
		ImageSprite: engine.NewImageSpriteFill(paddleX, paddleY, paddleWidth, paddleHeight,
			color.RGBA{0, 0xFF, 0, 0xFF}),
	}
	ball = BallSprite{
		ImageSprite: engine.NewImageSpriteFill(0, 0, ballWidth, ballHeight,
			color.RGBA{0, 0, 0, 0xFF}),
	}
	level = engine.Level{
		Layers: []*engine.Layer{
			&engine.Layer{
				Sprites: []engine.Sprite{
					&background,
				},
			},
			&engine.Layer{
				Sprites: []engine.Sprite{
					&paddle,
					&ball,
					&engine.ImageSprite{ // left
						X: -1, Y: 0,
						Width: 0, Height: screenHeight,
					},
					&engine.ImageSprite{ // right
						X: screenWidth, Y: 0,
						Width: 0, Height: screenHeight,
					},
					&engine.ImageSprite{ // top
						X: 0, Y: -1,
						Width: screenWidth, Height: 0,
					},
					&engine.ImageSprite{ // bottom
						X: 0, Y: screenHeight,
						Width: screenWidth, Height: 0,
					},
				},
			},
		},
	}
)

type PaddleSprite struct {
	engine.ImageSprite
}

func (sprt *PaddleSprite) Update() bool {
	// Always update to check for collisions.
	sprt.X += sprt.DX
	sprt.Y += sprt.DY
	return true
}

func (sprt *PaddleSprite) Collision(with engine.Sprite) {
	if ball, ok := with.(*BallSprite); ok {
		if ball.DY > 0 {
			ball.setDXDY(ball.DX + paddle.DX)
			ball.DY = -ball.DY
		}
	} else if sprt.X < 0 {
		sprt.X = 0
		sprt.DX = -sprt.DX / 2
	} else if sprt.X+sprt.Width > screenWidth {
		sprt.X = screenWidth - sprt.Width
		sprt.DX = -sprt.DX / 2
	}
}

type BallSprite struct {
	engine.ImageSprite
	speed float64
}

func (sprt *BallSprite) init() {
	sprt.speed = 4
	sprt.X = ballX
	sprt.Y = ballY
	sprt.setDXDY((rand.Float64() * sprt.speed / 2) - sprt.speed)
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
	sprt.DY = dy
}

func (sprt *BallSprite) Collision(with engine.Sprite) {
	if sprt.X < 0 {
		sprt.X = 0
		sprt.DX = -sprt.DX
	} else if sprt.Y < 0 {
		sprt.Y = 0
		sprt.DY = -sprt.DY
	} else if sprt.X+sprt.Width > screenWidth {
		sprt.X = screenWidth - sprt.Width
		sprt.DX = -sprt.DX
	} else if sprt.Y+sprt.Height > screenHeight {
		sprt.init()
	}
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

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		paddle.DX -= paddleAccel
		if -paddle.DX > paddleMaxSpeed {
			paddle.DX = -paddleMaxSpeed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		paddle.DX += paddleAccel
		if paddle.DX > paddleMaxSpeed {
			paddle.DX = paddleMaxSpeed
		}
	} else if paddle.DX > paddleDecel {
		paddle.DX -= paddleDecel
	} else if paddle.DX < -paddleDecel {
		paddle.DX += paddleDecel
	} else {
		paddle.DX = 0
	}

	level.Update()
	return nil
}

func (bo *breakout) Draw(screen *ebiten.Image) {
	level.Draw(screen)
}

func (bo *breakout) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ball.init()
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Breakout")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(&breakout{})
	if err != nil && !errors.Is(err, errQuit) {
		fmt.Fprintf(os.Stderr, "breakout failed: %s\n", err)
		os.Exit(1)
	}
}
