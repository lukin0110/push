package file

import (
    "os"
    "net/http"
    "io/ioutil"
    "fmt"
    "strings"
    "path"
    "golang.org/x/crypto/openpgp"
    "golang.org/x/crypto/openpgp/packet"
    "path/filepath"
    "io"
)

func UploadFile(url string, file string, email string) (string, error) {
    f, err := os.Open(file)
    if (err != nil) {
        return "", err
    }
    defer f.Close()

    req, err := http.NewRequest("PUT", url, f)
    if (err != nil) {
        return "", err
    }

    if email != "" {
        req.Header.Set("x-email", email)
    }
    fileStat, err1 := f.Stat()
    if err1 == nil {
        req.ContentLength = fileStat.Size()
    }

    //tr := &http.Transport{
    //    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    //}
    //client := &http.Client{Transport: tr}
    client := &http.Client{}

    res, err := client.Do(req)
    if (err != nil) {
        return "", err
    } else {
        defer res.Body.Close()
    }

    if(res.StatusCode != http.StatusOK) {
        bs, _ := ioutil.ReadAll(res.Body)
        return "", fmt.Errorf("bad status: %s, %s", res.Status, strings.Replace(string(bs), "\n", "", -1))
    } else {
        bs, err1 := ioutil.ReadAll(res.Body)
        if (err1 != nil) {
            return "", err
        }
        return strings.Replace(string(bs), "\n", "", -1), nil
    }
}

func Encrypt(fullPath string, passPhrase string) (resultFile *os.File, err error) {
    filename := path.Base(fullPath)
    stats, err := os.Stat(fullPath)
    if err != nil {
        return
    }
    hints := &openpgp.FileHints{IsBinary: true, FileName: filename, ModTime: stats.ModTime()}

    var config *packet.Config = &packet.Config{
        DefaultCompressionAlgo: packet.CompressionNone,
        CompressionConfig: &packet.CompressionConfig{packet.DefaultCompression},
    }

    sourceFile, err := os.Open(fullPath)
    if err != nil {
        return
    }
    defer sourceFile.Close()

    resultFile, err = ioutil.TempFile(filepath.Dir(fullPath), ".pushKiwiPGP")
    if err != nil {
        return
    }
    defer resultFile.Close()

    plaintext, err := openpgp.SymmetricallyEncrypt(resultFile, []byte(passPhrase), hints, config)
    if err != nil {
        return
    }
    defer plaintext.Close()

    // Copy original contents the the writer created previously
    _, err = io.Copy(plaintext, sourceFile)

    return
}
