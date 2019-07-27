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

	"github.com/stretchr/testify/assert"
)

func TestWalletSyncer(t *testing.T) {
	uid := "U002"
	conf := MockConfig()
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	wdb := NewWalletDB(log, conf)
	wdb.setChain(newMockChain(log))
	defer wdb.Close()

	// Open.
	{
		dir := "/tmp/tss"
		os.RemoveAll(dir)

		err := wdb.Open(dir)
		assert.Nil(t, err)
	}

	// Get.
	{
		err := wdb.CreateWallet(uid, mockCliMasterPubKey)
		assert.Nil(t, err)
	}

	// New address.
	{
		for i := 0; i < 3; i++ {
			addr, err := wdb.NewAddress(uid, "")
			assert.Nil(t, err)
			t.Logf("addr:%+v", addr)
		}
	}

	syncer := wdb.syncer
	syncer.Sync()
}
