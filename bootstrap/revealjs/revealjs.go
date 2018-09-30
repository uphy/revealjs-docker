package revealjs

//go:generate rice embed-go

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

type RevealJS struct {
	config        *Config
	directory     string
	dataDirectory string
	indexTemplate string
	EmbedSection  bool
}

const dataDirectoryName = "data"

func NewRevealJS(dir string) (*RevealJS, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, errors.New("`dir` not exist")
	}
	dataDirectory := filepath.Join(dir, dataDirectoryName)
	indexTemplate := filepath.Join(dataDirectory, "index.html.tmpl")
	return &RevealJS{nil, dir, dataDirectory, indexTemplate, true}, nil
}

func (r *RevealJS) reloadConfig() error {
	configFile := filepath.Join(r.dataDirectory, "config.yml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := r.generateInitialData(); err != nil {
			return err
		}
	}
	c, err := LoadConfigFile(configFile)
	if err != nil {
		return err
	}
	r.config = c
	return nil
}

func (r *RevealJS) generateInitialData() error {
	box, err := rice.FindBox("defaults")
	if err != nil {
		return err
	}
	box.Walk("/", func(path string, info os.FileInfo, err error) error {
		dest := filepath.Join(r.dataDirectory, path)
		if _, err := os.Stat(dest); !os.IsNotExist(err) {
			return nil
		}
		if info.IsDir() {
			if err := os.MkdirAll(dest, 0700); err != nil {
				return err
			}
		} else {
			in, err := box.Open(path)
			if err != nil {
				return err
			}
			defer in.Close()
			out, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, in)
			return err
		}
		return nil
	})
	return nil
}

func (r *RevealJS) Start() error {
	r.Reconfigure()
	go r.runRevealJS()
	watcher, err := NewWatcher(r)
	if err != nil {
		return err
	}
	go watcher.Start()
	return nil
}

func (r *RevealJS) runRevealJS() error {
	cmd := exec.Command("npm", "start")
	cmd.Dir = r.directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *RevealJS) Reconfigure() {
	if err := r.reloadConfig(); err != nil {
		log.Println("failed to reload config.yml: ", err)
	}
	if err := r.generateIndexHTML(); err != nil {
		log.Println("failed to generate index.html: ", err)
	}
}

func (r *RevealJS) generateIndexHTML() error {
	b, err := ioutil.ReadFile(r.indexTemplate)
	tmpl, err := template.New("index.html.tmpl").Parse(string(b))
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(r.directory, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := tmpl.Execute(f, map[string]interface{}{
		"config":   r.config,
		"sections": r.generateSections(),
	}); err != nil {
		return err
	}
	return nil
}

func (r *RevealJS) generateSections() []string {
	sections := make([]string, 0)
	if r.config.Slides == nil || len(r.config.Slides) == 0 {
		if err := filepath.Walk(r.dataDirectory, func(path string, info os.FileInfo, err error) error {
			p, _ := filepath.Rel(r.dataDirectory, path)
			section := r.sectionFor(p)
			if len(section) > 0 {
				sections = append(sections, section)
			}
			return nil
		}); err != nil {
			log.Println("failed to walk data directory: ", err)
		}
	} else {
		for _, s := range r.config.Slides {
			section := r.sectionFor(s)
			if len(section) > 0 {
				sections = append(sections, section)
			} else {
				log.Println("unsupported slide file: ", s)
			}
		}
	}
	return sections
}

func (r *RevealJS) sectionFor(file string) string {
	f := filepath.Join(dataDirectoryName, file)
	if r.EmbedSection {
		path := filepath.Join(r.directory, f)

		switch filepath.Ext(path) {
		case ".html":
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("failed to load file %s: %s", path, err)
				return ""
			}
			return fmt.Sprintf(`%s`, string(content))
		case ".md":
			return fmt.Sprintf(`<section data-markdown="%s"></section>`, f)
		default:
			return ""
		}
	} else {
		switch filepath.Ext(f) {
		case ".html":
			return fmt.Sprintf(`<section data-external="%s"></section>`, f)
		case ".md":
			return fmt.Sprintf(`<section data-markdown="%s"></section>`, f)
		default:
			return ""
		}
	}
}

func (r *RevealJS) UpdateSlideFile(file string) {
	r.Reconfigure()
}
