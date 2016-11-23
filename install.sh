#!/bin/sh

# Define the latest version
VERSION="0.0.2beta"
readonly VERSION

#####################################################################
# Check if a command with the name 'push' already exists on the
# system. Will exit with an error if it is not ours.
#####################################################################
# 'Our' push command should write to the stderr
#checker=$(push --unknown 2>&1 >/dev/null)
#echo $checker
#
#if (echo "$checker" | grep -q "flag provided but not defined: -unknown"); then
#    echo "Installing 'push'"
#else
#    echo "Error: could not install 'push'" 1>&2
#    echo "A 'push' command already exists in $(which push)" 1>&2
#	exit 1
#fi;


#####################################################################
# 1. Download to a temporary file
#    GitHub redirects to S3 (-L follows redirects)
# 2. Move & fix exec rights
#####################################################################

case "$(uname -s)" in
    Darwin)
        echo 'Mac OS X'
        temp_file=$(mktemp)
        curl -o "$temp_file" -L https://github.com/lukin0110/push/releases/download/$VERSION/push.x86.darwin
        mv "$temp_file" /usr/local/bin/push
        chmod 755 /usr/local/bin/push
        echo "Installed in /usr/local/bin/push"
    ;;

    Linux)
        temp_file=$(mktemp)
        arch=$(uname -p)
        if [ "$arch" = "unknown" ]; then
            arch="x86"
        fi;
        if [ "$arch" = "x86_64" ]; then
            arch="x86"
        fi;

        echo "Installing on Linux $arch"
        curl --fail -o "$temp_file" -L https://github.com/lukin0110/push/releases/download/$VERSION/push.$arch.linux
        mv "$temp_file" /usr/local/bin/push
        chmod 755 /usr/local/bin/push
        echo "Installed in /usr/local/bin/push"
    ;;

    *)
        echo 'other OS'
        echo 'Not supported ... yet. Open an issue on https://github.com/lukin0110/push/issues'
    ;;
esac
