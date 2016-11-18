#!/usr/bin/env bash

# Fix auto=completion of PyCharm
# Adding symbolic links in the vendor map own packages
rm -R vendor/github.com/lukin0110/push/ || true
mkdir -p vendor/github.com/lukin0110/push/
ln -sf "$(pwd)/version/" vendor/github.com/lukin0110/push/version
ln -sf "$(pwd)/file/" vendor/github.com/lukin0110/push/file
rm file/file || true
rm version/version || true
