package revealjs

//go:generate rice embed-go

import (
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
)

type RevealJS struct {
	config        *Config
	directory     string
	dataDirectory string
	indexTemplate string
	EmbedHTML     bool
	EmbedMarkdown bool
}

const (
	dataDirectoryName = "data"
	markdownSection   = `<section data-markdown="%s" data-separator="^\r?\n---\r?\n$" data-separator-vertical="^\r?\n~~~\r?\n$"></section>`
)

var box *rice.Box

func NewRevealJS(dir string) (*RevealJS, error) {
	if !exist(dir) {
		return nil, errors.New("`dir` not exist")
	}
	var err error
	box, err = rice.FindBox("files")
	if err != nil {
		panic(err)
	}
	dataDirectory := filepath.Join(dir, dataDirectoryName)
	indexTemplate := filepath.Join(dataDirectory, "index.html.tmpl")
	return &RevealJS{nil, dir, dataDirectory, indexTemplate, true, false}, nil
}

func (r *RevealJS) reloadConfig() error {
	configFile := filepath.Join(r.dataDirectory, "config.yml")
	if !exist(configFile) {
		if err := Generate(FilesetNames[0], r.dataDirectory, false); err != nil {
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
	path := filepath.Join(r.directory, f)

	switch filepath.Ext(path) {
	case ".html":
		if r.EmbedHTML {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("failed to load file %s: %s", path, err)
				return ""
			}
			return fmt.Sprintf(`%s`, string(content))
		}
		return fmt.Sprintf(`<section data-external="%s"></section>`, f)
	case ".md":
		if r.EmbedMarkdown {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("failed to read markdown file: %s", path)
			}
			return fmt.Sprintf(`<section data-markdown data-separator="^\r?\n---\r?\n$" data-separator-vertical="^\r?\n~~~\r?\n$">%s</section>`, html.EscapeString(string(b)))
		}
		return fmt.Sprintf(`<section data-markdown="%s" data-separator="^\r?\n---\r?\n$" data-separator-vertical="^\r?\n~~~\r?\n$"></section>`, f)
	default:
		return ""
	}
}

func (r *RevealJS) UpdateSlideFile(file string) {
	r.Reconfigure()
}

func (r *RevealJS) DataDirectory() string {
	return r.dataDirectory
}

func (r *RevealJS) Build() error {
	r.Reconfigure()
	dst := filepath.Join(r.directory, "build")
	// clean
	files, err := ioutil.ReadDir(dst)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.RemoveAll(filepath.Join(dst, f.Name())); err != nil {
			return err
		}
	}
	// build
	for _, src := range []string{"dist", "plugin", "index.html", "data"} {
		if err := copy(filepath.Join(r.directory, src), dst); err != nil {
			return err
		}
	}
	return nil
}

func copy(src, dst string) error {
	_, srcname := filepath.Split(src)
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		p := filepath.Join(dst, srcname, rel)
		if info.IsDir() {
			return os.MkdirAll(p, 0700)
		}
		reader, err := os.Open(path)
		if err != nil {
			return err
		}
		defer reader.Close()

		writer, err := os.Create(p)
		if err != nil {
			return err
		}
		defer writer.Close()

		_, err = io.Copy(writer, reader)
		return err
	})
}

func exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
