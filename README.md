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

## Server

```
./bin/threshwallet-server -c conf/server.json.sample  -vcode off
```

## Client

```
./bin/threshwallet-client -mobile=10086 -apiurl=http://localhost:9099
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
| getnewaddress    | getnewaddress                          | getnewaddress                                                       |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| getsendfees      | getsendfees <address> <value>          | getsendfees tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw 10000        |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| sendtoaddress    | sendtoaddress <address> <value> <fees> | sendtoaddress tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw 10000 1000 |
+------------------+----------------------------------------+---------------------------------------------------------------------+
| sendalltoaddress | sendalltoaddress <address>             | sendalltoaddress tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw         |
+------------------+----------------------------------------+---------------------------------------------------------------------+
(13 rows)
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

## License

thresh-wallet is released under the GPLv3 License.


## References

[1] Y. Lindell. [Fast Secure Two-Party ECDSA Signing](https://eprint.iacr.org/2017/552.pdf)
