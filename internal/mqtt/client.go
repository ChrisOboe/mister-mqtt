package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
)

// Client wraps the MQTT client functionality
type Client struct {
	client      mqtt.Client
	topicPrefix string
	nodeID      string
}

// Config holds MQTT client configuration
type Config struct {
	BrokerAddr  string
	ClientID    string
	Username    string
	Password    string
	TopicPrefix string
}

// HomeAssistantDiscovery represents a HomeAssistant MQTT discovery message
type HomeAssistantDiscovery struct {
	Name              string            `json:"name"`
	StateTopic        string            `json:"state_topic"`
	UniqueID          string            `json:"unique_id"`
	Device            Device            `json:"device"`
	AvailabilityTopic string            `json:"availability_topic"`
	PayloadAvailable  string            `json:"payload_available"`
	PayloadNotAvailable string          `json:"payload_not_available"`
}

// Device represents the device information for HomeAssistant discovery
type Device struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Model        string   `json:"model"`
	Manufacturer string   `json:"manufacturer"`
	SWVersion    string   `json:"sw_version"`
}

// NewClient creates a new MQTT client
func NewClient(config Config) (*Client, error) {
	// Read hostname from /etc/hostname
	nodeID, err := getNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to get node ID from hostname: %w", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", config.BrokerAddr))
	opts.SetClientID(config.ClientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	})

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	client := mqtt.NewClient(opts)
	
	c := &Client{
		client:      client,
		topicPrefix: config.TopicPrefix,
		nodeID:      nodeID,
	}

	return c, nil
}

// getNodeID reads the hostname from /etc/hostname and returns it as the node ID
func getNodeID() (string, error) {
	content, err := os.ReadFile("/etc/hostname")
	if err != nil {
		return "", fmt.Errorf("failed to read /etc/hostname: %w", err)
	}
	
	hostname := strings.TrimSpace(string(content))
	if hostname == "" {
		return "", fmt.Errorf("hostname is empty")
	}
	
	return hostname, nil
}

// Connect establishes connection to MQTT broker
func (c *Client) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	// Publish availability
	availTopic := fmt.Sprintf("%s/sensor/%s/availability", c.topicPrefix, c.nodeID)
	if token := c.client.Publish(availTopic, 0, true, "online"); token.Wait() && token.Error() != nil {
		log.Printf("Failed to publish availability: %v", token.Error())
	}

	return nil
}

// Disconnect closes the MQTT connection
func (c *Client) Disconnect() {
	// Publish offline status
	availTopic := fmt.Sprintf("%s/sensor/%s/availability", c.topicPrefix, c.nodeID)
	if token := c.client.Publish(availTopic, 0, true, "offline"); token.Wait() && token.Error() != nil {
		log.Printf("Failed to publish offline status: %v", token.Error())
	}

	c.client.Disconnect(250)
}

// PublishDiscovery publishes HomeAssistant discovery messages
func (c *Client) PublishDiscovery() error {
	sensors := []struct {
		ID   string
		Name string
	}{
		{"corename", "MiSTer Core Name"},
		{"activegame", "MiSTer Active Game"},
		{"rbfname", "MiSTer RBF Name"},
	}

	device := Device{
		Identifiers:  []string{c.nodeID},
		Name:         "MiSTer FPGA",
		Model:        "MiSTer",
		Manufacturer: "MiSTer Project",
		SWVersion:    "1.0.0",
	}

	for _, sensor := range sensors {
		discovery := HomeAssistantDiscovery{
			Name:                fmt.Sprintf("%s %s", device.Name, sensor.Name),
			StateTopic:          fmt.Sprintf("%s/sensor/%s/%s/state", c.topicPrefix, c.nodeID, sensor.ID),
			UniqueID:           fmt.Sprintf("%s_%s", c.nodeID, sensor.ID),
			Device:             device,
			AvailabilityTopic:  fmt.Sprintf("%s/sensor/%s/availability", c.topicPrefix, c.nodeID),
			PayloadAvailable:   "online",
			PayloadNotAvailable: "offline",
		}

		payload, err := json.Marshal(discovery)
		if err != nil {
			return fmt.Errorf("failed to marshal discovery message for %s: %w", sensor.ID, err)
		}

		topic := fmt.Sprintf("%s/sensor/%s/%s/config", c.topicPrefix, c.nodeID, sensor.ID)
		if token := c.client.Publish(topic, 0, true, payload); token.Wait() && token.Error() != nil {
			return fmt.Errorf("failed to publish discovery for %s: %w", sensor.ID, token.Error())
		}

		log.Printf("Published discovery for %s", sensor.Name)
	}

	return nil
}

// PublishState publishes sensor state
func (c *Client) PublishState(sensorID, state string) error {
	topic := fmt.Sprintf("%s/sensor/%s/%s/state", c.topicPrefix, c.nodeID, sensorID)
	
	if token := c.client.Publish(topic, 0, false, state); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish state for %s: %w", sensorID, token.Error())
	}

	log.Printf("Published state for %s: %s", sensorID, state)
	return nil
}
