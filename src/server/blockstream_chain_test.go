// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"testing"

	"xlog"

	"github.com/stretchr/testify/assert"
)

func TestBlockstreamChain(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	conf := MockConfig()

	chain := NewBlockstreamChain(log, conf)

	unspent, err := chain.GetUTXO("mkTcdRZAtJLYggv6zspQg6GP22JKS3DSXo")
	assert.Nil(t, err)
	assert.NotNil(t, unspent)
}
