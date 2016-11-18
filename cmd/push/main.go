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
const MaxBytes int64 = 100 * 1024 * 1024 // 100 MegaByte
const UsageString string =
`Usage: push [OPTIONS] file...
       push [--help | --version]

Share a file from the command line. It returns an unique url to share. The file argument is
required, you can specify multiple files.

Options:

 -e, --email        Share files via email
 -p, --passphrase   Protect files with a password
 -z, --zip          Compress files to one archive
 -h, --help         Print usage
 -v, --version      Print version information and quit

Examples:

$ push ./nginx.conf
$ push --email=jeffrey@lebowski.org ./nginx.conf
$ push --passphrase=Security007 ./nginx.conf
`


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

func HandleFile(filePath string, passPhrase string, email string) (url string, err error) {
    var stat os.FileInfo
    fullPath, err := filepath.Abs(filePath)

    if stat, err = os.Stat(fullPath); err == nil {
        var uploadPath string = fullPath
        var filename string = filepath.Base(fullPath)

        if stat.Size() < MaxBytes {
            if passPhrase != "" {
                var f *os.File
                f, err = file.Encrypt(fullPath, passPhrase)
                uploadPath, _ = filepath.Abs(f.Name())
                filename += ".gpg"
                defer func() {
                    os.Remove(uploadPath)
                }()
            }

            // Console progress bar for Golang
            // https://github.com/cheggaaa/pb
            uploadFile, err := os.Open(uploadPath)
            if err != nil {
                return "", err
            }
            stats, err := os.Stat(uploadPath)
            if err != nil {
                return "", err
            }
            bar := pb.New(int(stats.Size())).SetUnits(pb.U_BYTES).Prefix(filename)
            bar.Start()
            reader := bar.NewProxyReader(uploadFile)

            url, err = file.UploadFile(Url + filename, *uploadFile, reader, email)
            //fmt.Printf("Uploading: %s, %s\n", uploadPath, *email); result := filename
            bar.Finish()
        } else {
            err = fmt.Errorf("Max file size (%s) exceeded for %s", units.BytesSize(float64(MaxBytes)), fullPath)
        }
    }

    return
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
    zip := flagBool("zip", "z", false, "Compress files to one archive")
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

    var results = make([]string, 0) // List of urls all successful uploads
    var errors = make([]error, 0)   // List of failed uploads (file not found, failed to upload, etc)

    if *zip {
        total, err := file.CalculateSize(flag.Args()...)

        if err == nil {
            bar := pb.New(int(total)).SetUnits(pb.U_BYTES).Prefix("Compressing")
            bar.Start()

            var progressCount int = 0
            reporter := func(w int64) {
                progressCount += int(w)
                bar.Set(progressCount)
            }

            zipfileName := "zippy.zip"
            err := file.ZippyMcZipface(zipfileName, reporter, flag.Args()...)
            bar.Finish()

            if err == nil {
                // Bye bye zippy
                defer os.Remove(zipfileName)

                var url string
                url, err = HandleFile(zipfileName, *passPhrase, *email)
                if url != "" {
                    results = append(results, url)
                }
            }
        }

        if err != nil {
            errors = append(errors, err)
        }

    } else {
        // Loopy MacLoopface
        for _, v := range flag.Args() {
            url, err := HandleFile(v, *passPhrase, *email)
            if url != "" {
                results = append(results, url)
            } else {
                errors = append(errors, err)
            }
        }
    }

    // Print all results
    for _, v := range results {
        if v != "" {
            fmt.Println(v)
        }
    }

    // Print all errors
    if len(errors) > 0 {
        fmt.Println("\nErrors:")
        for _, v := range errors {
            fmt.Printf(" - %s\n", v)
        }
    }
}
