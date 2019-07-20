// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"os"
	"testing"

	"xlog"

	"github.com/tokublock/tokucore/network"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
)

func TestWalletStore(t *testing.T) {
	defer leaktest.Check(t)()

	wallet := NewWallet()
	wallet.net = network.TestNet
	wallet.UID = mockUID
	wallet.SvrMasterPrvKey = mockSvrMasterPrvKey
	wallet.CliMasterPubKey = mockCliMasterPubKey

	dir := "/tmp/tss"
	os.RemoveAll(dir)

	conf := MockConfig()
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	wstore := NewWalletStore(log, conf)

	// Open.
	{
		err := wstore.Open(dir)
		assert.Nil(t, err)
	}

	// Write.
	{
		err := wstore.Write(wallet)
		assert.Nil(t, err)
	}

	// Read.
	{
		path := "/tmp/tss/10086.json"
		got, err := wstore.Read(path)
		assert.Nil(t, err)
		assert.Equal(t, wallet, got)
	}

	// re-Open.
	{
		wstore2 := NewWalletStore(log, conf)
		err := wstore2.Open(dir)
		assert.Nil(t, err)
	}
}
