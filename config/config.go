package config

import "fmt"

type ServerConfiguration struct {
	Host              string
	Port              string
	MessageBufferSize int
	MessageQueueSize  int
}

type BusConfiguration struct {
	I2cAddress string
}

func (c *ServerConfiguration) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (b *BusConfiguration) Address() string {
	return b.I2cAddress
}
