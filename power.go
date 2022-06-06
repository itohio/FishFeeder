package main

import (
	"machine"

	"tinygo.org/x/drivers/axp192"
	"tinygo.org/x/drivers/i2csoft"
)

func initPower() {
	i2c := i2csoft.New(machine.IO22, machine.IO21)
	i2c.Configure(i2csoft.I2CConfig{Frequency: 100e3})
	axp := axp192.New(i2c)
	axp.SetLDOEnable(2, true)
}
