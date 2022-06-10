package main

import (
	"machine"
	"time"
)

var servo Servo

type Servo machine.Pin

func (s Servo) SetAngle(a float64) {
	b := 1000 + int16(a*1000/90)

	N := time.Microsecond * time.Duration(b)
	M := time.Millisecond*10 - N
	for i := 0; i < 50; i++ {
		machine.Pin(s).High()
		time.Sleep(N)
		machine.Pin(s).Low()
		time.Sleep(M)
	}
}

func newServo(pin machine.Pin) Servo {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	return Servo(pin)
}
