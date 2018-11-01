package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// Watcher 监听配置文件变化
func watcher(cfgChan chan *Config, path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("配置文件发生变化")
					cfg, err := readConfFile(path)
					if err != nil {
						return
					}
					cfgChan <- cfg
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
