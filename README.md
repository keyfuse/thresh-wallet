# thresh-wallet

A cross-platform keyless Bitcoin wallet service powered by Breakthrough Cryptography.

[![Build Status](https://travis-ci.org/keyfuse/thresh-wallet.png)](https://travis-ci.org/keyfuse/thresh-wallet) [![Go Report Card](https://goreportcard.com/badge/github.com/keyfuse/thresh-wallet)](https://goreportcard.com/report/github.com/keyfuse/thresh-wallet)


## What is a keyless thresh wallet?

No private keys to worry about.

The wallet security is **distributed** between your device and the server.

The server knows **absolutely** no information related to your private key.

## Platforms

- iOS
- OSX
- Linux
- Android

## How to Build

To build thresh-wallet from the source code you need to have a working Go environment with [version 1.12 or greater installed](https://golang.org/doc/install).

#### Building

```
$ git clone https://github.com/keyfuse/thresh-wallet
$ cd thresh-wallet
$ export GOPATH=`pwd`
$ make build
```

#### IOS Library

```
$ make buildosx
```

#### Android Library

```
$ make buildandroid
```

## Try the Demo

####  Server
```
./bin/threshwallet-server -c conf/server.json.sample  -vcode off
```

####  Client

```
./bin/threshwallet-client -uid=xx@xx.com -apiurl=http://localhost:9099
+------------------+----------------------------------------+---------------------------------------------------------------------+
|     commands     |                 usage                  |                               example                               |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| help             | help                                   | help                                                                |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| dumpkey          | dumpkey                                | help                                                                |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| gettoken         | gettoken <vcode>                       | gettoken 666888                                                     |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| checkwallet      | checkwallet                            | checkwallet                                                         |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| createwallet     | createwallet                           | createwallet                                                        |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| backupwallet     | backupwallet                           | backupwallet                                                        |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| recoverwallet    | recoverwallet                          | recoverwallet                                                       |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| getbalance       | getbalance                             | getbalance                                                          |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| gettxs           | gettxs                                 | gettxs                                                              |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| getaddresses     | getaddresses                           | getaddresses                                                        |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| getnewaddress    | getnewaddress                          | getnewaddress                                                       |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| getsendfees      | getsendfees <address> <value>          | getsendfees tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw 10000        |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| sendtoaddress    | sendtoaddress <address> <value> <fees> | sendtoaddress tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw 10000 1000 |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| sendalltoaddress | sendalltoaddress <address>             | sendalltoaddress tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw         |
+------------------+----------------------------------------+---------------------------------------------------------------------+
(14 rows)

threshwallet@testnet> gettoken xx
+--------+
| status |
+--------+
| OK     |
+--------+
(1 rows)
threshwallet@testnet> checkwallet
+-------------+---------------+
| user_exists | backup_exists |
+-------------+---------------+
| false       | false         |
+-------------+---------------+
(1 rows)
threshwallet@testnet> createwallet
+--------+
| status |
+--------+
| OK     |
+--------+
(1 rows)
threshwallet@testnet> checkwallet
+-------------+---------------+
| user_exists | backup_exists |
+-------------+---------------+
| true        | false         |
+-------------+---------------+
(1 rows)
threshwallet@testnet> backupwallet
+--------+
| status |
+--------+
| OK     |
+--------+
(1 rows)
threshwallet@testnet> checkwallet
+-------------+---------------+
| user_exists | backup_exists |
+-------------+---------------+
| true        | true          |
+-------------+---------------+
(1 rows)
threshwallet@testnet> getnewaddress
+--------------------------------------------+---------+
|                  address                   | postion |
+--------------------------------------------+---------+
| tb1qzg2p5je5z82elv3k08lzpyjcwv9830he0qt4c7 |       0 |
+--------------------------------------------+---------+
(1 rows)
threshwallet@testnet> getbalance
+-----------------+
| current_balance |
+-----------------+
|               0 |
+-----------------+
(1 rows)
threshwallet@testnet> getbalance
+-----------------+
| current_balance |
+-----------------+
|         1000000 |
+-----------------+
(1 rows)
threshwallet@testnet> getnewaddress
+--------------------------------------------+---------+
|                  address                   | postion |
+--------------------------------------------+---------+
| tb1qx0ehpa5a7ld2mhv99nkvzpsx343ap7aqhmtca0 |       1 |
+--------------------------------------------+---------+
(1 rows)
threshwallet@testnet> sendalltoaddress tb1qx0ehpa5a7ld2mhv99nkvzpsx343ap7aqhmtca0
+--------------------------------------------+------------+-----------+------------------------------------------------------------------+
|                 toaddress                  | value(sat) | fees(sat) |                               txid                               |
+--------------------------------------------+------------+-----------+------------------------------------------------------------------+
| tb1qx0ehpa5a7ld2mhv99nkvzpsx343ap7aqhmtca0 |     999775 |       225 | 1d73d5e424b3fce1b7d5fb221e069a52929276cd1eaec8d6ed1a4295a17149c4 |
+--------------------------------------------+------------+-----------+------------------------------------------------------------------+
(1 rows)
threshwallet@testnet> getbalance
+-----------------+
| current_balance |
+-----------------+
|          999775 |
+-----------------+
(1 rows)
threshwallet@testnet>
```

## Can I trust this code?
*Don't trust. Verify.*

## License

thresh-wallet is released under the GPLv3 License.


## References

[1] Y. Lindell. [Fast Secure Two-Party ECDSA Signing](https://eprint.iacr.org/2017/552.pdf)
