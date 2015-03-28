package blog

import (
	"github.com/howeyc/fsnotify"
	"log"
)

func watch(dir string, callback func(e *fsnotify.FileEvent, err error)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				callback(ev, nil)
			case err := <-watcher.Error:
				callback(nil, err)
				done <- true
			}
		}
	}()
	if err := watcher.Watch(dir); err != nil {
		log.Fatal(err)
	}
	<-done
}
