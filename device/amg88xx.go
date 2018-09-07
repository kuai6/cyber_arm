package device

import (
	"golang.org/x/exp/io/i2c"
	"sync"
	"time"
)

const (
	AMG88xx_PCTL byte = 0x00
	AMG88xx_RST  byte = 0x01
	AMG88xx_FPSC byte = 0x02
	AMG88xx_INTC byte = 0x03
	AMG88xx_STAT byte = 0x04
	AMG88xx_SCLR byte = 0x05
	//0x06 reserved
	AMG88xx_AVE          byte = 0x07
	AMG88xx_INTHL        byte = 0x08
	AMG88xx_INTHH        byte = 0x09
	AMG88xx_INTLL        byte = 0x0A
	AMG88xx_INTLH        byte = 0x0B
	AMG88xx_IHYSL        byte = 0x0C
	AMG88xx_IHYSH        byte = 0x0D
	AMG88xx_TTHL         byte = 0x0E
	AMG88xx_TTHH         byte = 0x0F
	AMG88xx_INT_OFFSET   byte = 0x010
	AMG88xx_PIXEL_OFFSET byte = 0x80

	// power
	AMG88xx_NORMAL_MODE byte = 0x00
	AMG88xx_SLEEP_MODE  byte = 0x01
	AMG88xx_STAND_BY_60 byte = 0x20
	AMG88xx_STAND_BY_10 byte = 0x21

	//reset
	AMG88xx_FLAG_RESET    byte = 0x30
	AMG88xx_INITIAL_RESET byte = 0x3F

	// frame rate
	AMG88xx_FPS_10 byte = 0x00
	AMG88xx_FPS_1  byte = 0x01

	// interrupt
	AMG88xx_INT_DISABLED byte = 0x00
	AMG88xx_INT_ENABLED  byte = 0x01

	// init modes
	AMG88xx_DIFFERENCE     byte = 0x00
	AMG88xx_ABSOLUTE_VALUE byte = 0x01

	// MISC

	AMG88xx_PIXEL_ARRAY_SIZE      int     = 64
	AMG88xx_PIXEL_TEMP_CONVERSION float64 = 0.25
	AMG88xx_THERMISTOR_CONVERSION float64 = 0.0625
)

// power control register
type pctl struct {
	// 0x00 = Normal Mode
	// 0x01 = Sleep Mode
	// 0x20 = Stand-by mode (60 sec intermittence)
	// 0x21 = Stand-by mode (10 sec intermittence)
	PCTL uint8
}

func (r *pctl) get() uint8 {
	return r.PCTL
}

//reset register
type rst struct {
	//0x30 = flag reset (all clear status reg 0x04, interrupt flag and interrupt table)
	//0x3F = initial reset (brings flag reset and returns to initial setting)
	RST uint8
}

func (r *rst) get() uint8 {
	return r.RST
}

//frame rate register
type fpsc struct {
	//0 = 10FPS
	//1 = 1FPS
	FPS uint8
}

func (r *fpsc) get() uint8 {
	return r.FPS & 0x01
}

//interrupt control register
type intc struct {
	// 0 = INT output reactive (Hi-Z)
	// 1 = INT output active
	INTEN uint8

	// 0 = Difference interrupt mode
	// 1 = absolute value interrupt mode
	INTMOD uint8
}

func (r *intc) get() byte {
	return r.INTMOD<<1 | r.INTEN&0x03
}

//status register
type stat struct {
	//interrupt outbreak (val of interrupt table reg)
	INTF uint8

	//temperature output overflow (val of temperature reg)
	OVF_IRS uint8

	//thermistor temperature output overflow (value of thermistor)
	OVF_THS uint8
}

func (r *stat) get() uint8 {
	return ((r.OVF_THS << 3) | (r.OVF_IRS << 2) | (r.INTF << 1)) & 0x07
}

//status clear register
//write to clear overflow flag and interrupt flag
//after writing automatically turns to 0x00
type sclr struct {
	//interrupt flag clear
	INTCLR uint8
	//temp output overflow flag clear
	OVS_CLR uint8
	//thermistor temp output overflow flag clear
	OVT_CLR uint8
}

