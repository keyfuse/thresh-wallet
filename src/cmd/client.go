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
	flagUID       string
	flagAPIurl    string
	flagChainNet  string
	flagMasterKey string
)

func init() {
	flag.StringVar(&flagAPIurl, "apiurl", "http://localhost:19099", "wallet server, default(http://localhost:19099)")
	flag.StringVar(&flagUID, "uid", "", "mobile or email")
	flag.StringVar(&flagChainNet, "chainnet", "testnet", "chainnet(testnet|mainnet), default(testnet)")
	flag.StringVar(&flagMasterKey, "masterkey", "", "master key wif string(default is random)")
}

func usage() {
	fmt.Println("Usage: " + os.Args[0] + " -uid=[...]")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = func() { usage() }
	flag.Parse()
	if flagUID == "" {
		usage()
		os.Exit(0)
	}
	cli := client.NewClient(flagAPIurl, flagUID, flagChainNet, flagMasterKey)
	cli.Start()
}
