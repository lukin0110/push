# Push

Share a file from the command line.

```
$ push slack.dmg
https://push.kiwi/HrfUfSAzj/slack.dmg
```

## Install

OSX / Linux:
```
curl -sL https://raw.githubusercontent.com/lukin0110/push/master/install.sh | sudo -E bash -
```

## Usage

```
Usage: push [OPTIONS] file...
       push [--help | --version]

Share a file from the command line. It returns an unique url to share. The file argument is
required, you can specify multiple files.

Options:

 -e, --email        Share files via email
 -p, --passphrase   Protect files with a password
 -h, --help         Print usage
 -v, --version      Print version information and quit

Examples:

$ push ./nginx.conf
$ push --email=jeffrey@lebowski.org ./nginx.conf
$ push --passphrase=Security007 ./nginx.conf
```

## Roadmap

* [Add zip compression](https://github.com/lukin0110/push/issues/3)
* [Parallel uploads](https://github.com/lukin0110/push/issues/4)
* [Capture sender info](https://github.com/lukin0110/push/issues/5)
