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

func TestVcode(t *testing.T) {
	conf := DefaultConfig()
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))

	vcode := NewVcode(log, conf)

	// Add.
	{
		vcode.Add("13888888888", "88886666")
	}

	// UID Error.
	{
		err := vcode.Check("10087", "88886666")
		assert.NotNil(t, err)
	}

	// Code error.
	{
		err := vcode.Check("13888888888", "8886666")
		assert.NotNil(t, err)
	}

	// OK.
	{
		err := vcode.Check("13888888888", "88886666")
		assert.Nil(t, err)
	}
}
