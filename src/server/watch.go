package server

import (
	"github.com/howeyc/fsnotify"
	"log"
	"math/rand"
	"path/filepath"
	"time"
)

var curBuildTask int // rebuild() will invoked repeatedly

func rebuildArticles(task int) {
	if task != curBuildTask {
		return
	}
	log.Println("rebuild task:", task)
	blogBuilder.RebuildAll()
}

func watch(dir string) {
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
				log.Println("event:", eventString(ev))
				curBuildTask = rand.Intn(10000)
				go func(taskid int) {
					time.Sleep(time.Second * 2)
					rebuildArticles(taskid)
				}(curBuildTask)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(dir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
	/* ... do stuff ... */
	watcher.Close()
}

func eventString(e *fsnotify.FileEvent) string {
	_, file := filepath.Split(e.Name)
	var events string = ""
	if e.IsCreate() {
		events += "|" + "CREATE"
	}
	if e.IsDelete() {
		events += "|" + "DELETE"
	}
	if e.IsModify() {
		events += "|" + "MODIFY"
	}
	if e.IsRename() {
		events += "|" + "RENAME"
	}
	if e.IsAttrib() {
		events += "|" + "ATTRIB"
	}
	if len(events) > 0 {
		events = events[1:]
	}
	return file + " " + events
}
