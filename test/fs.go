package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"github.com/fsnotify/fsnotify"
)

func getCurrentBrightness(path string) int {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading brightness: %v", err)
		return -1
	}
	value, err := strconv.Atoi(string(data[:len(data)-1]))
	if err != nil {
		log.Printf("Error parsing brightness: %v", err)
		return -1
	}
	return value
}

func Fswatch() {
	backlightPathGlob := "/sys/class/backlight/*/brightness"
	paths, err := filepath.Glob(backlightPathGlob)
	if err != nil || len(paths) == 0 {
		log.Fatalf("No brightness files found in %s", backlightPathGlob)
	}
	brightnessFile := paths[0]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(brightnessFile)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Watching brightness changes...")

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				current := getCurrentBrightness(brightnessFile)
				fmt.Printf("%d\n", current/960)

				// You can also send this to a socket, update a panel, or call notify-send
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}
