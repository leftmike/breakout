package main

import (
	"errors"
	"fmt"
	"image/color"
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

	paddleWidth    = 60
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

	background = engine.Sprite{}
	paddle     = engine.Sprite{
		X:         paddleX,
		Y:         paddleY,
		Collision: paddleCollision,
	}
	ball = engine.Sprite{
		Collision: ballCollision,
	}
	level = engine.Level{
		Layers: []*engine.Layer{
			&engine.Layer{
				Sprites: []*engine.Sprite{
					background.NewImageFill(screenWidth, screenHeight,
						color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}),
				},
			},
			&engine.Layer{
				Sprites: []*engine.Sprite{
					paddle.NewImageFill(paddleWidth, paddleHeight, color.RGBA{0, 0xFF, 0, 0xFF}),
					ball.NewImageFill(ballWidth, ballHeight, color.RGBA{0, 0, 0, 0xFF}),
					&engine.Sprite{X: -1, Y: 0, W: 0, H: screenHeight},          // left
					&engine.Sprite{X: screenWidth, Y: 0, W: 0, H: screenHeight}, // right
					&engine.Sprite{X: 0, Y: -1, W: screenWidth, H: 0},           // top
					&engine.Sprite{X: 0, Y: screenHeight, W: screenWidth, H: 0}, // bottom
				},
			},
		},
	}
)

func paddleCollision(sprt, with *engine.Sprite) {
	if with == &ball {
		ball.DY = -ball.DY
	} else if sprt.X < 0 {
		sprt.X = 0
		sprt.DX = -sprt.DX / 2
	} else if sprt.X+sprt.W > screenWidth {
		sprt.X = screenWidth - sprt.W
		sprt.DX = -sprt.DX / 2
	}
}

func initBall(sprt *engine.Sprite) {
	sprt.X = ballX
	sprt.Y = ballY
	sprt.DX = rand.Float64() * 4
	sprt.DY = 5 - sprt.DX
}

func ballCollision(sprt, with *engine.Sprite) {
	if sprt.X < 0 {
		sprt.X = 0
		sprt.DX = -sprt.DX
	} else if sprt.Y < 0 {
		sprt.Y = 0
		sprt.DY = -sprt.DY
	} else if sprt.X+sprt.W > screenWidth {
		sprt.X = screenWidth - sprt.W
		sprt.DX = -sprt.DX
	} else if sprt.Y+sprt.H > screenHeight {
		initBall(sprt)
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
	initBall(&ball)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Breakout")
	//ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(&breakout{})
	if err != nil && !errors.Is(err, errQuit) {
		fmt.Fprintf(os.Stderr, "breakout failed: %s\n", err)
		os.Exit(1)
	}
}
