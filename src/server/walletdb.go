// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"sort"
	"sync"

	"xlog"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore/bip32"
)

// WalletDB --
type WalletDB struct {
	mu     sync.Mutex
	log    *xlog.Log
	conf   *Config
	net    *network.Network
	chain  Chain
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

	chain := NewChainProxy(log, conf)
	store := NewWalletStore(log, conf)
	syncer := NewWalletSyncer(log, conf, chain, store)
	return &WalletDB{
		log:    log,
		net:    net,
		conf:   conf,
		chain:  chain,
		store:  store,
		syncer: syncer,
	}
}

func (wdb *WalletDB) setChain(chain Chain) {
	wdb.mu.Lock()
	defer wdb.mu.Unlock()

	log := wdb.log
	conf := wdb.conf
	store := wdb.store

	syncer := wdb.syncer
	syncer.Stop()

	// Set new syncer.
	newsyncer := NewWalletSyncer(log, conf, chain, store)
	wdb.syncer = newsyncer
	newsyncer.Start()
	wdb.chain = chain
}

// Open -- used to load all the wallets who in the disk to the cache.
func (wdb *WalletDB) Open(dir string) error {
	wdb.mu.Lock()
	defer wdb.mu.Unlock()

	if err := wdb.store.Open(dir); err != nil {
		return err
	}
	wdb.syncer.Start()
	return nil
}

// Close -- used to close the db.
func (wdb *WalletDB) Close() {
	wdb.mu.Lock()
	defer wdb.mu.Unlock()
	wdb.syncer.Stop()
}

func (wdb *WalletDB) check(uid string, walletMasterPubKey string, cliMasterPubKey string) error {
	log := wdb.log

	if walletMasterPubKey != cliMasterPubKey {
		log.Error("wdb.openwallet[%v].check.wallet.invalid.req.climasterpubkey:%v, svr.climasterpubkey:%v", uid, cliMasterPubKey, walletMasterPubKey)
		return fmt.Errorf("wdb.uid[%v].req.masterpubkey[%v].invalid", uid, cliMasterPubKey)
	}
	return nil
}

// OpenUIDWallet -- used to open(or create if not exists) a wallet file.
func (wdb *WalletDB) OpenUIDWallet(uid string, cliMasterPubKey string) (*Wallet, error) {
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

// NewAddress -- used to generate new address of this uid.
func (wdb *WalletDB) NewAddress(uid string, cliMasterPubKey string) (*Address, error) {
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

func (wdb *WalletDB) MasterPrvKey(uid string, cliMasterPubKey string) (string, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return "", fmt.Errorf("wdb.master.prvkey.uid[%v].cant.found", uid)
	}

	// Check.
	if err := wdb.check(uid, wallet.CliMasterPubKey, cliMasterPubKey); err != nil {
		return "", err
	}
	return wallet.SvrMasterPrvKey, nil
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

// UnspentsByAmount -- used to return unspent which all the value upper than the amount.
func (wdb *WalletDB) Unspents(uid string, amount uint64) ([]UTXO, error) {
	var rsp []UTXO
	var utxos []UTXO
	var thresh uint64
	var balance uint64

	net := wdb.net
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.unspentsbyamount[%v].cant.found", uid)
	}
	wallet.Lock()
	defer wallet.Unlock()

	for _, addr := range wallet.Address {
		for _, unspent := range addr.Unspents {
			svrpubkey, err := createSvrChildPubKey(addr.Pos, wallet.SvrMasterPrvKey, net)
			if err != nil {
				return nil, err
			}
			utxos = append(utxos, UTXO{
				Pos:          addr.Pos,
				Txid:         unspent.Txid,
				Vout:         unspent.Vout,
				Value:        unspent.Value,
				Address:      addr.Address,
				Confirmed:    unspent.Confirmed,
				SvrPubKey:    svrpubkey,
				Scriptpubkey: unspent.Scriptpubkey,
			})
		}
		balance += addr.Balance.AllBalance
	}

	// Check.
	if balance <= amount {
		return nil, fmt.Errorf("wdb.unpsentsbyamount[%v].suffient.req.amount[%v].allbalance[%v]", uid, amount, balance)
	}

	// Sort by value desc.
	sort.Slice(utxos, func(i, j int) bool { return utxos[i].Value > utxos[j].Value })

	// Patch.
	for _, utxo := range utxos {
		thresh += utxo.Value
		rsp = append(rsp, utxo)
		if thresh > amount {
			break
		}
	}
	return rsp, nil
}
