package main

import (
	"image/color"
	"time"

	"github.com/itohio/FishFeeder/icons"
)

var (
	selectedTime  = time.Now()
	selectedNudge = -1
	iconImages    = [][]uint16{
		icons.FoodPng,
		icons.AquariumPng,
		icons.FilterPng,
	}

	colors = []color.RGBA{
		{G: 8},
		{B: 8},
		{R: 8},
	}

	nudges = []Nudge{
		{
			timestamp: time.Now(),
			delay:     time.Hour * 12,
			nudge:     makeNudgeWithCount(colors[0], icons.FoodPng, 2),
		},
		{
			timestamp: time.Now(),
			delay:     time.Hour * 24 * 14,
			nudge:     makeNudge(colors[1], icons.AquariumPng),
		},
		{
			timestamp: time.Now(),
			delay:     time.Hour * 24 * 31,
			nudge:     makeNudge(colors[2], icons.FilterPng),
		},
	}
)

func maxColor(c color.RGBA) color.RGBA {
	m := uint8(255)
	r := 0.0
	if 255-c.R < m {
		m = 255 - c.R
		r = 255 / float64(c.R)
	}
	if 255-c.G < m {
		m = 255 - c.G
		r = 255 / float64(c.G)
	}
	if 255-c.B < m {
		m = 255 - c.B
		r = 255 / float64(c.B)
	}
	return color.RGBA{
		R: uint8(float64(c.R) * r),
		G: uint8(float64(c.G) * r),
		B: uint8(float64(c.B) * r),
	}
}

func fill(c color.RGBA) {
	display.FillScreen(c)
}

func flash(c color.RGBA, n *Nudge) bool {
	for j := 0; j < 2; j++ {
		time.Sleep(time.Millisecond * 333)
		for i := 0; i < 3; i++ {
			fill(maxColor(c))
			ledPin.Low()
			time.Sleep(time.Millisecond * 100)
			fill(color.RGBA{R: 0x16, G: 0x16, B: 0x16})
			ledPin.High()
			time.Sleep(time.Millisecond * 100)

			if clicked() {
				return true
			}
		}
		fill(color.RGBA{})
	}

	return false
}

func clicked() bool {
	if cmd, ok := <-command; ok && (cmd >= RESET || cmd == FEED) {
		if cmd == FEED {
			go func() {
				command <- FEED
			}()
		}
		return true
	}
	return false
}

func sleep(n *Nudge, icon []uint16) bool {
	for i := 0; i < 30; i++ {
		fill(color.RGBA{})
		display.DrawRGBBitmap((160-40)/2, (80-40)/2+(int16(i%2))*5, icon, 40, 40)
		time.Sleep(time.Millisecond * 100)
		if clicked() {
			return true
		}
	}
	return false
}

func makeNudge(c color.RGBA, icon []uint16) func(*Nudge) bool {
	return makeNudgeWithCount(c, icon, -1)
}

func makeNudgeWithCount(c color.RGBA, icon []uint16, N int) func(*Nudge) bool {
	return func(n *Nudge) bool {
		for i := 0; i < N || N == -1; i++ {
			if flash(c, n) || sleep(n, icon) {
				display.FillScreen(color.RGBA{})
				return true
			}
		}
		return false
	}
}
