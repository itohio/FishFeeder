package main

import (
	"image/color"
	"time"
)

var (
	feedingTime   = time.Now()
	feedingCount  = 0
	feedingDelays = []time.Duration{
		time.Second * 30,
		time.Second * 30,
		time.Second * 30,
		time.Second * 3,
	}
)

func handleFeeding() {
	if feedingCount >= len(feedingDelays) {
		feedingCount = 0
	}
	if feedingCount == 0 {
		return
	}

	var (
		sum    time.Duration
		cur    time.Duration
		curTmp time.Duration
	)
	for i, d := range feedingDelays {
		if i == len(feedingDelays)-1 {
			break
		}
		sum += d
		if i+1 < feedingCount {
			cur += d
		}
	}
	cur += time.Since(feedingTime)

	for i, d := range feedingDelays {
		if i == len(feedingDelays)-1 {
			break
		}
		curTmp += d
		w := int16(160 * (curTmp.Seconds() / sum.Seconds()))
		display.FillRectangle(w, 0, 1, 10, color.RGBA{G: 40})
	}

	w := int16(160 * (cur.Seconds() / sum.Seconds()))
	display.FillRectangle(0, 0, w, 7, color.RGBA{G: 70})

	if time.Since(feedingTime) > feedingDelays[feedingCount-1] {
		feedingTime = time.Now()
		go func() {
			command <- FEED
		}()
	}
}

func feed() {
	feedingTime = time.Now()
	feedingCount += 1

	dispenseFood()

	if feedingCount >= len(feedingDelays) {
		go func() {
			command <- RESET_1
		}()
		feedingCount = 0
	}
}

func dispenseFood() {
	display.FillScreen(color.RGBA{})
	display.DrawRGBBitmap((160-40)/2, (80-40)/2, FoodPng, 40, 40)

	servo.SetAngle(5)
	time.Sleep(time.Millisecond * 200)
	servo.SetAngle(92)
	display.FillScreen(color.RGBA{})
}
