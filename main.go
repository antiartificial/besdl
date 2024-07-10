package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

type Button struct {
	rect           sdl.Rect
	color          sdl.Color
	borderColor    sdl.Color
	hoverColor     sdl.Color
	isHovered      bool
	borderWidth    int32
	gradientOffset int32
}

func (b *Button) Draw(renderer *sdl.Renderer) {
	// Draw button fill
	renderer.SetDrawColor(b.color.R, b.color.G, b.color.B, b.color.A)
	renderer.FillRect(&b.rect)

	// Draw gradient border
	b.drawGradientBorder(renderer)

	// Draw button border when hovered
	if b.isHovered {
		b.drawHoverBorder(renderer)
	}
}

func (b *Button) drawGradientBorder(renderer *sdl.Renderer) {
	width, height := b.rect.W, b.rect.H
	for i := int32(0); i < b.borderWidth; i++ {
		ratio := float32(i) / float32(b.borderWidth)
		color := sdl.Color{
			R: uint8(255 * ratio),
			G: uint8(255 * (1 - ratio)),
			B: 0,
			A: 255,
		}
		renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		// Top border
		renderer.DrawLine(b.rect.X+i, b.rect.Y+i, b.rect.X+width-i-1, b.rect.Y+i)
		// Bottom border
		renderer.DrawLine(b.rect.X+i, b.rect.Y+height-i-1, b.rect.X+width-i-1, b.rect.Y+height-i-1)
		// Left border
		renderer.DrawLine(b.rect.X+i, b.rect.Y+i, b.rect.X+i, b.rect.Y+height-i-1)
		// Right border
		renderer.DrawLine(b.rect.X+width-i-1, b.rect.Y+i, b.rect.X+width-i-1, b.rect.Y+height-i-1)
	}
}

func (b *Button) drawHoverBorder(renderer *sdl.Renderer) {
	width, height := b.rect.W, b.rect.H
	for i := int32(0); i < b.borderWidth; i++ {
		color := b.hoverColor
		renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		renderer.DrawRect(&sdl.Rect{X: b.rect.X + i, Y: b.rect.Y + i, W: width - 2*i, H: height - 2*i})
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("could not initialize sdl: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("SDL Example", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("could not create renderer: %v", err)
	}
	defer renderer.Destroy()

	button := Button{
		rect:        sdl.Rect{X: 350, Y: 250, W: 100, H: 50},
		color:       sdl.Color{R: 0, G: 0, B: 255, A: 255},
		hoverColor:  sdl.Color{R: 255, G: 255, B: 0, A: 255},
		borderWidth: 5,
	}

	updateChan := make(chan int32)

	// Start a background goroutine to update the gradient offset
	go func() {
		for {
			time.Sleep(16 * time.Millisecond) // ~60 FPS
			updateChan <- button.gradientOffset + 1
		}
	}()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				x, y := t.X, t.Y
				button.isHovered = x >= button.rect.X && x <= (button.rect.X+button.rect.W) &&
					y >= button.rect.Y && y <= (button.rect.Y+button.rect.H)
			}
		}

		select {
		case newOffset := <-updateChan:
			button.gradientOffset = newOffset
		default:
		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		button.Draw(renderer)
		renderer.Present()

		sdl.Delay(16) // ~60 FPS
	}
}
