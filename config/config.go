package config

import "fmt"

type ServerConfiguration struct {
	Host string
	Port string
	MessageBufferSize int
	MessageQueueSize int
}

func (c *ServerConfiguration) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
