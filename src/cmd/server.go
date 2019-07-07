// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"server"
	"xlog"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "c", "", "config file")
}

func usage() {
	fmt.Println("Usage: " + os.Args[0] + " [-c] <config-file>")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))

	// Load config.
	flag.Usage = func() { usage() }
	flag.Parse()
	if flagConf == "" {
		usage()
		os.Exit(0)
	}
	conf, err := server.LoadConfig(flagConf)
	if err != nil {
		log.Panic("server.load.config.error[%v]", err)
	}
	log.Info("%+v", conf)

	router := server.NewAPIRouter(log, conf)
	go http.ListenAndServe(conf.Endpoint, router)

	// Handle SIGINT and SIGTERM signals.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Info("server.got.signal:%+v", <-ch)
	log.Info("server.exit.done")
}
