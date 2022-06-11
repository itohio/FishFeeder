//go:build esp32
// +build esp32

package i2c

import (
	"device"
)

// wait waits for half the time of the SCL operation interval. It is set to
// about 50 kHz.
func (i2c *I2C) wait() {
	wait := 50
	for i := 0; i < wait; i++ {
		device.Asm(`nop`)
	}
}
