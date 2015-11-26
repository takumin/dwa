package main

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
)

func Download(url, directory, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))

	var fullpath string
	if params["filename"] != "" {
		fullpath = path.Join(directory, params["filename"])
	} else if filename != "" {
		fullpath = path.Join(directory, filename)
	} else {
		return errors.New("Require filename.")
	}

	dst, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer dst.Close()

	n, err := io.Copy(dst, resp.Body)
	if n < 0 {
		return errors.New("Fatal size.")
	}
	if err != nil {
		return err
	}

	return nil
}
