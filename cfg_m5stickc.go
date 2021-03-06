package main

import (
	"machine"

	"tinygo.org/x/drivers/axp192"
	"tinygo.org/x/drivers/i2csoft"
)

var (
	ledPin     = machine.IO10
	servoPin   = machine.IO26
	btnMain    = machine.IO37
	btnSide    = machine.IO39
	displaySCK = machine.IO13
	displaySDO = machine.IO15
	displayRST = machine.IO18
	displayDC  = machine.IO23
	displayCS  = machine.IO5
	displayBL  = machine.NoPin

	// Grove connector
	i2cCLK  = machine.IO33
	i2cDATA = machine.IO32

	// I2C bus used for internal periphery
	i2cInternalCLK  = machine.IO22
	i2cInternalDATA = machine.IO21
)

// some M5StickCs have LCD backlight off by default, so we must turn on the LCD ldo.
func initPower() {
	i2c := i2csoft.New(i2cInternalCLK, i2cInternalDATA)
	i2c.Configure(i2csoft.I2CConfig{Frequency: 100e3})
	axp := axp192.New(i2c)
	axp.SetLDOEnable(2, true)
}
