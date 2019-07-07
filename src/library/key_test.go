// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"testing"
)

func TestNewMasterKey(t *testing.T) {
	body := NewMasterPrvKey("testnet")
	t.Logf("body:%+v", body)

	mkrsp := &MasterPrvKeyResponse{}
	unmarshal(body, mkrsp)

	body = GetMasterPubKey("testnet", mkrsp.MasterPrvKey)
	t.Logf("body:%+v", body)
}
