package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"strconv"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	// window
	cfg := pixelgl.WindowConfig{
		Title:  "Starting",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	// sprite
	pic, err := loadPicture("hiking.png")
	if err != nil {
		panic(err)
	}
	dizzyGohper := pixel.NewSprite(pic, pic.Bounds())
	angleDizzyGohper := 0.0

	// txt dash
	txtDash := text.New(pixel.V(100, 500), text.NewAtlas(basicfont.Face7x13, text.ASCII))
	lastSPFChecked := time.Now()
	lastTick := time.Now()
	framesSinceLastFPSCheck := 0
	lastFPS := 0

	// channel - timer
	chTick := time.Tick(time.Second)

	// event loop
	for !win.Closed() {
		// txt dash
		spf := time.Since(lastSPFChecked).Seconds()
		lastSPFChecked = time.Now()
		framesSinceLastFPSCheck++
		select {
		case tick := <-chTick:
			lastTick = tick
			lastFPS = framesSinceLastFPSCheck
			framesSinceLastFPSCheck = 0
		default:
		}
		txtDash.Clear()
		fmt.Fprintln(txtDash, "delta time - seconds per frame: ", spf)
		fmt.Fprintln(txtDash, "tick: ", lastTick)
		fmt.Fprintln(txtDash, "framerate: ", lastFPS)

		// sprite
		angleDizzyGohper -= math.Pi * 2 * spf // single rotation per sec
		matDizzyGohper := pixel.IM.Scaled(pixel.ZV, 0.5)
		matDizzyGohper = matDizzyGohper.Rotated(pixel.ZV, angleDizzyGohper)
		matDizzyGohper = matDizzyGohper.Moved(win.Bounds().Center())

		// update
		win.SetTitle("SPF: " + strconv.FormatFloat(spf, 'f', 7, 64) +
			" // " + lastTick.String() +
			" // " + " Framerate: " + strconv.Itoa(lastFPS),
		)
		win.Clear(colornames.Darkorange)
		dizzyGohper.Draw(win, matDizzyGohper)
		txtDash.Draw(win, pixel.IM.Scaled(txtDash.Orig, 2))

		win.Update()
	} // for
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
