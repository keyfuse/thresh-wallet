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
(9 rows)
threshwallet@testnet> gettoken xx
+--------+
| status |
+--------+
| OK     |
+--------+
(1 rows)
threshwallet@testnet> getnewaddress
+--------------------------------------------+---------+
|                  address                   | postion |
+--------------------------------------------+---------+
| tb1qv8v4grqqvjuhwmn6xk9wnkdk7gem8m0uua4w62 |       0 |
+--------------------------------------------+---------+
(1 rows)
threshwallet@testnet> getbalance
+-----------------+
| current_balance |
+-----------------+
|         7055957 |
+-----------------+
(1 rows)
threshwallet@testnet> getnewaddress
+--------------------------------------------+---------+
|                  address                   | postion |
+--------------------------------------------+---------+
| tb1qalf2fk6g3xva2pxhzvgfy37009nd47tvdn57cv |       1 |
+--------------------------------------------+---------+
(1 rows)
threshwallet@testnet> sendalltoaddress tb1qalf2fk6g3xva2pxhzvgfy37009nd47tvdn57cv
+--------------------------------------------+------------+-----------+------------------------------------------------------------------+
|                 toaddress                  | value(sat) | fees(sat) |                               txid                               |
+--------------------------------------------+------------+-----------+------------------------------------------------------------------+
| tb1qalf2fk6g3xva2pxhzvgfy37009nd47tvdn57cv |    7055732 |       225 | 61cb003e443d1c48eb4ab1bfea98101a5e6f1488ebdac4c00371a7853c8d516b |
+--------------------------------------------+------------+-----------+------------------------------------------------------------------+
(1 rows)
threshwallet@testnet> getbalance
+-----------------+
| current_balance |
+-----------------+
|         7055732 |
+-----------------+
(1 rows)
threshwallet@testnet>
```

## License

thresh-wallet is released under the GPLv3 License.


## References

[1] Y. Lindell. [Fast Secure Two-Party ECDSA Signing](https://eprint.iacr.org/2017/552.pdf)

[2] [Bitcoin testnet3 faucet](https://coinfaucet.eu/en/btc-testnet/)
