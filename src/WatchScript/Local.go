package watcher

import (
	"github.com/howeyc/fsnotify"
	"log"
)

var directory = "E://temp//"

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(directory) // folder watch
	if err != nil {
		log.Fatal(err)
	}

	<-done
	watcher.Close()
}
