package main

import (
	"github.com/kuai6/cyber_arm/command"
	"github.com/kuai6/cyber_arm/device"
	s "github.com/kuai6/cyber_arm/server"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

//var (
//	host              string
//	port              string
//	messageBufferSize string
//	messageQueueSize  string
//)

func main() {
	cyberArmAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:10001")
	if err != nil {
		log.Fatal(err)
	}

	pca9685 := new(device.PCA9685)
	err = pca9685.Start()
	if err != nil {
		panic(err)
	}

	//amg88xx := device.AMG88XX{}
	//err = amg88xx.Start()
	//if err != nil {
	//    panic(err)
	//}
	//fmt.Printf("%v", d.ReadPixels())

	s.ListenCyberArmCommands(cyberArmAddr, func(command *command.Command) {
		switch command.Name {
		case "ROTATE":
			alpha, err := strconv.ParseFloat(command.Arguments[0], 32)
			if err != nil {
				log.Printf("Failed to parse argument: %s", err)
			}
			beta, err := strconv.ParseFloat(command.Arguments[1], 32)
			if err != nil {
				log.Printf("Failed to parse argument: %s", err)
			}
			log.Printf("Perform cyber-arm rotation to (%f,%f)\n", alpha, beta)
			//rotate(alpha, beta)
			xChannel := pca9685.GetChannel(0)
			xChannel.SetPercentage(float32(alpha))

			yChannel := pca9685.GetChannel(1)
			yChannel.SetPercentage(float32(beta))

		case "FIRE":
			log.Printf("Perform fire action\n")
			//fire()
		}
	})
	s.ConnectServer(cyberArmAddr, []byte(`{"name":"ROTATE","arguments":["50", "100"]}`))
	//s.ConnectServer(cyberArmAddr, []byte(`{"name":"FIRE"}`))

	thermalSensorAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:10002")
	if err != nil {
		log.Fatal(err)
	}

	s.StreamThermalSensorData(thermalSensorAddr)
	//s.ConnectServer(thermalSensorAddr, nil)

	//var server, thermalServer *s.Server
	//
	//var rootCmd = &cobra.Command{Use: "cyber-arm-service"}
	//var start = &cobra.Command{
	//	Use:   "start",
	//	Short: "Start server",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		serverConfig, err := readConfiguration()
	//		if err != nil {
	//			log.Fatalf("Failed to initialize server configuration: %s", err)
	//		}
	//
	//		server, err = s.NewServer(serverConfig)
	//		if err != nil {
	//			log.Fatalf("Failed to create server: %s", err)
	//		}
	//
	//		err = server.Start(func(message *s.Message) error {
	//			log.Printf("Received message: %s\n", string(message.Data))
	//			return nil
	//		})
	//		if err != nil {
	//			log.Fatalf("Failed to start server: %s", err)
	//		}
	//	},
	//}
	//
	//start.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "server host")
	//start.PersistentFlags().StringVar(&port, "port", "10001", "server port")
	//start.PersistentFlags().StringVar(&messageBufferSize, "messageBufferSize", "1024", "message buffer size (in bytes)")
	//start.PersistentFlags().StringVar(&messageQueueSize, "messageQueueSize", "10", "message queue size")
	//
	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM)
	signal.Notify(stopSignal, syscall.SIGINT)

	//rootCmd.AddCommand(start)
	//go rootCmd.Execute()
	//
	<-stopSignal
	log.Println("Server is shutting down...")
}

//func readConfiguration() (*config.ServerConfiguration, error) {
//	configuration := new(config.ServerConfiguration)
//	configuration.Host = host
//	configuration.Port = port
//
//	if messageBufferSize, err := strconv.Atoi(messageBufferSize); err != nil {
//		return nil, err
//	} else {
//		configuration.MessageBufferSize = messageBufferSize
//	}
//
//	if messageQueueSize, err := strconv.Atoi(messageQueueSize); err != nil {
//		return nil, err
//	} else {
//		configuration.MessageQueueSize = messageQueueSize
//	}
//
//	return configuration, nil
//}
