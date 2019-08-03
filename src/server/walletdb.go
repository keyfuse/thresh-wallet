// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"sync"
	"time"

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

// CreateWallet -- used to create a wallet file.
func (wdb *WalletDB) CreateWallet(uid string, cliMasterPubKey string) error {
	net := wdb.net
	store := wdb.store

	wallet := store.Get(uid)
	if wallet == nil {
		masterKey, err := bip32.NewHDKeyRand()
		if err != nil {
			return err
		}
		svrMasterPrvKey := masterKey.ToString(net)
		wallet = &Wallet{
			net:             net,
			UID:             uid,
			Address:         make(map[string]*Address),
			CliMasterPubKey: cliMasterPubKey,
			SvrMasterPrvKey: svrMasterPrvKey,
		}
		return store.Write(wallet)
	} else {
		return fmt.Errorf("wdb.wallet[%v, %v].create.error:wallet.exists", uid, cliMasterPubKey)
	}
}

// NewAddress -- used to generate new address of this uid.
func (wdb *WalletDB) NewAddress(uid string, typ string) (*Address, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.newaddress.uid[%v].cant.found", uid)
	}

	address, err := wallet.NewAddress(typ)
	if err != nil {
		return nil, err
	}

	// Write to db.
	if err := store.Write(wallet); err != nil {
		return nil, err
	}
	return address, nil
}

func (wdb *WalletDB) MasterPrvKey(uid string) (string, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return "", fmt.Errorf("wdb.master.prvkey.uid[%v].cant.found", uid)
	}

	return wallet.SvrMasterPrvKey, nil
}

// Wallet -- used to get the wallet.
func (wdb *WalletDB) Wallet(uid string) *Wallet {
	store := wdb.store
	return store.Get(uid)
}

// Balance --used to return balance of the wallet.
func (wdb *WalletDB) Balance(uid string) (*Balance, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.balance.uid[%v].cant.found", uid)
	}
	return wallet.Balance(), nil
}

// Unspents -- used to return unspent which all the value upper than the amount.
func (wdb *WalletDB) Unspents(uid string, amount uint64) ([]UTXO, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.unspents.uid[%v].cant.found", uid)
	}
	return wallet.Unspents(amount)
}

// Txs -- used to returns tx list.
func (wdb *WalletDB) Txs(uid string, offset int, limit int) ([]Tx, error) {
	var ret []Tx
	store := wdb.store
	chain := wdb.chain

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.txs.uid[%v].cant.found", uid)
	}
	txs := wallet.Txs(offset, limit)
	for _, tx := range txs {
		tx.Link = fmt.Sprintf(chain.GetTxLink(), tx.Txid)
		ret = append(ret, tx)
	}
	return ret, nil
}

// Addresses -- used to get address list.
func (wdb *WalletDB) Addresses(uid string, offset int, limit int) ([]AddressPos, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.addresses.uid[%v].cant.found", uid)
	}
	ret := wallet.AddressPoss(offset, limit)
	return ret, nil
}

// SendFees -- returns the fee info for this send.
func (wdb *WalletDB) SendFees(uid string, priority string, sendAmount uint64) (*SendFees, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return nil, fmt.Errorf("wdb.send.fees.uid[%v].cant.found", uid)
	}

	feesperkb := store.FeesPerKB(priority)
	return wallet.SendFees(sendAmount, feesperkb)
}

func (wdb *WalletDB) StoreBackup(uid string, email string, did string, cloudService string, encryptedPrvKey string, encryptionPubKey string) error {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return fmt.Errorf("wdb.store.backup.uid[%v].cant.found", uid)
	}
	wallet.Lock()
	wallet.Backup.Email = email
	wallet.Backup.DeviceID = did
	wallet.Backup.CloudService = cloudService
	wallet.Backup.EncryptedPrvKey = encryptedPrvKey
	wallet.Backup.EncryptionPubKey = encryptionPubKey
	wallet.Backup.Time = time.Now().Unix()
	wallet.Unlock()
	return store.Write(wallet)
}

func (wdb *WalletDB) GetBackup(uid string) (Backup, error) {
	store := wdb.store

	// Get wallet.
	wallet := store.Get(uid)
	if wallet == nil {
		return Backup{}, fmt.Errorf("wdb.backup.pubkey.uid[%v].cant.found", uid)
	}
	return wallet.Backup, nil
}
