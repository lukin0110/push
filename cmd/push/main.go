package main

import (
    "flag"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
    "io"
    "golang.org/x/crypto/openpgp"
    "golang.org/x/crypto/openpgp/packet"
    "path"
    ver "github.com/lukin0110/push/version"
    "github.com/docker/go-units"
)

const Url string = "https://push.kiwi/"
const UsageString string =
`Usage: push [OPTIONS] file...
       push [--help | --version]

Share a file from the command line. It returns an unique url to share.

Options:

 -e, --email        Share files via email
 -p, --passphrase   Protect files with a password
 -h, --help         Print usage
 -v, --version      Print version information and quit

Examples:

$ push ./nginx.conf
$ push --email dude@example.com ./nginx.conf
`
const MAX_BYTES int64 = 100 * 1024 * 1024 // 100 MegaByte


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

func flagBool(nameLong string, nameShort string, value bool, usage string) *bool {
    result := flag.Bool(nameLong, value, usage)
    flag.BoolVar(result, nameShort, value, usage)
    return result
}

func flagString(nameLong string, nameShort string, value string, usage string) *string {
    result := flag.String(nameLong, value, usage)
    flag.StringVar(result, nameShort, value, usage)
    return result
}

func encrypt(fullPath string, passPhrase string) (resultFile *os.File, err error) {
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

// Main function for the push command. Parses a few flag (optional) options and takes a few file, or more files, as
// command line arguments. All passed and existing files will be pushed, they will all return a shareable url.
func main() {
    // runtime.GOMAXPROCS(runtime.NumCPU() * 2)
    // Overwrite the default help message
    flag.Usage = func() {
    	fmt.Fprintln(os.Stderr, "See 'push --help'.")
    }

    email := flagString("email", "e", "", "Share files via email")
    passPhrase := flagString("passphrase", "p", "", "Protect files with a password")
    help := flagBool("help", "h", false, "Print usage")
    version := flagBool("version", "v", false, "Print version information and quit")
    flag.Parse()

    if *help {
        fmt.Println(UsageString)
        os.Exit(0)
    } else if *version {
        fmt.Println(ver.Full())
        os.Exit(0)
    }

    if len(flag.Args()) == 0 {
        fmt.Println(UsageString)
        os.Exit(0)
    }

    var toRemove = make([]string, len(flag.Args()))

    for index, v := range flag.Args() {
        var err error
        var stat os.FileInfo
        fullPath, err := filepath.Abs(v)

        if stat, err = os.Stat(fullPath); err == nil {
            var uploadPath string = fullPath
            var filename string = filepath.Base(fullPath)

            if stat.Size() < MAX_BYTES {
                if *passPhrase != "" {
                    var file *os.File
                    file, err = encrypt(fullPath, *passPhrase)
                    uploadPath, _ = filepath.Abs(file.Name())
                    filename += ".gpg"
                    toRemove[index] = uploadPath
                }

                var result string
                result, err = UploadFile(Url + filename, uploadPath, *email)
                //fmt.Printf("Uploading: %s, %s\n", uploadPath, *email); result := filename

                if err == nil {
                    fmt.Println(result)
                }
            } else {
                err = fmt.Errorf("Max file size (%s) exceeded for %s", units.BytesSize(float64(MAX_BYTES)), fullPath)
            }
        }

        if err != nil {
            fmt.Println(err)
        }
    }

    // Cleanup temporary encrypted files
    for _, v := range toRemove {
        if v != "" {
            os.Remove(v)
        }
    }
}
