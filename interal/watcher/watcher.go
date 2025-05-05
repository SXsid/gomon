package watcher

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	config "github.com/SXsid/gomon/interal/Config"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	config       *config.Config
	isDebouncing bool
	lastEvent    time.Time
	fsWatcher    *fsnotify.Watcher
	eventChannel chan struct{} //this is the common chaneel across the file
	DoneChannel  chan struct{} // to shudown the watcher
}

func NewWatcher(eventChan chan struct{}, Cfg *config.Config) (*FileWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		config:       Cfg,
		eventChannel: eventChan,
		DoneChannel:  make(chan struct{}),
		fsWatcher:    fsWatcher,
		isDebouncing: false,
		lastEvent:    time.Now(),
	}, nil

}

func (fw *FileWatcher) Start() error {
	log.Println("watching the file")
	err := filepath.Walk(fw.config.WatchDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//ignore if it's included in the pattern to ignore
		for _, pattern := range fw.config.IgnorePattern {
			matched, err := filepath.Match(pattern, path)
			if err != nil {
				return err
			}
			if matched {
				// stop recusrion for this folder
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.IsDir() {
			color.Yellow("Watching directory: %s", path)

			if err := fw.fsWatcher.Add(path); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	//start watching on the folder
	go fw.Watchevent()

	return nil
}

// long running process
func (fw *FileWatcher) Watchevent() {
	for {
		select {

		case events, ok := <-fw.fsWatcher.Events:
			if !ok {
				//kill the process
				return
			}

			if !shouldWatchFile(events.Name) {
				//ignore
				continue
			}

			if events.Op&fsnotify.Write == fsnotify.Write {

				info, err := os.Stat(events.Name)
				if err != nil {
					log.Println("warning:couldnt detect the changes in ", events.Name)
				} else if info.IsDir() {
					continue
				}

			}
			if events.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(events.Name)
				if err == nil && info.IsDir() {
					// if we creat a new subdir
					log.Println("Adding new subdirectory to watcher:", events.Name)
					fw.fsWatcher.Add(events.Name)
				}
			}

			if events.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {

				fw.handleEvents()
			}

		case err, ok := <-fw.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching files: %v", err)

		case <-fw.DoneChannel:
			//kill the process by maintainer
			return
		}
	}
}

func shouldWatchFile(fileName string) bool {

	if strings.HasSuffix(fileName, ".go") ||
		strings.HasSuffix(fileName, ".mod") ||
		strings.HasSuffix(fileName, ".sum") {
		return true
	}
	dataExtensions := []string{
		".json", ".yaml", ".yml", ".toml", ".xml",
		".csv", ".txt", ".env", ".ini", ".conf",
	}

	for _, extnsion := range dataExtensions {
		if strings.HasSuffix(fileName, extnsion) {
			return true
		}
	}
	return false

}

//send the instru to the controller while keeping th beboudning in check
//bebouncing=>

func (fw *FileWatcher) handleEvents() {
	now := time.Now()
	if fw.isDebouncing {
		fw.lastEvent = now
	} else {
		fw.isDebouncing = true
		fw.lastEvent = now

		go func() {
			firstDebounceTime := time.Now()
			durationtime := time.Duration(300) * time.Millisecond
			for {

				timer := time.NewTimer(durationtime)
				select {

				case <-timer.C:
					//if some changes occur btw deboucing re start the debouning but force build after 2 sec
					if time.Since(fw.lastEvent) >= durationtime || time.Since(firstDebounceTime) > time.Duration(2)*time.Second {
						fw.eventChannel <- struct{}{}
						fw.isDebouncing = false
						return
					}
					//stop the routine on if parent is stopped
				case <-fw.DoneChannel:
					color.Red("it done now shutting down")
					return
				}

			}

		}()

	}

}
