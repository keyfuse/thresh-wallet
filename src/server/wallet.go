// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"fmt"
	"sort"
	"sync"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore"
)

// Ticker --
type Ticker struct {
	One5M  float64 `json:"15m"`
	Last   float64 `json:"last"`
	Buy    float64 `json:"buy"`
	Sell   float64 `json:"sell"`
	Symbol string  `json:"symbol"`
}

// SendFees --
type SendFees struct {
	Fees          uint64 `json:"fees"`
	TotalValue    uint64 `json:"total_value"`
	SendableValue uint64 `json:"sendable_value"`
}

// Tx --
type Tx struct {
	Txid        string `json:"txid"`
	Fee         int64  `json:"fee"`
	Data        string `json:"data"`
	Link        string `json:"link"`
	Value       int64  `json:"value"`
	Confirmed   bool   `json:"confirmed"`
	BlockTime   int64  `json:"block_time"`
	BlockHeight int64  `json:"block_height"`
}

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
	TotalBalance       uint64 `json:"total_balance"`
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

// AddressPos --
type AddressPos struct {
	Pos     uint32 `json:"pos"`
	Address string `json:"address"`
}

// Address --
type Address struct {
	mu       sync.Mutex
	Pos      uint32    `json:"pos"`
	Address  string    `json:"address"`
	Balance  Balance   `json:"balance"`
	Txs      []Tx      `json:"txs"`
	Unspents []Unspent `json:"unspents"`
}

// Wallet --
type Wallet struct {
	mu              sync.Mutex
	net             *network.Network
	UID             string              `json:"uid"`
	DID             string              `json:"did"`
	Backup          Backup              `json:"backup"`
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
func (w *Wallet) Addresses() []AddressPos {
	w.Lock()
	defer w.Unlock()

	var addrs []AddressPos
	for _, addr := range w.Address {
		addrs = append(addrs, AddressPos{
			Address: addr.Address,
			Pos:     addr.Pos,
		})
	}
	sort.Slice(addrs, func(i, j int) bool {
		return addrs[i].Pos > addrs[j].Pos
	})
	return addrs
}

// NewAddress -- used to generate new address.
func (w *Wallet) NewAddress(typ string) (*Address, error) {
	net := w.net

	// New address.
	w.Lock()
	defer w.Unlock()

	pos := w.LastPos
	addr, err := createSharedAddress(pos, w.SvrMasterPrvKey, w.CliMasterPubKey, net, typ)
	if err != nil {
		return nil, err
	}

	address := &Address{
		Pos:     pos,
		Address: addr,
	}
	w.Address[addr] = address
	w.LastPos++

	return address, nil
}

// UpdateUnspents -- update the address balance/unspent which fetchs from the chain.
func (w *Wallet) UpdateUnspents(addr string, unspents []Unspent) {
	w.Lock()
	address := w.Address[addr]
	w.Unlock()

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
	address.Balance.TotalBalance = balance
	address.Balance.UnconfirmedBalance = unconfirmedBalance
}

// UpdateTxs -- update the tx fetchs from the chain.
func (w *Wallet) UpdateTxs(addr string, txs []Tx) {
	w.Lock()
	address := w.Address[addr]
	w.Unlock()

	address.mu.Lock()
	defer address.mu.Unlock()
	address.Txs = txs
}

// Balance --used to return balance of the wallet.
func (w *Wallet) Balance() *Balance {
	w.Lock()
	defer w.Unlock()
	balance := &Balance{}
	for _, addr := range w.Address {
		balance.TotalBalance += addr.Balance.TotalBalance
		balance.UnconfirmedBalance += addr.Balance.UnconfirmedBalance
	}
	return balance
}

// Unspents -- used to return unspent which all the value upper than the amount.
func (w *Wallet) Unspents(sendAmount uint64) ([]UTXO, error) {
	var rsp []UTXO
	var utxos []UTXO
	var thresh uint64
	var balance uint64
	net := w.net

	w.Lock()
	defer w.Unlock()

	for _, addr := range w.Address {
		for _, unspent := range addr.Unspents {
			svrpubkey, err := createSvrChildPubKey(addr.Pos, w.SvrMasterPrvKey, net)
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
		balance += addr.Balance.TotalBalance
	}

	// Check.
	if balance < sendAmount {
		return nil, fmt.Errorf("unpsents.suffient.req.amount[%v].allbalance[%v]", sendAmount, balance)
	}

	// Sort by value desc.
	sort.Slice(utxos, func(i, j int) bool { return utxos[i].Value > utxos[j].Value })

	// Patch.
	for _, utxo := range utxos {
		thresh += utxo.Value
		rsp = append(rsp, utxo)
		if thresh >= sendAmount {
			break
		}
	}
	return rsp, nil
}

// Txs -- used to return the txs starts from offset to offset+limit.
func (w *Wallet) Txs(offset int, limit int) []Tx {
	var txs []Tx
	var confirmedtxs []Tx
	var unconfirmedtxs []Tx

	w.Lock()
	for _, addr := range w.Address {
		txs = append(txs, addr.Txs...)
	}
	w.Unlock()

	for _, tx := range txs {
		if tx.BlockHeight == 0 {
			unconfirmedtxs = append(unconfirmedtxs, tx)
		} else {
			confirmedtxs = append(confirmedtxs, tx)
		}
	}

	// Sort txs.
	sort.Slice(confirmedtxs, func(i, j int) bool {
		return confirmedtxs[i].BlockTime > confirmedtxs[j].BlockTime
	})
	txs = txs[:0]
	txs = append(txs, unconfirmedtxs...)
	txs = append(txs, confirmedtxs...)

	size := len(txs)
	if offset >= size {
		return nil
	}
	if (offset + limit) > size {
		return txs[offset:]
	} else {
		return txs[offset : offset+limit]
	}
}

// AddressPoss -- used to return the AddressPoss from offset to offset+limit.
func (w *Wallet) AddressPoss(offset int, limit int) []AddressPos {
	addrs := w.Addresses()

	size := len(addrs)
	if offset >= size {
		return nil
	}
	if (offset + limit) > size {
		return addrs[offset:]
	} else {
		return addrs[offset : offset+limit]
	}
}

// SendFees -- used to get the send fees by send amount.
func (w *Wallet) SendFees(sendValue uint64, feesPerKB int) (*SendFees, error) {
	unspents, err := w.Unspents(sendValue)
	if err != nil {
		return nil, err
	}

	totalValue := w.Balance().TotalBalance
	estsize := xcore.EstimateNormalSize(len(unspents), 1+1)
	fees := uint64((estsize * int64(feesPerKB)) / 1000)

	if fees >= totalValue {
		return nil, fmt.Errorf("balace[%v].is.smaller.than.fees[%v]", totalValue, fees)
	}

	sendableValue := sendValue
	// Send all case.
	if (sendableValue + fees) > totalValue {
		sendableValue = totalValue - fees
	}
	return &SendFees{
		Fees:          fees,
		TotalValue:    totalValue,
		SendableValue: sendableValue,
	}, nil
}
