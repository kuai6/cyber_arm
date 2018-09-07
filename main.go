package main

import (
	"github.com/kuai6/cyber_arm/config"
	"github.com/kuai6/cyber_arm/server"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	host              string
	port              string
	messageBufferSize string
	messageQueueSize  string
)

func main() {
	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM)
	signal.Notify(stopSignal, syscall.SIGINT)

	var rootCmd = &cobra.Command{Use: "cyber-arm-service"}
	var start = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Run: func(cmd *cobra.Command, args []string) {
			configuration, err := readConfiguration()
			if err != nil {
				log.Fatalf("Failed to initialize server configuration: %s", err)
			}
			udpServer, messageChannel, err := server.NewServer(configuration)
			if err != nil {
				log.Fatalf("Failed to create server: %s", err)
			}

			go func() {
				udpServer.ReadMessage()
			}()

			for {
				select {
				case <-stopSignal:
					return
				case message := <-messageChannel:
					log.Printf("Received message: %s\n", string(message.Data))
				case <-time.After(1 * time.Second):
					log.Printf("Send sensor data")
					//addr, _ := net.ResolveUDPAddr("udp", "some_address")
					//udpServer.SendMessage(&server.Message{[]byte("sensor_data"), addr})
				}
			}
		},
	}

	start.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "server host")
	start.PersistentFlags().StringVar(&port, "port", "10001", "server port")
	start.PersistentFlags().StringVar(&messageBufferSize, "messageBufferSize", "1024", "message buffer size (in bytes)")
	start.PersistentFlags().StringVar(&messageQueueSize, "messageQueueSize", "10", "message queue size")

	rootCmd.AddCommand(start)
	go rootCmd.Execute()

	<-stopSignal
	log.Println("Server is shutting down...")
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
