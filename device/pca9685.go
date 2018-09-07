package device

import (
	"errors"
	"fmt"
	"golang.org/x/exp/io/i2c"
	"math"
	"time"
)

const (
	MODE1         byte = 0x00
	MODE2         byte = 0x01
	PRESCALE      byte = 0xFE
	LED0_ON_L     byte = 0x06
	LED0_ON_H     byte = 0x07
	LED0_OFF_L    byte = 0x08
	LED0_OFF_H    byte = 0x09
	ALL_LED_ON_L  byte = 0xFA
	ALL_LED_ON_H  byte = 0xFB
	ALL_LED_OFF_L byte = 0xFC
	ALL_LED_OFF_H byte = 0xFD
	OUTDRV        byte = 0x04
	ALLCALL       byte = 0x01
	SLEEP         byte = 0x10
	BYTE          byte = 0xFF

	DEFAULT_FREQ float32 = 100.0
	OSC_FREQ     float32 = 25000000.0
	STEP_COUNT   float32 = 4096.0
	MIN_PULSE    float32 = 100
	MAX_PULSE    float32 = 480
)

type PCA9685 struct {
	d *i2c.Device
}

type Channel struct {
	p   *PCA9685
	pin int
}

func (p *PCA9685) GetChannel(pin int) *Channel {
	return &Channel{
		p:   p,
		pin: pin,
	}
}

func (c *Channel) SetPulse(on int, off int) error {
	if on < 0 || on > off || off > 4096 {
		return errors.New(fmt.Sprintf(
			"On/Off (%d/%d) must be between 0 and %d",
			on,
			off,
			STEP_COUNT,
		))
	}

	c.p.setPwm(c.pin, on, off)

	return nil
}

func (c *Channel) SetPercentage(percent float32) {
	pulseLength := int((MAX_PULSE-MIN_PULSE)*percent/100 + MIN_PULSE)

	c.SetPulse(0, pulseLength)
}

func (p *PCA9685) Start() error {
	// Use i2creg I²C bus registry to find the first available I²C bus.
	var err error
	p.d, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x40)
	if err != nil {
		return err
	}
	p.SetAllPwm(0, 0)
	p.writeByte(MODE2, OUTDRV)
	p.writeByte(MODE1, ALLCALL)
	time.Sleep(5 * time.Millisecond)
	mode := p.readByte(MODE1)
	mode = mode & ^SLEEP
	p.writeByte(MODE1, mode)
	time.Sleep(5 * time.Millisecond)
	p.setPwmFrequency(DEFAULT_FREQ)

	return nil
}

// Stop - unregister device resources
func (p *PCA9685) Stop() {
	p.d.Close()
}

func (p *PCA9685) setPwmFrequency(freqHz float32) {

	preScaleValue := OSC_FREQ // 25MHz
	preScaleValue /= STEP_COUNT
	preScaleValue /= freqHz
	preScaleValue -= 1.0

	preScale := int(math.Floor(float64(preScaleValue + 0.5)))

	oldMode := p.read8(MODE1)

	newMode := (oldMode & 0x7F) | 0x10

	p.write8(MODE1, newMode)
	p.write8(PRESCALE, preScale)
	p.write8(MODE1, oldMode)

	time.Sleep(5 * time.Millisecond)

	p.write8(MODE1, oldMode|0x80)
}

func (p *PCA9685) setServoPulse(pwmNumber int, pulse int) {

	var pulseLength float32 = 1000000

	pulseLength /= float32(60)

	pulseLength /= STEP_COUNT

	pulseF := float32(pulse)

	pulseF /= pulseLength

	p.setPwm(pwmNumber, 0, int(pulseF))
}

func (p *PCA9685) SetAllPwm(on int, off int) {
	p.write8(ALL_LED_ON_L, on)
	p.write8(ALL_LED_ON_H, on>>8)
	p.write8(ALL_LED_OFF_L, off)
	p.write8(ALL_LED_OFF_H, off>>8)
}

func (p *PCA9685) setPwm(pwm int, on int, off int) {
	p.write8(LED0_ON_L+byte(4)*byte(pwm), on)
	p.write8(LED0_ON_H+byte(4)*byte(pwm), on>>8)
	p.write8(LED0_OFF_L+byte(4)*byte(pwm), off)
	p.write8(LED0_OFF_H+byte(4)*byte(pwm), off>>8)
}

func (p *PCA9685) write8(reg byte, intVal int) {
	byteVal := byte(intVal) & BYTE

	p.writeByte(reg, byteVal)
}

func (p *PCA9685) writeByte(reg byte, byteVal byte) {
	p.d.WriteReg(reg, []byte{byteVal})
}

func (p *PCA9685) read8(reg byte) int {
	byteVal := p.readByte(reg)

	return int(byteVal)
}

func (p *PCA9685) readByte(reg byte) byte {
	buf := make([]byte, 1)
	p.d.ReadReg(reg, buf)
	return buf[0]
}
