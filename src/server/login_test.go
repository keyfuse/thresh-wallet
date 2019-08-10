// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginType(t *testing.T) {
	tests := []struct {
		uid string
		res LoginType
	}{
		{uid: "13666666666", res: Mobile},
		{uid: "(+86)13666666666", res: Mobile},
		{uid: "a@b", res: Email},
		{uid: "a@b.com", res: Email},
	}

	for _, test := range tests {
		typ := loginType(test.uid)
		assert.Equal(t, test.res, typ)
	}
}
