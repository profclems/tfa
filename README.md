`tfa` is a Super cool two factor authentication command line app

## Usage

```sh
Usage:
  tfa [command]

Available Commands:
  add         Add a new account
  get         Get the OTP for an account
  help        Help about any command
  list        List the available 2fa codes
  version     show 2fa version information

Flags:
  -h, --help      help for tfa
  -v, --version   version for tfa

Use "tfa [command] --help" for more information about a command.
```

## Installation

Download a binary suitable for your OS at the [releases](https://github.com/profclems/tfa/releases/latest) page.

If you have `go` installed, run:
```
$ go get -v github.com/profclems/tfa/cmd/tfa
```
