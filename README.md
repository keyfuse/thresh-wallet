# thresh-wallet

A Bitcoin wallet powered by two-party ECDSA written in golang.

[![Build Status](https://travis-ci.org/keyfuse/thresh-wallet.png)](https://travis-ci.org/keyfuse/thresh-wallet) [![Go Report Card](https://goreportcard.com/badge/github.com/keyfuse/thresh-wallet)](https://goreportcard.com/report/github.com/keyfuse/thresh-wallet) [![codecov.io](https://codecov.io/gh/keyfuse/thresh-wallet/graphs/badge.svg)](https://codecov.io/gh/keyfuse/thresh-wallet/branch/master)

## How to Build

To build thresh-wallet from the source code you need to have a working
Go environment with [version 1.12 or greater installed](https://golang.org/doc/install).

#### Build Client/Server

```
$ git clone https://github.com/keyfuse/thresh-wallet
$ cd thresh-wallet
$ export GOPATH=`pwd`
$ make build
```

#### Build IOS Library

```
$ make buildosx
```

#### Build Android Library

```
$ make buildandroid
```

## License

thresh-wallet is released under the GPLv3 License.
