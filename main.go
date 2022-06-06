package main

// Icons:
// Fish Food icon by Icons8: https://icons8.com/icon/jQmKQCKAiKna/fish-food
// Aquarium icon by Icons8: https://icons8.com/icon/1Wi51AkN6dYv/aquarium
// Water Filter icon by Icons8: https://icons8.com/icon/dGlyR3EdBrXx/water-filter
// Thermometer icon by Icons8: https://icons8.com/icon/poFZHQZ-CjsC/thermometer

//go:generate tinygo build -target=m5stack -o fishfeeder.bin
//go:generate esptool -p COM3 write_flash -e 0x1000 fishfeeder.bin

import (
	"image/color"
	"machine"
	"time"
)

type Command int

const (
	NEXT = iota
	IDLE
	FEED
	DISPENSE
	RESET
	RESET_1
	RESET_2
	RESET_3
)

var (
	ledPin   = machine.IO10
	servoPin = machine.IO26
	btnMain  = machine.IO37
	btnSide  = machine.IO39

	display Display
	servo   Servo
	command = make(chan Command)
)

func initUART() {
	go func() {
		uart := machine.DefaultUART

		for {
			time.Sleep(time.Millisecond)
			if machine.DefaultUART.Bus.STATUS.Get()&255 < 1 {
				continue
			}
			b := machine.DefaultUART.Bus.RX_FIFO.Get()

			uart.WriteByte(b)

			switch b {
			case byte('f'):
				command <- FEED
				println("ok")
			case byte('d'):
				command <- DISPENSE
				println("ok")
			case byte('r'):
				command <- RESET
				println("ok")
			case byte('1'):
				command <- RESET_1
				println("ok")
			case byte('2'):
				command <- RESET_2
				println("ok")
			case byte('3'):
				command <- RESET_3
				println("ok")
			case byte('i'):
				println("ok ", len(nudges))
				for _, n := range nudges {
					println(n.ETA().String())
				}
			}
		}
	}()
}

func button(p machine.Pin) time.Duration {
	now := time.Now()
	time.Sleep(time.Millisecond * 10)
	for !p.Get() {
		time.Sleep(time.Millisecond)
	}
	return time.Since(now)
}

func initButtons() {
	btnMain.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnSide.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	go func() {
		for {
			if !btnMain.Get() {
				d := button(btnMain)
				if d > time.Second {
					command <- FEED
				} else {
					command <- RESET
				}
			}
			if !btnSide.Get() {
				_ = button(btnSide)
				command <- NEXT
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func echo() {
	for {
		time.Sleep(time.Millisecond)
		if machine.DefaultUART.Bus.STATUS.Get()&255 < 1 {
			continue
		}
		println(machine.DefaultUART.Bus.MEM_RX_STATUS.Get()&7, machine.DefaultUART.Bus.RXD_CNT.Get(), machine.DefaultUART.Bus.RX_FIFO.Get())
		// machine.DefaultUART.WriteByte(b)
	}
}

func main() {
	initPower()

	display = newDisplay()
	servo = newServo(servoPin)
	servo.SetAngle(88)

	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledPin.High()

	initUART()
	initButtons()

	display.FillScreen(color.RGBA{0, 0, 0, 255})

	go func() {
		for {
			command <- IDLE
			time.Sleep(time.Second)
		}
	}()

	for cmd := range command {
		for i := range nudges {
			nudges[i].Check()
		}

		handleCommand(cmd)
		draw()
	}
}

func handleCommand(c Command) {
	switch c {
	case IDLE:
		handleFeeding()
		return
	case FEED:
		feed()
	case DISPENSE:
		dispenseFood()
	case NEXT:
		selectedNudge = (selectedNudge + 1) % len(nudges)
		selectedTime = time.Now()
	case RESET:
		if selectedNudge >= 0 {
			nudges[selectedNudge].Reset()
		}
	case RESET_1:
		nudges[0].Reset()
	case RESET_2:
		nudges[1].Reset()
	case RESET_3:
		nudges[2].Reset()
	}
	println("Command: ", c)
}

var lastSelectedNudge = -1

func draw() {
	if time.Since(selectedTime).Seconds() > 30 {
		selectedNudge = -1
	}
	if lastSelectedNudge != selectedNudge {
		display.FillScreen(color.RGBA{})
		lastSelectedNudge = selectedNudge
	}

	for i, n := range nudges {
		x := 1 + int16(53*i)
		if selectedNudge == i {
			display.FillRectangle(x-1, 9, 42, 42, color.RGBA{R: 100})
		}
		display.DrawRGBBitmap(x, 10, icons[i], 40, 40)
		d := int16(40 * n.ETAPercent())
		display.Bar(x, 10+41, d, 40, maxColor(colors[i]))
	}
	display.DrawRGBBitmap(1, 80-ThermometerHeight-1, ThermometerPng, ThermometerWidth, ThermometerHeight)
}
