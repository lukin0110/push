#!/bin/sh

case "$1" in
    release)
        echo 'Mac OS X'
        env GOOS=darwin go install -v github.com/lukin0110/push/cmd/push/
        cp /go/bin/darwin_amd64/push /go/bin/push.x86.darwin
        echo 'Linux'
        env GOOS=linux go install -v github.com/lukin0110/push/cmd/push/
        cp /go/bin/push /go/bin/push.x86.linux
    ;;

    *)
        echo 'In docker'
        env GOOS=linux go install -v github.com/lukin0110/push/cmd/push/
    ;;
esac
