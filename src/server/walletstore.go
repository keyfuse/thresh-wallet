// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"xlog"

	"github.com/tokublock/tokucore/network"
)

// WalletStore --
type WalletStore struct {
	mu      sync.Mutex
	dir     string
	log     *xlog.Log
	conf    *Config
	net     *network.Network
	fees    map[string]float32
	wallets map[string]*Wallet
	tickers map[string]Ticker
}

// NewWalletStore -- creates new WalletStore.
func NewWalletStore(log *xlog.Log, conf *Config) *WalletStore {
	var net *network.Network
	switch conf.ChainNet {
	case testnet:
		net = network.TestNet
	case mainnet:
		net = network.MainNet
	}
	return &WalletStore{
		log:     log,
		conf:    conf,
		net:     net,
		fees:    make(map[string]float32),
		wallets: make(map[string]*Wallet),
		tickers: make(map[string]Ticker),
	}
}

// Open -- used to open the dir database.
func (s *WalletStore) Open(dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log := s.log
	s.dir = dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if x := os.MkdirAll(dir, os.ModePerm); x != nil {
			return x
		}
		return nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", dir, file.Name())
		wallet, err := s.Read(path)
		if err != nil {
			return err
		}
		s.wallets[wallet.UID] = wallet
		log.Info("wallet.store.load[%s/%v]", dir, wallet.UID)
	}
	return nil
}

// Write -- write wallet to the file.
// Thread-Safe.
func (s *WalletStore) Write(wallet *Wallet) error {
	dir := s.dir
	uid := wallet.UID

	s.mu.Lock()
	w, ok := s.wallets[uid]
	if !ok {
		s.wallets[uid] = wallet
	} else if (w.CliMasterPubKey != wallet.CliMasterPubKey) || (w.SvrMasterPrvKey != wallet.SvrMasterPrvKey) {
		return fmt.Errorf("storage.write.data.race.uid[%v]", uid)
	}
	s.mu.Unlock()

	wallet.Lock()
	defer wallet.Unlock()
	file := fmt.Sprintf("%s/%s.json", dir, uid)
	datas, err := json.MarshalIndent(wallet, "", " ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(file, datas, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Read -- reads a wallet from the file.
func (s *WalletStore) Read(path string) (*Wallet, error) {
	wallet := NewWallet()

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, wallet); err != nil {
		return nil, err
	}
	if wallet.Address == nil {
		wallet.Address = make(map[string]*Address)
	}
	wallet.net = s.net
	return wallet, nil
}

// Get -- used to get a wallet from the wallets list.
// Returns nil if not exists.
func (s *WalletStore) Get(uid string) *Wallet {
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, ok := s.wallets[uid]
	if !ok {
		return nil
	}
	return wallet
}

// AllUID -- used to clone all the wallet uids.
func (s *WalletStore) AllUID() []string {
	var uids []string

	s.mu.Lock()
	defer s.mu.Unlock()
	for uid := range s.wallets {
		uids = append(uids, uid)
	}
	return uids
}

func (s *WalletStore) updateFees(fees map[string]float32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fees = fees
}

func (s *WalletStore) updateTickers(tickers map[string]Ticker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tickers = tickers
}

func (s *WalletStore) getTicker(code string) (Ticker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ticker, ok := s.tickers[code]
	if !ok {
		return ticker, fmt.Errorf("wallet.store.get.ticker.code[%v].cant.found", code)
	}
	return ticker, nil
}

// FeesPerKB -- used to return the fees.
func (s *WalletStore) FeesPerKB(priority string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	fees := 1000
	switch strings.ToUpper(priority) {
	case "FAST":
		if v, ok := s.fees["2"]; ok {
			fees = int(v * 1000)
		}
	case "NORMAL":
		if v, ok := s.fees["4"]; ok {
			fees = int(v * 1000)
		}
	case "SLOW":
		if v, ok := s.fees["6"]; ok {
			fees = int(v * 1000)
		}
	default:
		if v, ok := s.fees["4"]; ok {
			fees = int(v * 1000)
		}
	}
	return fees
}
