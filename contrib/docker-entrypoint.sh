#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
# http://stackoverflow.com/questions/19622198/what-does-set-e-mean-in-a-bash-script
set -e

# Define help message
show_help() {
    echo """
Usage: docker run <imagename> COMMAND
Commands:
bash     : Start a bash shell
release  : Generate release builds
mac      : Generate release build for OSX
build    : Build locally for dev purposes
fetch    : Fetch a Go package, e.g: docker-compose run app fetch github.com/asaskevich/govalidator@v5
help     : Show this message
"""
}

# Run
case "$1" in
    bash)
        /bin/bash "${@:2}"
    ;;
    release)
        echo 'Mac OSX'
        env GOOS=darwin go install -v github.com/lukin0110/push/cmd/push/
        cp /go/bin/darwin_amd64/push /output/push.x86.darwin
        echo 'Linux'
        env GOOS=linux go install -v github.com/lukin0110/push/cmd/push/
        cp /go/bin/push /output/push.x86.linux
    ;;
    mac)
        echo 'Mac OSX'
        env GOOS=darwin go install -v github.com/lukin0110/push/cmd/push/
        cp /go/bin/darwin_amd64/push /output/push
    ;;
    build)
        echo 'In docker'
        env GOOS=linux go install -v github.com/lukin0110/push/cmd/push/
    ;;
    fetch)
        echo 'Fetching (with govendor)'
        govendor fetch "${@:2}"
    ;;
    *)
        show_help
    ;;
esac
