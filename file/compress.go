package file

import (
    "os"
    "archive/zip"
    "io"
    "path/filepath"
    "strings"
)

// http://blog.ralch.com/tutorial/golang-working-with-zip/
// TODO: avoid added targetName itself, report func
func ZippyMcZipface(targetName string, reporter func(r int64), items ...string) (err error) {
    zipfile, err := os.Create(targetName)
    if err != nil {
	return
    }
    defer zipfile.Close()

    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    for _, value := range items {
	info, err := os.Stat(value)
	if err != nil {
	    break
	}

	var baseDir string
	if info.IsDir() {
	    baseDir = filepath.Base(value)
	}

	filepath.Walk(value, func(path string, info os.FileInfo, err error) error {
	    if err != nil {
		return err
	    }

	    header, err := zip.FileInfoHeader(info)
	    if err != nil {
		return err
	    }

	    if baseDir != "" {
		header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, value))
	    }

	    if info.IsDir() {
		header.Name += "/"
	    } else {
		header.Method = zip.Deflate
	    }

	    writer, err := archive.CreateHeader(header)
	    if err != nil {
		return err
	    }

	    if info.IsDir() {
		return nil
	    }

	    file, err := os.Open(path)
	    if err != nil {
		return err
	    }
	    defer file.Close()
	    _, err = io.Copy(writer, file)
	    return err
	})
    }

    return
}
