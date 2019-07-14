// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"xlog"
)

// WalletStore --
type WalletStore struct {
	mu      sync.Mutex
	dir     string
	log     *xlog.Log
	conf    *Config
	wallets map[string]*Wallet
}

// NewWalletStore -- creates new WalletStore.
func NewWalletStore(log *xlog.Log, conf *Config) *WalletStore {
	return &WalletStore{
		log:     log,
		conf:    conf,
		wallets: make(map[string]*Wallet),
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
		log.Info("wallet.store.load[%v]", wallet.UID)
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
