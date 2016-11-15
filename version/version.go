package version

import (
    "fmt"
)

const (
    Proto = "0"
    Major = "0"
    Minor = "2"
)

func MajorMinor() string {
    return fmt.Sprintf("%s.%s", Major, Minor)
}

func Full() string {
    return fmt.Sprintf("%s.%s.%s", Proto, Major, Minor)
}

func Compat(client string, server string) bool {
    return client == server
}
