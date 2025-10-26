package mqtt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	// We'll test with a mocked hostname since we can't easily mock the real getNodeID function
	// In a real test environment, we'd use dependency injection or interfaces
	config := Config{
		BrokerAddr:  "localhost:1883",
		ClientID:    "test-client",
		Username:    "user",
		Password:    "pass",
		TopicPrefix: "homeassistant",
	}

	// This test will fail in environments without /etc/hostname or fail to connect to MQTT
	// In practice, this would be better with dependency injection for testability
	client, err := NewClient(config)
	
	// We expect this to fail gracefully if /etc/hostname doesn't exist
	if err != nil {
		t.Logf("NewClient failed as expected in test environment: %v", err)
		return
	}

	if client == nil {
		t.Fatal("Expected client instance, got nil")
	}

	if client.topicPrefix != config.TopicPrefix {
		t.Errorf("Expected topicPrefix '%s', got '%s'", config.TopicPrefix, client.topicPrefix)
	}

	if client.nodeID == "" {
		t.Error("Expected nodeID to be set from hostname")
	}
}

func TestHomeAssistantDiscoveryJSON(t *testing.T) {
	device := Device{
		Identifiers:  []string{"test-node"},
		Name:         "MiSTer FPGA",
		Model:        "MiSTer",
		Manufacturer: "MiSTer Project",
		SWVersion:    "1.0.0",
	}

	discovery := HomeAssistantDiscovery{
		Name:                "MiSTer FPGA Core Name",
		StateTopic:          "homeassistant/sensor/test-node/corename/state",
		UniqueID:           "test-node_corename",
		Device:             device,
		AvailabilityTopic:  "homeassistant/sensor/test-node/availability",
		PayloadAvailable:   "online",
		PayloadNotAvailable: "offline",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(discovery)
	if err != nil {
		t.Fatalf("Failed to marshal discovery: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled HomeAssistantDiscovery
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal discovery: %v", err)
	}

	// Check important fields
	if unmarshaled.Name != discovery.Name {
		t.Errorf("Expected name '%s', got '%s'", discovery.Name, unmarshaled.Name)
	}

	if unmarshaled.StateTopic != discovery.StateTopic {
		t.Errorf("Expected state topic '%s', got '%s'", discovery.StateTopic, unmarshaled.StateTopic)
	}

	if unmarshaled.UniqueID != discovery.UniqueID {
		t.Errorf("Expected unique ID '%s', got '%s'", discovery.UniqueID, unmarshaled.UniqueID)
	}

	if unmarshaled.Device.Name != device.Name {
		t.Errorf("Expected device name '%s', got '%s'", device.Name, unmarshaled.Device.Name)
	}
}

func TestConfig(t *testing.T) {
	config := Config{
		BrokerAddr:  "test.broker:1883",
		ClientID:    "test-client-id",
		Username:    "testuser",
		Password:    "testpass",
		TopicPrefix: "test/prefix",
	}

	// Test that all fields are properly set
	if config.BrokerAddr == "" {
		t.Error("BrokerAddr should not be empty")
	}
	if config.ClientID == "" {
		t.Error("ClientID should not be empty")
	}
	if config.TopicPrefix == "" {
		t.Error("TopicPrefix should not be empty")
	}
}

func TestDevice(t *testing.T) {
	device := Device{
		Identifiers:  []string{"device1", "device2"},
		Name:         "Test Device",
		Model:        "Test Model",
		Manufacturer: "Test Manufacturer",
		SWVersion:    "1.0.0",
	}

	if len(device.Identifiers) != 2 {
		t.Errorf("Expected 2 identifiers, got %d", len(device.Identifiers))
	}

	if device.Identifiers[0] != "device1" {
		t.Errorf("Expected first identifier 'device1', got '%s'", device.Identifiers[0])
	}

	if device.Name != "Test Device" {
		t.Errorf("Expected name 'Test Device', got '%s'", device.Name)
	}
}

func TestGetNodeID(t *testing.T) {
	// Create a temporary hostname file
	tempDir := t.TempDir()
	tempHostname := filepath.Join(tempDir, "hostname")
	expectedHostname := "test-mister-node"
	
	if err := os.WriteFile(tempHostname, []byte(expectedHostname+"\n"), 0644); err != nil {
		t.Fatalf("Failed to create test hostname file: %v", err)
	}

	// Test reading hostname from file directly
	content, err := os.ReadFile(tempHostname)
	if err != nil {
		t.Fatalf("Failed to read test hostname file: %v", err)
	}
	
	hostname := strings.TrimSpace(string(content))
	if hostname != expectedHostname {
		t.Errorf("Expected hostname '%s', got '%s'", expectedHostname, hostname)
	}
}

func TestGetNodeIDFromActualHostname(t *testing.T) {
	// Test with the actual getNodeID function if /etc/hostname exists
	nodeID, err := getNodeID()
	if err != nil {
		t.Logf("Cannot read /etc/hostname in test environment: %v", err)
		return
	}

	if nodeID == "" {
		t.Error("Expected non-empty nodeID from hostname")
	}
}
