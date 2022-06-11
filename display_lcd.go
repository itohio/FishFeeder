package main

import (
	"fmt"
	"image/color"
	"machine"
	"strings"
	"time"

	"github.com/itohio/FishFeeder/st7735"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/image/png"
)

var display Display

type Display struct {
	st7735.Device
}

var _ drivers.Displayer = &Display{}

func newDisplay() Display {
	spi := machine.SPI3
	spi.Configure(machine.SPIConfig{
		Frequency: 32000000,
		SCK:       displaySCK,
		SDO:       displaySDO,
	})
	time.Sleep(time.Second)
	display := st7735.New(spi, displayRST, displayDC, displayCS, displayBL)
	display.Configure(st7735.Config{
		Model:        st7735.MINI80x160,
		Width:        80,
		Height:       160,
		ColumnOffset: 26,
		RowOffset:    1,
		Rotation:     st7735.ROTATION_270,
	})

	return Display{
		Device: display,
	}
}

func (d *Display) Bar(x, y, i, w int16, c color.RGBA) {
	if i < 0 {
		i = 0
	}
	if i > w {
		i = w
	}
	if i != 0 {
		d.FillRectangle(x, y, i, 7, c)
	}
	if i != w {
		d.FillRectangle(x+i, y, w-i, 7, color.RGBA{20, 20, 20, 255})
	}
}

var buffer [3 * 256]uint16

// NOTE: This part does not work with tinygo version 0.23.0 windows/amd64 (using go version go1.18 and LLVM version 14.0.0)
// The effect is that some memory overwrite occurs (serios BUG that should be isolated) and panics happen.
// However, it is mostly not needed due to the fact, that raw bytes of the icons are smaller than embedded PNG files.
func (d *Display) DrawPng(x0, y0 int16, pngImage string) error {
	p := strings.NewReader(pngImage)
	png.SetCallback(buffer[:], func(data []uint16, x, y, w, h, width, height int16) {
		W, H := d.Size()
		println(x0+x, y0+y, w, h, width, height, W, H, len(data))
		err := d.DrawRGBBitmap(x0+x, y0+y, data[:w*h], w, h)
		if err != nil {
			println(fmt.Errorf("error drawPng: %s", err))
		}
	})

	println("decode")
	w, h := d.Size()
	println(w, h)
	_, err := png.Decode(p)
	println("decode done")
	w, h = d.Size()
	println(w, h)
	if err != nil {
		println(fmt.Errorf("error drawPng: %s", err))
	}
	return err
}
