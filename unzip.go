package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Extracts struct {
	Src string
	Dst string
	Rep string
	Cut uint
}

func Unzip(ext Extracts) error {
	r, err := zip.OpenReader(ext.Src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		file := f.Name

		if ext.Rep != "" {
			file = regexp.MustCompile(ext.Rep).ReplaceAllString(file, "")
		}

		if ext.Cut > 0 {
			file = strings.Join(regexp.MustCompile("/").Split(file, -1)[ext.Cut:], "/")
		}

		path := filepath.Join(ext.Dst, file)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
