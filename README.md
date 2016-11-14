# Push

Share a file from the command line.

```
$ push README.md
```

## Install

OSX / Debian / Ubuntu:
```
curl -sL https://raw.githubusercontent.com/lukin0110/push/master/install.sh | sudo -E bash -
```

## Usage

```
Usage: push [OPTIONS] file...
       push [--help | --version]

Share a file from the command line. It returns an unique url to share.

Options:

 -e, --email        Share files via email
 -p, --passphrase   Protect files with a password
 -h, --help         Print usage
 -v, --version      Print version information and quit

Examples:

$ push ./nginx.conf
$ push --email dude@example.com ./nginx.conf
$ push -p Security123 slack.dmg
```
