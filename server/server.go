package server

import (
	"net"
	"log"
	"os"
	"os/signal"
	"syscall"
	"cyber_arm/config"
)

type MessageChannel chan string

func Start(config *config.ServerConfiguration) {
	connection, err := net.ListenPacket("udp", config.Address())
	if err != nil {
		log.Fatalf("Failed to start UDP server: %s", err)
	}
	defer connection.Close()

	log.Printf("Started UDP server: %s\n", connection.LocalAddr().String())

	messageChannel := listenPackets(connection, config.MessageBufferSize, config.MessageQueueSize)
	stopSignalChannel := listenInterruption()

	for {
		select {
		case message := <-messageChannel:
			{
				log.Printf("Received message: %s\n", message)
			}
		case <-stopSignalChannel:
			{
				log.Println("Server is shutting down...")
				return
			}
		}
	}
}

func listenPackets(connection net.PacketConn, bufferSize int, queueSize int) MessageChannel {
	buffer := make([]byte, bufferSize)
	channel := make(MessageChannel, queueSize)

	go func() {
		for {
			n, _, err := connection.ReadFrom(buffer)
			if err != nil {
				log.Printf("Failed to receive UDP packet: %s\n", err)
			}
			message := string(buffer[0:n])
			channel <- message
		}
	}()

	return channel
}

func listenInterruption() chan os.Signal {
	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM)
	signal.Notify(stopSignal, syscall.SIGINT)

	return stopSignal
}
