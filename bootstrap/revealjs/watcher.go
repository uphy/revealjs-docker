package revealjs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	revealjs *RevealJS
}

func NewWatcher(revealjs *RevealJS) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{w, revealjs}, err
}

func (w *Watcher) Start() {
	w.watcher.Add(w.revealjs.dataDirectory)
	filepath.Walk(w.revealjs.dataDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			w.watcher.Add(path)
		}
		return nil
	})
	w.watcher.Add(w.revealjs.indexTemplate)
	for evt := range w.watcher.Events {
		op := evt.Op
		if op&fsnotify.Create != 0 {
			if s, err := os.Stat(evt.Name); !os.IsNotExist(err) && s.IsDir() {
				w.watcher.Add(evt.Name)
			} else {
				w.revealjs.Reconfigure()
			}
		} else if op&fsnotify.Remove != 0 {
			w.watcher.Remove(evt.Name)
			w.revealjs.Reconfigure()
		} else if op&fsnotify.Write != 0 {
			if _, file := filepath.Split(evt.Name); file == "config.yml" || strings.HasSuffix(file, ".tmpl") {
				w.revealjs.Reconfigure()
			} else {
				w.revealjs.UpdateSlideFile(evt.Name)
			}
		}
	}
}
