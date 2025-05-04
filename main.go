package main

//go:generate tinygo build -target=m5stick-c -o fishfeeder.bin
//go:generate esptool -p COM5 write_flash -e 0x1000 fishfeeder.bin

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	"github.com/itohio/FishFeeder/i2c"
	"github.com/itohio/FishFeeder/icons"
	"github.com/itohio/FishFeeder/mlx90614"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
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

var temperature *mlx90614.Device
var command = make(chan Command)

func initUART() {
	go func() {
		uart := machine.DefaultUART

		for {
			b, _ := uart.ReadByte()
			// time.Sleep(time.Millisecond)
			// if machine.DefaultUART.Bus.STATUS.Get()&255 < 1 {
			// 	continue
			// }
			// b := machine.DefaultUART.Bus.RX_FIFO.Get()

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
			case byte('t'):
				println("ok 5 ")
				id := temperature.ID()
				println(id[0], id[1], id[2], id[3])
				println(temperature.Emissivity())
				println(temperature.Ambient())
				println(temperature.Object1())
				println(temperature.Object2())
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
				d := button(btnSide)
				if d > time.Second {
					command <- DISPENSE
				} else {
					command <- NEXT
				}
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func initTemperature() {
	bus := i2c.New(i2cCLK, i2cDATA)
	bus.Configure(i2c.I2CConfig{Frequency: 20e3})

	for i := 0; i < 255; i++ {
		if bus.Tx(uint16(i), []byte{}, []byte{0}) == nil {
			println(i)
		}
	}

	temperature = mlx90614.New(bus, mlx90614.I2C_ADDR)
	temperature.Configure(.90) // emissivity of water
}

func echo() {
	for {
		time.Sleep(time.Millisecond)
		if machine.DefaultUART.Bus.STATUS.Get()&255 < 1 {
			continue
		}
		// println(machine.DefaultUART.Bus.MEM_RX_STATUS.Get()&7, machine.DefaultUART.Bus.RXD_CNT.Get(), machine.DefaultUART.Bus.RX_FIFO.Get())
		// machine.DefaultUART.WriteByte(b)
	}
}

func main() {
	initPower()
	initTemperature()

	display = newDisplay()
	servo = newServo(servoPin)
	servo.SetAngle(88)

	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledPin.High()

	// initUART()
	initButtons()

	display.FillScreen(color.RGBA{0, 0, 0, 255})

	go func() {
		for {
			command <- IDLE
			time.Sleep(time.Second)
		}
	}()

	originalNudge := nudges[0].nudge
	nudges[0].nudge = func(n *Nudge) bool {
		if originalNudge(n) {
			return true
		}

		go func() {
			command <- FEED
		}()
		return true
	}

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
		display.DrawRGBBitmap(x, 10, iconImages[i], 40, 40)
		d := int16(40 * n.ETAPercent())
		display.Bar(x, 10+41, d, 40, maxColor(colors[i]))
	}

	display.DrawRGBBitmap(1, 80-icons.ThermometerHeight-1, icons.ThermometerPng, icons.ThermometerWidth, icons.ThermometerHeight)

	display.FillRectangle(20, 80-icons.ThermometerHeight, 140, 14, color.RGBA{})
	ambient := strconv.FormatFloat(temperature.Ambient(), 'f', 1, 32)
	// ambient := fmt.Sprintf("%0.1f", temperature.Ambient())
	water := strconv.FormatFloat(temperature.Object1(), 'f', 1, 32)
	// water := fmt.Sprintf("%0.1f", temperature.Object1())
	tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 21, 80-icons.ThermometerHeight-1+12, water, color.RGBA{100, 100, 255, 255})
	tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 21+60, 80-icons.ThermometerHeight-1+12, ambient, color.RGBA{100, 100, 100, 255})
}
