// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"os"
	"sync"
	"testing"

	"xlog"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
)

func TestWalletDB(t *testing.T) {
	defer leaktest.Check(t)()

	uid := "U002"
	conf := MockConfig()
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	wdb := NewWalletDB(log, conf)
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
		_, err := wdb.OpenWalletByUID(uid, mockCliMasterPubKey)
		assert.Nil(t, err)
	}

	// New address.
	{
		for i := 0; i < 10; i++ {
			addr, err := wdb.NewAddressByUID(uid, mockCliMasterPubKey)
			assert.Nil(t, err)
			t.Logf("addr:%+v", addr)
		}
	}

	// Thread-Safe check.
	{
		var wg sync.WaitGroup

		for i := 0; i < 30; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				addr, err := wdb.NewAddressByUID(uid, mockCliMasterPubKey)
				assert.Nil(t, err)
				t.Logf("addr:%+v", addr)
			}()
		}
		wg.Wait()
	}
}
