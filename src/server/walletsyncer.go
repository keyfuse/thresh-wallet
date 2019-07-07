// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"sync"
	"time"

	"xlog"
)

// WalletSyncer --
type WalletSyncer struct {
	mu     sync.Mutex
	wg     sync.WaitGroup
	log    *xlog.Log
	conf   *Config
	done   chan bool
	store  *WalletStore
	chain  Chain
	ticker *time.Ticker
}

// NewWalletSyncer -- creates new WalletSyncer.
func NewWalletSyncer(log *xlog.Log, conf *Config, store *WalletStore) *WalletSyncer {
	return &WalletSyncer{
		log:    log,
		store:  store,
		conf:   conf,
		done:   make(chan bool),
		chain:  NewChainProxy(log, conf),
		ticker: time.NewTicker(time.Duration(time.Second * time.Duration(conf.WalletSyncIntervalSec))),
	}
}

// Start -- used to start the sync worker talk with chain.
func (ws *WalletSyncer) Start() {
	ws.wg.Add(1)
	go func(syncer *WalletSyncer) {
		defer syncer.wg.Done()
		defer syncer.ticker.Stop()

		for {
			select {
			case <-syncer.ticker.C:
				syncer.Sync()
			case <-syncer.done:
				return
			}
		}
	}(ws)
}

// Sync -- the sync worker.
func (ws *WalletSyncer) Sync() {
	log := ws.log
	store := ws.store
	chain := ws.chain

	ws.mu.Lock()
	defer ws.mu.Unlock()

	uids := store.UIDs()
	for _, uid := range uids {
		wallet := store.Get(uid)
		if wallet != nil {
			addresses := wallet.Addresses()
			for _, addr := range addresses {
				unspents, err := chain.GetUTXO(addr)
				if err != nil {
					log.Error("walletsyncer.address[%v].get.utxo.error:%v", addr, err)
					continue
				}
				wallet.UpdateAddressUnspent(addr, unspents)
			}
			if err := store.Write(wallet); err != nil {
				log.Error("walletsyncer.wallet[%v].store.write.error:%v", wallet.UID, err)
			}
		}
	}
}

// Stop -- used to stop the sync worker.
func (ws *WalletSyncer) Stop() {
	close(ws.done)
	ws.wg.Wait()
}
