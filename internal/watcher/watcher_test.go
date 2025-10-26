package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNewFileWatcher(t *testing.T) {
	called := false
	onChange := func(filename, content string) {
		called = true
	}

	fw, err := NewFileWatcher(onChange)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fw == nil {
		t.Fatal("Expected FileWatcher instance, got nil")
	}

	if fw.onChange == nil {
		t.Fatal("Expected onChange callback to be set")
	}

	if len(fw.files) != 3 {
		t.Errorf("Expected 3 files to watch, got %d", len(fw.files))
	}

	// Test callback works
	fw.onChange("test", "content")
	if !called {
		t.Error("Expected onChange callback to be called")
	}

	fw.Stop()
}

func TestFileWatcherWithTempFiles(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Override the file paths for testing
	testFile := filepath.Join(tempDir, "TESTFILE")
	
	// Create a custom watcher with a test file
	fw := &FileWatcher{
		files: []string{testFile},
	}

	var err error
	fw.watcher, err = createTestWatcher()
	if err != nil {
		t.Fatalf("Failed to create test watcher: %v", err)
	}
	defer fw.Stop()

	// Track onChange calls
	var lastFilename, lastContent string
	fw.onChange = func(filename, content string) {
		lastFilename = filename
		lastContent = content
	}

	// Create the test file
	initialContent := "initial content"
	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Start watching
	if err := fw.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	// Give some time for initial read
	time.Sleep(100 * time.Millisecond)

	// Check initial read
	if lastFilename != "TESTFILE" {
		t.Errorf("Expected filename 'TESTFILE', got '%s'", lastFilename)
	}
	if lastContent != initialContent {
		t.Errorf("Expected content '%s', got '%s'", initialContent, lastContent)
	}

	// Update the file
	newContent := "updated content"
	if err := os.WriteFile(testFile, []byte(newContent), 0644); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	// Give some time for the watcher to detect changes
	time.Sleep(200 * time.Millisecond)

	// Check if the change was detected
	if lastContent != newContent {
		t.Errorf("Expected updated content '%s', got '%s'", newContent, lastContent)
	}
}

// createTestWatcher creates a watcher for testing purposes
func createTestWatcher() (*fsnotify.Watcher, error) {
	return fsnotify.NewWatcher()
}

func TestReadFile(t *testing.T) {
	fw := &FileWatcher{}
	
	// Create a temporary file
	tempFile := filepath.Join(t.TempDir(), "testfile")
	expectedContent := "test content"
	
	if err := os.WriteFile(tempFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	content, err := fw.readFile(tempFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, content)
	}
}

func TestCreateFile(t *testing.T) {
	fw := &FileWatcher{}
	
	tempFile := filepath.Join(t.TempDir(), "newfile")
	
	if err := fw.createFile(tempFile); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}

	// Check if file is empty
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty file, got content: %s", string(content))
	}
}
