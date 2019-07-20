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

	assert.NotNil(t, chain)
}
