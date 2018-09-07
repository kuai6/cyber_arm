package device

import (
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

type AMG8831 struct {
	d *i2c.Dev
	b i2c.BusCloser
}

// Start function initialize device
func (ts *AMG8831) Start() error {
	// Use i2creg I²C bus registry to find the first available I²C bus.
	var err error
	ts.b, err = i2creg.Open("AMG8831")
	if err != nil {
		return err
	}
	defer ts.b.Close()

	// Dev is a valid conn.Conn.
	ts.d = &i2c.Dev{Addr: 23, Bus: ts.b}
	return nil
}

// Stop - unregister device resources
func (ts *AMG8831) Stop() {
	ts.b.Close()
}
