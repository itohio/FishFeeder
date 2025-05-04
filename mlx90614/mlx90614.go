package mlx90614

import (
	"math"
	"time"

	"tinygo.org/x/drivers"
)

const (
	I2C_ADDR     = 0x5A
	REG_TAMBIENT = 0x06
	REG_TOBJECT1 = 0x07
	REG_TOBJECT2 = 0x08

	// EEPROM
	REG_EMISSIVITY = 0x24
)

// Device wraps an SPI connection.
type Device struct {
	bus  drivers.I2C
	addr uint8
}

func New(bus drivers.I2C, addr uint8) *Device {
	return &Device{
		bus:  bus,
		addr: addr,
	}
}

func (d *Device) Configure(emissivity float64) error {
	e := uint32(emissivity * 0xFFFF)
	if e > 0xFFFF {
		e = 0xFFFF
	}
	d.write16(REG_EMISSIVITY, 0)
	time.Sleep(time.Millisecond * 10)
	d.write16(REG_EMISSIVITY, uint16(e))
	return nil
}

func (d *Device) ID() [4]uint16 {
	var id [4]uint16
	for i := range id {
		id[i] = d.read16(0x3C + uint8(i))
	}
	return id
}

func (d *Device) Temperature(reg uint8) float64 {
	temp := d.read16(reg)
	if temp == 0 {
		return math.NaN()
	}

	return float64(temp)*.02 - 273.15
}

func (d *Device) Ambient() float64 {
	return d.Temperature(REG_TAMBIENT)
}

func (d *Device) Object1() float64 {
	return d.Temperature(REG_TOBJECT1)
}

func (d *Device) Object2() float64 {
	return d.Temperature(REG_TOBJECT2)
}

func (d *Device) Emissivity() float64 {
	return float64(d.read16(REG_EMISSIVITY)) / 0xFFFF
}

func crc8(d []uint8) uint8 {
	var crc uint8
	for _, inbyte := range d {
		for i := 0; i < 8; i++ {
			carry := (crc ^ inbyte) & 0x80
			crc <<= 1
			if carry != 0 {
				crc ^= 0x7
			}
			inbyte <<= 1
		}
	}
	return crc
}

func (d *Device) write16(reg uint8, data uint16) error {
	db := [4]byte{d.addr, reg, byte(data & 0xFF), byte((data & 0xFF00) >> 8)}
	pec := crc8(db[:])
	db[0] = db[1]
	db[1] = db[2]
	db[2] = db[3]
	db[3] = pec
	// return d.bus.WriteRegister(d.addr, reg, db[2:])
	return d.bus.Tx(uint16(d.addr), db[1:], nil)
}

func (d *Device) read16(reg uint8) uint16 {
	db := [3]byte{}
	err := d.bus.Tx(uint16(d.addr), []byte{reg}, db[:])
	// err := d.bus.ReadRegister(d.addr, reg, db[:])
	if err != nil {
		return 0
	}
	return (uint16(db[1]) << 8) | uint16(db[0])
}
