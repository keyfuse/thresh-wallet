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

func TestInfoHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	{

		httpRsp, err := proto.NewRequest().Get(ts.URL + "/api/server/info")
		assert.Nil(t, err)
		rsp := &proto.ServerInfoResponse{}
		httpRsp.Json(rsp)
		t.Logf("rsp:%+v", rsp)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
