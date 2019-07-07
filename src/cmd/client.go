// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"client"
)

var (
	flagAPIurl    string
	flagMobile    string
	flagChainNet  string
	flagMasterKey string
)

func init() {
	flag.StringVar(&flagAPIurl, "apiurl", "http://localhost:19099", "wallet server, default(http://localhost:19099)")
	flag.StringVar(&flagMobile, "mobile", "", "mobile number")
	flag.StringVar(&flagChainNet, "chainnet", "testnet", "chainnet(testnet|mainnet), default(testnet)")
	flag.StringVar(&flagMasterKey, "masterkey", "", "master key wif string(default is random)")
}

func usage() {
	fmt.Println("Usage: " + os.Args[0] + " --mobile=[...]")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = func() { usage() }
	flag.Parse()
	if flagMobile == "" {
		usage()
		os.Exit(0)
	}
	cli := client.NewClient(flagAPIurl, flagMobile, flagChainNet, flagMasterKey)
	cli.Start()
}
