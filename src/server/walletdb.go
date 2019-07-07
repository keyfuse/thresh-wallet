// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"

	"xlog"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore/bip32"
)

// WalletDB --
type WalletDB struct {
	log    *xlog.Log
	conf   *Config
	net    *network.Network
	store  *WalletStore
	syncer *WalletSyncer
}

// NewWalletDB -- creates new WalletDB.
func NewWalletDB(log *xlog.Log, conf *Config) *WalletDB {
	var net *network.Network
	switch conf.ChainNet {
	case testnet:
		net = network.TestNet
	case mainnet:
		net = network.MainNet
	}

	store := NewWalletStore(log, conf)
	syncer := NewWalletSyncer(log, conf, store)
	return &WalletDB{
		log:    log,
		net:    net,
		conf:   conf,
		store:  store,
		syncer: syncer,
	}
}

func (wdb *WalletDB) check(uid string, walletMasterPubKey string, cliMasterPubKey string) error {
	log := wdb.log

	if walletMasterPubKey != cliMasterPubKey {
		log.Error("wdb.openwallet[%v].check.wallet.invalid.req.climasterpubkey:%v", uid, cliMasterPubKey)
		return fmt.Errorf("wdb.uid[%v].req.masterpubkey[%v].invalid", uid, cliMasterPubKey)
	}
	return nil
}

// Open -- used to load all the wallets who in the disk to the cache.
func (wdb *WalletDB) Open(dir string) error {
	store := wdb.store
	syncer := wdb.syncer

	if err := store.Open(dir); err != nil {
		return err
	}
	syncer.Start()
	return nil
}

// Close -- used to close the db.
func (wdb *WalletDB) Close() {
	wdb.syncer.Stop()
}

// OpenWalletByUID -- used to open(or create if not exists) a wallet file.
func (wdb *WalletDB) OpenWalletByUID(uid string, cliMasterPubKey string) (*Wallet, error) {
	net := wdb.net
	store := wdb.store

	wallet := store.Get(uid)
	if wallet == nil {
		masterKey, err := bip32.NewHDKeyRand()
		if err != nil {
			return nil, err
		}
		svrMasterPrvKey := masterKey.ToString(net)
		wallet = &Wallet{
			UID:             uid,
			Address:         make(map[string]*Address),
			CliMasterPubKey: cliMasterPubKey,
			SvrMasterPrvKey: svrMasterPrvKey,
		}
		return wallet, store.Write(wallet)
	}

	// Check.
	if err := wdb.check(uid, wallet.CliMasterPubKey, cliMasterPubKey); err != nil {
		return nil, err
	}
	return wallet, nil
}

// NewAddressByUID -- used to generate new address of this uid.
func (wdb *WalletDB) NewAddressByUID(uid string, cliMasterPubKey string) (*Address, error) {
	net := wdb.net
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.newaddress.uid[%v].cant.found", uid)
	}

	// Check.
	if err := wdb.check(uid, wallet.CliMasterPubKey, cliMasterPubKey); err != nil {
		return nil, err
	}

	// New address.
	wallet.Lock()
	pos := wallet.LastPos
	addr, err := createSharedAddress(pos, wallet.SvrMasterPrvKey, wallet.CliMasterPubKey, net)
	if err != nil {
		wallet.Unlock()
		return nil, err
	}

	address := &Address{
		Pos:     pos,
		Address: addr,
	}
	wallet.Address[addr] = address
	wallet.LastPos++
	wallet.Unlock()

	// Write to db.
	store.Write(wallet)
	return address, nil
}

// Balance --used to return balance of the wallet.
func (wdb *WalletDB) Balance(uid string) (*Balance, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.balance.uid[%v].cant.found", uid)
	}
	wallet.Lock()
	defer wallet.Unlock()
	balance := &Balance{}
	for _, addr := range wallet.Address {
		balance.AllBalance += addr.Balance.AllBalance
		balance.UnconfirmedBalance += addr.Balance.UnconfirmedBalance
	}
	return balance, nil
}
