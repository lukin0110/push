package main

import "fmt"
import (
	"flag"
	"os"
)

func main() {
	flag.Usage = func() {
	    fmt.Fprintln(os.Stderr, "See 'push --help'.")
	}

	email := flag.String("e", "", "Email address")
	help := flag.Bool("help", false, "Print usage")
	flag.Parse()

	if(*help) {
		fmt.Println("Print usage")
	} else {
		fmt.Printf("Email: %s\n", *email)
		fmt.Println("tail:", flag.Args())
	}
}
