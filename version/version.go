package version

import (
    "fmt"
)

const (
    Proto = "0"
    Major = "0"
    Minor = "2"
)

// $ figlet -m2 push.kiwi
const Kiwi = `
                    _         _     _            _
 _ __   _   _  ___ | |__     | | __(_)__      __(_)
| '_ \ | | | |/ __|| '_ \    | |/ /| |\ \ /\ / /| |
| |_) || |_| |\__ \| | | | _ |   < | | \ V  V / | |
| .__/  \__,_||___/|_| |_|(_)|_|\_\|_|  \_/\_/  |_|
|_|

`

func MajorMinor() string {
    return fmt.Sprintf("%s.%s", Major, Minor)
}

func Full() string {
    return fmt.Sprintf("%s.%s.%s", Proto, Major, Minor)
}

func Compat(client string, server string) bool {
    return client == server
}
