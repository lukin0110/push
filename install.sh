#!/bin/sh

#####################################################################
# 1. Download to a temporary file
#    GitHub redirects to S3 (-L follows redirects)
# 2. Move & fix exec rights
#####################################################################

case "$(uname -s)" in
    Darwin)
        echo 'Mac OS X'
        temp_file=$(mktemp)
        curl -o "$temp_file" -L https://github.com/lukin0110/push/releases/download/0.0.1a/push.x86.darwin
        mv "$temp_file" /usr/local/bin/push
        chmod 755 /usr/local/bin/push
    ;;

    Linux)
        echo 'Linux'
        temp_file=$(mktemp)
        curl -o "$temp_file" -L https://github.com/lukin0110/push/releases/download/0.0.1a/push.x86.linux
        mv "$temp_file" /usr/local/bin/push
        chmod 755 /usr/local/bin/push
    ;;

    CYGWIN*|MINGW32*|MSYS*)
        echo 'MS Windows'
    ;;

    *)
        echo 'other OS'
        echo 'Not supported ... yet. Open an issue on https://github.com/lukin0110/push/issues'
    ;;
esac
