// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"testing"

	"proto"

	"github.com/stretchr/testify/assert"
)

func TestLoginVCodeHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// VCode.
	{
		req := &proto.VCodeRequest{
			UID: mockUID,
		}
		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/login/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestLoginTokenHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// OK.
	{
		req := &proto.TokenRequest{
			UID: mockUID,
		}

		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/login/token", req)
		assert.Nil(t, err)
		rsp := &proto.TokenResponse{}
		httpRsp.Json(rsp)
		t.Logf("rsp:%+v", rsp)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
