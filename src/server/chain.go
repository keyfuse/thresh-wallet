// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"xlog"
)

const (
	testnet = "testnet"
	mainnet = "mainnet"
)

// Chain --
type Chain interface {
	GetTxs(address string) ([]Tx, error)
	GetFees() (map[string]float32, error)
	GetUTXO(address string) ([]Unspent, error)
	GetTickers() (map[string]Ticker, error)
	PushTx(hex string) (string, error)
}

// NewChainProxy -- creates new Chain, default provider is blockstream.info.
func NewChainProxy(log *xlog.Log, conf *Config) Chain {
	return NewBlockstreamChain(log, conf)
}
