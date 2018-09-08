package main

import (
	s "github.com/kuai6/cyber_arm/server"
	"log"
	"net"
	"os"
	"os/signal"
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

	s.ListenCyberArmCommands(cyberArmAddr)
	//s.ConnectServer(cyberArmAddr, []byte(`{"name":"ROTATE","arguments":["1.1", "1.2"]}`))
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
