package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mister-mqtt/internal/mqtt"
	"mister-mqtt/internal/watcher"
)

var (
	brokerAddr  = flag.String("broker", "localhost:1883", "MQTT broker address")
	clientID    = flag.String("client-id", "mister-mqtt", "MQTT client ID")
	username    = flag.String("username", "", "MQTT username")
	password    = flag.String("password", "", "MQTT password")
	topicPrefix = flag.String("topic-prefix", "homeassistant", "MQTT topic prefix")
)

func main() {
	flag.Parse()

	fmt.Printf("mister-mqtt starting...\n")
	fmt.Printf("MQTT Broker: %s\n", *brokerAddr)
	fmt.Printf("Client ID: %s\n", *clientID)
	fmt.Printf("Topic Prefix: %s\n", *topicPrefix)

	// Initialize MQTT client
	mqttConfig := mqtt.Config{
		BrokerAddr:  *brokerAddr,
		ClientID:    *clientID,
		Username:    *username,
		Password:    *password,
		TopicPrefix: *topicPrefix,
	}

	mqttClient, err := mqtt.NewClient(mqttConfig)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}

	if err := mqttClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()

	// Set up HomeAssistant discovery
	if err := mqttClient.PublishDiscovery(); err != nil {
		log.Fatalf("Failed to publish discovery messages: %v", err)
	}

	// Initialize file watcher
	fileWatcher, err := watcher.NewFileWatcher(func(filename, content string) {
		content = strings.TrimSpace(content)
		var sensorID string
		
		switch filename {
		case "CORENAME":
			sensorID = "corename"
		case "ACTIVEGAME":
			sensorID = "activegame"
		case "RBFNAME":
			sensorID = "rbfname"
		default:
			log.Printf("Unknown file: %s", filename)
			return
		}

		if err := mqttClient.PublishState(sensorID, content); err != nil {
			log.Printf("Failed to publish state for %s: %v", sensorID, err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}

	if err := fileWatcher.Start(); err != nil {
		log.Fatalf("Failed to start file watcher: %v", err)
	}
	defer fileWatcher.Stop()

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("mister-mqtt is running. Press Ctrl+C to exit.")
	<-c

	fmt.Println("\nShutting down mister-mqtt...")
}
