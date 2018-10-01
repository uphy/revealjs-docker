package revealjs

//go:generate rice embed-go

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var FilesetNames = []string{"default", "demo"}

func Generate(name string, dest string, force bool) error {
	if !isSupportedFilesetName(name) {
		return errors.New("unsupported fileset name: " + name)
	}
	if err := extractAll(name, dest, force); err != nil {
		return err
	}
	if err := extract("/index.html.tmpl", filepath.Join(dest, "index.html.tmpl"), nil, force); err != nil {
		return err
	}
	return nil
}

func isSupportedFilesetName(name string) bool {
	for _, n := range FilesetNames {
		if name == n {
			return true
		}
	}
	return false
}

func extractAll(srcDir string, destDir string, force bool) error {
	box.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if srcDir == path {
			return nil
		}
		rel, _ := filepath.Rel(srcDir, path)
		dest := filepath.Join(destDir, rel)
		return extract(path, dest, info, force)
	})
	return nil
}

func extract(src string, dest string, info os.FileInfo, force bool) error {
	if info == nil {
		f, err := box.Open(src)
		if err != nil {
			return err
		}
		i, _ := f.Stat()
		if err := f.Close(); i == nil || err != nil {
			return err
		}
		info = i
	}
	exist := exist(dest)
	if !force && exist {
		log.Println("Skipped.  File already exist:", dest)
		return nil
	}

	if force {
		if err := os.RemoveAll(dest); err != nil {
			return err
		}
	}

	if info.IsDir() {
		if err := os.MkdirAll(dest, 0700); err != nil {
			return err
		}
	} else {
		in, err := box.Open(src)
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
}
