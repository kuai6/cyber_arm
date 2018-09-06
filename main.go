package main

import (
	"github.com/spf13/cobra"
	"strconv"
	"log"
)

var (
	host string
	port string
	messageBufferSize string
	messageQueueSize string
)

func main() {
	var rootCmd = &cobra.Command{Use: "cyber-arm-service"}
	var start = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := readConfig()
			if err != nil {
				log.Fatalf("Failed to initialize server configuration: %s", err)
			}
			StartServer(config)
		},
	}

	start.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "server host")
	start.PersistentFlags().StringVar(&port, "port", "10001", "server port")
	start.PersistentFlags().StringVar(&messageBufferSize, "messageBufferSize",  "1024", "message buffer size (in bytes)")
	start.PersistentFlags().StringVar(&messageQueueSize, "messageQueueSize", "10", "message queue size")

	rootCmd.AddCommand(start)
	rootCmd.Execute()
}

func readConfig() (*ServerConfiguration, error) {
	config := new(ServerConfiguration)
	config.Host = host
	config.Port = port

	if messageBufferSize, err := strconv.Atoi(messageBufferSize); err != nil {
		return nil, err
	} else {
		config.MessageBufferSize = messageBufferSize
	}

	if messageQueueSize, err := strconv.Atoi(messageQueueSize); err != nil {
		return nil, err
	} else {
		config.MessageQueueSize = messageQueueSize
	}

	return config, nil
}
