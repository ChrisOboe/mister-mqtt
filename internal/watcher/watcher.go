package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const (
	CoreNameFile   = "/tmp/CORENAME"
	ActiveGameFile = "/tmp/ACTIVEGAME"
	RBFNameFile    = "/tmp/RBFNAME"
)

// FileWatcher monitors MiSTer status files for changes
type FileWatcher struct {
	watcher   *fsnotify.Watcher
	onChange  func(filename, content string)
	files     []string
}

// NewFileWatcher creates a new file watcher instance
func NewFileWatcher(onChange func(filename, content string)) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	fw := &FileWatcher{
		watcher:  watcher,
		onChange: onChange,
		files:    []string{CoreNameFile, ActiveGameFile, RBFNameFile},
	}

	return fw, nil
}

// Start begins watching the files
func (fw *FileWatcher) Start() error {
	// Add files to watcher
	for _, file := range fw.files {
		// Create file if it doesn't exist
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if err := fw.createFile(file); err != nil {
				log.Printf("Warning: could not create %s: %v", file, err)
			}
		}

		// Add to watcher
		if err := fw.watcher.Add(file); err != nil {
			log.Printf("Warning: could not watch %s: %v", file, err)
		}
	}

	// Start watching in a goroutine
	go fw.watch()

	// Read initial values
	for _, file := range fw.files {
		if content, err := fw.readFile(file); err == nil {
			fw.onChange(filepath.Base(file), content)
		}
	}

	return nil
}

// Stop stops the file watcher
func (fw *FileWatcher) Stop() error {
	return fw.watcher.Close()
}

func (fw *FileWatcher) watch() {
	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				content, err := fw.readFile(event.Name)
				if err != nil {
					log.Printf("Error reading file %s: %v", event.Name, err)
					continue
				}
				fw.onChange(filepath.Base(event.Name), content)
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (fw *FileWatcher) readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (fw *FileWatcher) createFile(filename string) error {
	return os.WriteFile(filename, []byte(""), 0644)
}
