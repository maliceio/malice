package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/maliceio/malice/config"
)

func cmdWatch(folderName string, logs bool) error {

	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Info("Malice watching folder: ", folderName)

	info, err := os.Stat(folderName)

	// Check that folder exists
	if os.IsNotExist(err) {
		log.Error("error: folder does not exist.")
		return nil
	}
	// Check that path is a folder and not a file
	if info.IsDir() {
		NewWatcher(folderName)
	} else {
		log.Error("error: path is not a folder")
	}

	return nil
}

// NewWatcher creates a new watcher for the user supplied folder
func NewWatcher(folder string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("modified file:", event.Name)
					// Scan new sample in watch folder
					ScanSample(event.Name)
				}
			case err := <-watcher.Errors:
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
