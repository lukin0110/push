package main

import (
    "flag"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
)

const Url string = "https://push.kiwi/"
const Version string = "0.0.1"
const UsageString string =
`Usage: push [OPTIONS] file...
       push [--help | --version]

Share a file from the command line. It returns an unique url to share.

Options:

 -e, --email    Share via email
 -h, --help     Print usage
 -v, --version  Print version information and quit

Examples:

$ push ./nginx.conf
$ push --email dude@example.com ./nginx.conf
`

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

    //req.Header.Set("Content-Type", "text/markdown; charset=UTF-8")

    //tr := &http.Transport{
    //    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    //}
    //client := &http.Client{Transport: tr}
    client := &http.Client{}

    res, err := client.Do(req)
    defer res.Body.Close()
    if (err != nil) {
        return "", nil
    }

    if(res.StatusCode != http.StatusOK) {
        return "", fmt.Errorf("bad status: %s", res.Status)
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

// Main function for the push command. Parses a few flag (optional) options and takes a few file, or more files, as
// command line arguments. All passed and existing files will be pushed, they will all return a shareable url.
func main() {
    // Overwrite the default help message
    flag.Usage = func() {
    	fmt.Fprintln(os.Stderr, "See 'push --help'.")
    }

    email := flagString("email", "e", "", "Share via email")
    help := flagBool("help", "h", false, "Print usage")
    version := flagBool("version", "v", false, "Print version information and quit")
    flag.Parse()

    if *help {
        fmt.Println(UsageString)
        os.Exit(0)
    } else if *version {
        fmt.Println(Version)
        os.Exit(0)
    }

    if len(flag.Args()) == 0 {
        fmt.Println(UsageString)
        os.Exit(0)
    }

    for _, v := range flag.Args() {
        var err error
        fullPath, err := filepath.Abs(v)

        if _, err1 := os.Stat(fullPath); err1 == nil {
            filename := filepath.Base(fullPath)
            result, err := UploadFile(Url + filename, fullPath, *email)

            if (err == nil) {
                fmt.Println(result)
            }
        } else {
            err = fmt.Errorf("Doesn't exist: %s\n", fullPath)
        }

        if err != nil {
            fmt.Print(err)
        }
    }
}
