// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"sync"
)

// UTXO --
type UTXO struct {
	Pos          uint32 `json:"pos"`
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	Value        uint64 `json:"value"`
	Address      string `json:"address"`
	Confirmed    bool   `json:"confirmed"`
	SvrPubKey    string `json:"svrpubkey"`
	Scriptpubkey string `json:"Scriptpubkey"`
}

// Balance --
type Balance struct {
	AllBalance         uint64 `json:"all_balance"`
	UnconfirmedBalance uint64 `json:"unconfirmed_balance"`
}

// Unspent --
type Unspent struct {
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	Value        uint64 `json:"value"`
	Confirmed    bool   `json:"confirmed"`
	BlockTime    uint32 `json:"block_time"`
	BlockHeight  uint32 `json:"block_height"`
	Scriptpubkey string `json:"Scriptpubkey"`
}

// Address --
type Address struct {
	mu       sync.Mutex
	Pos      uint32    `json:"pos"`
	Address  string    `json:"address"`
	Unspents []Unspent `json:"unspents"`
	Balance  Balance   `json:"balance"`
}

type Backup struct {
	Time     uint32 `json:"time"`
	Email    string `json:"email"`
	EncKey   string `json:"enckey"`
	DeviceID string `json:"device_id"`
}

// Wallet --
type Wallet struct {
	mu              sync.Mutex
	UID             string              `json:"uid"`
	DID             string              `json:"did"`
	Backup          *Backup             `json:"backup"`
	LastPos         uint32              `json:"lastpos"`
	Address         map[string]*Address `json:"address"`
	SvrMasterPrvKey string              `json:"svrmasterprvkey"`
	CliMasterPubKey string              `json:"climasterpubkey"`
}

// NewWallet -- creates new Wallet.
func NewWallet() *Wallet {
	return &Wallet{
		Address: make(map[string]*Address),
	}
}

// Lock -- used to lock the wallet entry for thread-safe purposes.
func (w *Wallet) Lock() {
	w.mu.Lock()
}

// Unlock -- used to unlock the wallet entry for thread-safe purposes.
func (w *Wallet) Unlock() {
	w.mu.Unlock()
}

// Addresses -- used to returns all the address of the wallet.
func (w *Wallet) Addresses() []string {
	var addrs []string

	w.mu.Lock()
	defer w.mu.Unlock()
	for addr := range w.Address {
		addrs = append(addrs, addr)
	}
	return addrs
}

// UpdateUnspents -- update the address balance/unspent which fetchs from the chain.
func (w *Wallet) UpdateUnspents(addr string, unspents []Unspent) {
	w.mu.Lock()
	address := w.Address[addr]
	w.mu.Unlock()

	address.mu.Lock()
	defer address.mu.Unlock()

	var balance, unconfirmedBalance uint64
	for _, unspent := range unspents {
		if !unspent.Confirmed {
			unconfirmedBalance += unspent.Value
		}
		balance += unspent.Value
	}
	address.Unspents = unspents
	address.Balance.AllBalance = balance
	address.Balance.UnconfirmedBalance = unconfirmedBalance
}
