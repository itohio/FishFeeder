package main

import (
	"machine"
	"time"

	srv "tinygo.org/x/drivers/servo"
)

type Servo struct {
	srv.Servo
	pin machine.Pin
}

func (s Servo) SetAngle(a float64) {
	b := 1000 + int16(a*1000/90)
	//s.SetMicroseconds(b)

	N := time.Microsecond * time.Duration(b)
	M := time.Millisecond*10 - N
	for i := 0; i < 50; i++ {
		s.pin.High()
		time.Sleep(N)
		s.pin.Low()
		time.Sleep(M)
	}
}

func newServo(pin machine.Pin) Servo {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	return Servo{
		pin: pin,
	}
}
