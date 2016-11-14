#!/usr/bin/env bash

# Fix auto=completion of PyCharm
# Adding symbolic links in the vendor map own packages
mkdir -p vendor/github.com/lukin0110/push/
ln -s "$(pwd)/version/" vendor/github.com/lukin0110/push/version
ln -s "$(pwd)/file/" vendor/github.com/lukin0110/push/file
