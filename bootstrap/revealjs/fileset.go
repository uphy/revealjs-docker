package revealjs

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var FilesetNames = []string{"default", "demo"}

func Generate(name string, dest string, force bool) error {
	if !isSupportedFilesetName(name) {
		return errors.New("unsupported fileset name: " + name)
	}
	// Create 'dest' directory if not exist
	if err := os.MkdirAll(dest, 0700); err != nil {
		return err
	}
	if err := extractAll(name, dest, force); err != nil {
		return err
	}
	if err := extract("files/index.html.tmpl", filepath.Join(dest, "index.html.tmpl")); err != nil {
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
	// Clean destDir if force==true
	if force {
		if exist(destDir) {
			files, err := os.ReadDir(destDir)
			if err != nil {
				return err
			}
			for _, f := range files {
				if err := os.RemoveAll(filepath.Join(destDir, f.Name())); err != nil {
					return err
				}
			}
		}
	}

	embedSrcDir := filepath.Join("files", srcDir)
	return fs.WalkDir(embedFS, embedSrcDir, func(path string, d fs.DirEntry, err error) error {
		if embedSrcDir == path {
			return nil
		}
		rel, _ := filepath.Rel(embedSrcDir, path)
		dest := filepath.Join(destDir, rel)
		return extract(path, dest)
	})
}

func extract(src string, dest string) error {
	// Get file info of src file
	f, err := embedFS.Open(src)
	if err != nil {
		return err
	}
	info, _ := f.Stat()
	if err := f.Close(); err != nil {
		return err
	}

	// If dest file already exist, skip
	exist := exist(dest)
	if exist {
		log.Println("Skipped.  File already exist:", dest)
		return nil
	}

	if info.IsDir() {
		// Create dir
		if err := os.MkdirAll(dest, 0700); err != nil {
			return err
		}
	} else {
		// Copy file
		in, err := embedFS.Open(src)
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
