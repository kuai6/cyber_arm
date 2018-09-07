package main

import (
	"github.com/kuai6/cyber_arm/config"
	"github.com/kuai6/cyber_arm/device"
	"github.com/kuai6/cyber_arm/server"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var (
	host              string
	port              string
	messageBufferSize string
	messageQueueSize  string
)

func main() {

	d := device.PCA9685{}
	d.Start()
	c := d.GetChannel(0)
	v, _ := strconv.ParseFloat(os.Args[1], 32)
	c.SetPercentage(float32(v))

	var rootCmd = &cobra.Command{Use: "cyber-arm-service"}
	var start = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Run: func(cmd *cobra.Command, args []string) {
			configuration, err := readConfiguration()
			if err != nil {
				log.Fatalf("Failed to initialize server configuration: %s", err)
			}
			server.Start(configuration)
		},
	}

	start.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "server host")
	start.PersistentFlags().StringVar(&port, "port", "10001", "server port")
	start.PersistentFlags().StringVar(&messageBufferSize, "messageBufferSize", "1024", "message buffer size (in bytes)")
	start.PersistentFlags().StringVar(&messageQueueSize, "messageQueueSize", "10", "message queue size")

	rootCmd.AddCommand(start)
	rootCmd.Execute()
}

func readConfiguration() (*config.ServerConfiguration, error) {
	configuration := new(config.ServerConfiguration)
	configuration.Host = host
	configuration.Port = port

	if messageBufferSize, err := strconv.Atoi(messageBufferSize); err != nil {
		return nil, err
	} else {
		configuration.MessageBufferSize = messageBufferSize
	}

	if messageQueueSize, err := strconv.Atoi(messageQueueSize); err != nil {
		return nil, err
	} else {
		configuration.MessageQueueSize = messageQueueSize
	}

	return configuration, nil
}
