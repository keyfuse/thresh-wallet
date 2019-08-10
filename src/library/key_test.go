// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMasterKey(t *testing.T) {
	body := NewMasterPrvKey("testnet")
	t.Logf("body:%+v", body)

	rsp := &MasterPrvKeyResponse{}
	unmarshal(body, rsp)
	assert.Equal(t, 200, rsp.Code)
}
