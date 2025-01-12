package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// watchFile наблюдает за изменениями файла filePath, и отправляет новые данные,
// записанные в файл, в канал newData.
func watchFile(filePath string, newData chan<- []byte) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Println(err.Error())
		return
	}
	//nolint:all
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Set initial read position to the end of file
	readPosition := stat.Size()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err.Error())
		return
	}
	//nolint:all
	defer watcher.Close()

	err = watcher.Add(watchedFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println(event)

			stat, err := file.Stat()
			if err != nil {
				log.Println(err.Error())
				return
			}

			// Only read from file if it's size has increased
			size := stat.Size()
			if size > readPosition {
				// Make buffer of size just enough to read all new content
				buf := make([]byte, size-readPosition)
				_, err = file.ReadAt(buf, readPosition)
				if err != nil && err.Error() != "EOF" {
					log.Printf("error reading from log file: %s", err.Error())
					return
				}

				// Send data via channel
				newData <- buf
			}

			// Again, set initial read position to the end of file
			readPosition = size
		case errors := <-watcher.Errors:
			log.Println(errors)
		}
	}
}
