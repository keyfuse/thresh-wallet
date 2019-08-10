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

func TestAPIHello(t *testing.T) {
	body := APIHello("jack")
	rsp := &HelloResponse{}
	unmarshal(body, rsp)
	assert.Equal(t, 200, rsp.Code)
}
