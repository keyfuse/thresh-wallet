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

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
)

func TestWalletStore(t *testing.T) {
	defer leaktest.Check(t)()

	wallet := NewWallet()
	wallet.UID = "U001"
	wallet.SvrMasterPrvKey = "tprv8ZgxMBicQKsPdZaiD1bZC55UcWpif1Nk9SD4iqjemtPFNRYcXMFRGiWGyGejLWJpXqXffi9zdiYkDqtgF3Gn2ShmbhQYGMsCm3Q8jGPFDLR"
	wallet.CliMasterPubKey = "tpubD6NzVbkrYhZ4X7Cn1qGQ7XReumN4yFvgP3ms8dPtTiLD7wpP95cqmbaAkk5WSZaSrgpgtmPQhpNGmkxVRezP3WN486xEddsWHU22a6F7yJZ"

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
		path := "/tmp/tss/U001.json"
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