func (r *sclr) get() uint8 {
	return ((r.OVT_CLR << 3) | (r.OVS_CLR << 2) | (r.INTCLR << 1)) & 0x07
}

//average register
//for setting moving average output mode
type ave struct {
	//1 = twice moving average mode
	MAMOD uint8
}

func (r *ave) get() uint8 {
	return r.MAMOD << 5
}

//interrupt level registers
//for setting upper / lower limit hysteresis on interrupt level

//interrupt level upper limit setting. Interrupt output
// and interrupt pixel table are set when value exceeds set value
type inthl struct {
	INT_LVL_H uint8
}

func (r *inthl) get() uint8 {
	return r.INT_LVL_H
}

type inthh struct {
	INT_LVL_H uint8
}

func (r *inthh) get() uint8 {
	return r.INT_LVL_H
}

//interrupt level lower limit. Interrupt output
//and interrupt pixel table are set when value is lower than set value
type intll struct {
	INT_LVL_L uint8
}

func (r *intll) get() uint8 {
	return r.INT_LVL_L
}

type intlh struct {
	INT_LVL_L uint8
}

func (r *intlh) get() uint8 {
	return r.INT_LVL_L & 0xF
}

//setting of interrupt hysteresis level when interrupt is generated.
//should not be higher than interrupt level
type ihysl struct {
	INT_HYS uint8
}

func (r *ihysl) get() uint8 {
	return r.INT_HYS
}

type ihysh struct {
	INT_HYS uint8
}

func (r *ihysh) get() uint8 {
	return r.INT_HYS & 0xF
}

//thermistor register
//SIGNED MAGNITUDE FORMAT
type tthl struct {
	TEMP uint8
}

func (r *tthl) get() uint8 {
	return r.TEMP
}

type tthh struct {
	TEMP uint8
	SIGN uint8
}

func (r *tthh) get() uint8 {
	return ((r.SIGN << 3) | r.TEMP) & 0xF
}

type AMG88XX struct {
	d  *i2c.Device
	mu sync.RWMutex
}

// Start function initialize device
func (a *AMG88XX) Start() error {
	// Use i2creg I²C bus registry to find the first available I²C bus.
	a.mu.Lock()
	defer a.mu.Unlock()
	// Use i2creg I²C bus registry to find the first available I²C bus.
	var err error
	a.d, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x69)
	if err != nil {
		return err
	}
	// noraml mode
	_pctl := pctl{PCTL: AMG88xx_NORMAL_MODE}
	err = a.writeByte(AMG88xx_PCTL, _pctl.get())
	if err != nil {
		return err
	}
	//software reset
	_rst := rst{RST: AMG88xx_INITIAL_RESET}
	err = a.writeByte(AMG88xx_RST, _rst.get())
	if err != nil {
		return err
	}
	// disabling interrupt
	_intc := intc{INTEN: 0}
	err = a.writeByte(AMG88xx_INTC, _intc.get())
	if err != nil {
		return err
	}
	//set to 10 FPS
	_fpsc := fpsc{FPS: AMG88xx_FPS_10}
	err = a.writeByte(AMG88xx_FPSC, _fpsc.get())
	if err != nil {
		return err
	}
	//delay
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (p *AMG88XX) writeByte(reg byte, byteVal byte) error {
	return p.d.WriteReg(reg, []byte{byteVal})
}

func (p *AMG88XX) ReadPixelsRAW() []byte {
	buf := make([]byte, 64)
	p.d.ReadReg(AMG88xx_PIXEL_OFFSET, buf)
	return buf
}

func (p *AMG88XX) ReadPixels() []float64 {
	raw := p.ReadPixelsRAW()
	buf := make([]float64, 64)

	for i, b := range raw {
		recast := uint16(b<<8) | uint16(b)
		buf[i] = p.signedMag12ToFloat(recast) * AMG88xx_PIXEL_TEMP_CONVERSION
	}
	return buf
}

func (p *AMG88XX) signedMag12ToFloat(val uint16) float64 {
	//take first 11 bits as absolute val
	absVal := (val & 0x7FF)

	if val&0x8000 == 1 {
		return 0 - float64(absVal)
	}
	return float64(absVal)
}

// Stop - unregister device resources
func (a *AMG88XX) Stop() {
	a.d.Close()
}
