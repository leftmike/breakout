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
)

var (
	start = true

	errQuit = errorQuit{}

	background = engine.NewImageSprite(0, 0,
		engine.NewImageFill(screenWidth, screenHeight, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}))
	paddle = PaddleSprite{
		ImageSprite: engine.NewImageSprite(paddleX, paddleY,
			engine.NewImageFill(paddleWidth, paddleHeight, color.RGBA{0, 0xFF, 0, 0xFF})),
	}
	ballImg   = engine.NewImageFill(ballWidth, ballHeight, color.RGBA{0, 0, 0, 0xFF})
	gameLayer = engine.Layer{
		Sprites: []engine.Sprite{
			&paddle,
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
	}
	level = engine.Level{
		Layers: []*engine.Layer{
			&engine.Layer{
				Sprites: []engine.Sprite{
					&background,
				},
			},
			&gameLayer,
		},
	}
)

type PaddleSprite struct {
	engine.ImageSprite
}

func (sprt *PaddleSprite) Update() bool {
	if start && ebiten.IsKeyPressed(ebiten.KeySpace) {
		start = false

		ball := BallSprite{
			ImageSprite: engine.NewImageSprite(0, 0, ballImg),
		}
		ball.speed = 4
		ball.X = sprt.X + (paddleWidth-ballWidth)/2
		ball.Y = sprt.Y - ballHeight
		ball.setDXDY((rand.Float64() * ball.speed / 2) - ball.speed)

		gameLayer.Sprites = append(gameLayer.Sprites, &ball)
	}

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
			ball.setDXDY(ball.DX + sprt.DX)
		}
	} else if sprt.X < 0 {
		sprt.X = 0
		sprt.DX = -sprt.DX / 2
	} else if sprt.X+sprt.Width > screenWidth {
		sprt.X = screenWidth - sprt.Width
		sprt.DX = -sprt.DX / 2
	}
}

func (sprt *PaddleSprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
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
		start = true
		sprt.Delete()
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
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Breakout")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	err := ebiten.RunGame(&breakout{})
	if err != nil && !errors.Is(err, errQuit) {
		fmt.Fprintf(os.Stderr, "breakout failed: %s\n", err)
		os.Exit(1)
	}
}
