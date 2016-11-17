package file

import (
    "os"
    "archive/zip"
    "io"
    "path/filepath"
    "strings"
    "errors"
    "io/ioutil"
)

// Similar to: http://stackoverflow.com/questions/26050380/go-tracking-post-request-progress
type ProgressWriter struct {
    io.Writer
    Reporter func(r int64)
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
    n, err = pw.Writer.Write(p)
    pw.Reporter(int64(n))
    return
}

// Calculates to total amount of bytes from a list of files/directories
func CalculateSize(items ...string) (total int64, err error) {
    total = 0

    for _, value := range items {
	var info os.FileInfo
	info, err = os.Stat(value)
	if err != nil {
	    return
	}

	if !info.IsDir() {
	    total += info.Size()
	} else {
	    filepath.Walk(value, func(path string, info os.FileInfo, err error) error {
		if err != nil {
		    return err
		}

		info, err = os.Stat(path)
		if err != nil {
		    return err
		}

		if !info.IsDir() {
		    total += info.Size()
		}
		return nil
	    })
	}
    }

    return
}

// Creates a zipfile from a list of files/directories
//
// http://blog.ralch.com/tutorial/golang-working-with-zip/
func ZippyMcZipface(targetName string, reporter func(w int64), items ...string) (err error) {
    if reporter == nil {
	return errors.New("reporter function can't be nil")
    }

    tempFile, err := ioutil.TempFile(os.TempDir(), "pushKiwiTmp")
    if err != nil {
	return
    }
    defer tempFile.Close()

    archive := zip.NewWriter(tempFile)
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
	    proxyWriter := &ProgressWriter{writer, reporter}

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
	    _, err = io.Copy(proxyWriter, file)
	    return err
	})
    }

    // Move temp file to target
    absPath, err := filepath.Abs(targetName)
    if err != nil {
	return
    }

    // Won't work within a docker container if volumes are mounted. Will produce an error:
    // rename /tmp/pushKiwiTmp706257135 /go/src/github.com/lukin0110/push/test.zip: invalid cross-device link
    err = os.Rename(tempFile.Name(), absPath)
    return
}
