package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "github.com/docker/go-units"
    ver "github.com/lukin0110/push/version"
    "github.com/lukin0110/push/file"
    "github.com/cheggaaa/pb"
)

const Url string = "https://push.kiwi/"
const UsageString string =
`Usage: push [OPTIONS] file...
       push [--help | --version]

Share a file from the command line. It returns an unique url to share. The file argument is
required, you can specify multiple files.

Options:

 -e, --email        Share files via email
 -p, --passphrase   Protect files with a password
 -h, --help         Print usage
 -v, --version      Print version information and quit

Examples:

$ push ./nginx.conf
$ push --email=jeffrey@lebowski.org ./nginx.conf
`
const MAX_BYTES int64 = 100 * 1024 * 1024 // 100 MegaByte


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

    email := flagString("email", "e", "", "Share files via email")
    passPhrase := flagString("passphrase", "p", "", "Protect files with a password")
    help := flagBool("help", "h", false, "Print usage")
    kiwi := flagBool("kiwi", "k", false, "Show a ascii art")
    version := flagBool("version", "v", false, "Print version information and quit")
    flag.Parse()

    if *help {
        fmt.Println(UsageString)
        os.Exit(0)
    } else if *version {
        fmt.Println(ver.Full())
        os.Exit(0)
    } else if *kiwi {
        fmt.Print(ver.Kiwi)
        os.Exit(0)
    }

    if *email != "" && !file.IsEmail(*email) {
        fmt.Printf("Invalid email: %s\n", *email)
        os.Exit(0)
    }

    if len(flag.Args()) == 0 {
        fmt.Println(UsageString)
        os.Exit(0)
    }

    var toRemove = make([]string, len(flag.Args()))
    var results = make([]string, len(flag.Args()))

    for index, v := range flag.Args() {
        var err error
        var stat os.FileInfo
        fullPath, err := filepath.Abs(v)

        if stat, err = os.Stat(fullPath); err == nil {
            var uploadPath string = fullPath
            var filename string = filepath.Base(fullPath)

            if stat.Size() < MAX_BYTES {
                if *passPhrase != "" {
                    var f *os.File
                    f, err = file.Encrypt(fullPath, *passPhrase)
                    uploadPath, _ = filepath.Abs(f.Name())
                    filename += ".gpg"
                    toRemove[index] = uploadPath
                }

                // Console progress bar for Golang
                // https://github.com/cheggaaa/pb
                uploadFile, err := os.Open(uploadPath)
                stats, err := os.Stat(uploadPath)
                bar := pb.New(int(stats.Size())).SetUnits(pb.U_BYTES).Prefix(filename)
                bar.Start()
                reader := bar.NewProxyReader(uploadFile)

                var result string
                result, err = file.UploadFile(Url + filename, *uploadFile, reader, *email)
                //fmt.Printf("Uploading: %s, %s\n", uploadPath, *email); result := filename
                bar.Finish()

                if err == nil {
                    results[index] = result
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

    // Print all results
    for _, v := range results {
        if v != "" {
            fmt.Println(v)
        }
    }
}
