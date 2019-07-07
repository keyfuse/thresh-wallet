// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"fmt"
	"net/http"

	"proto"
)

// VCodeRequest --
type VcodeResponse struct {
	Status
}

// APIGetVCode -- the vcode api.
func APIGetVCode(url string, uid string) string {
	rsp := &VcodeResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/vcode", url)

	req := &proto.VCodeRequest{
		UID: uid,
	}

	httpRsp, err := proto.NewRequest().Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.Code = httpRsp.StatusCode()
	return marshal(rsp)
}

// TokenResponse --
type TokenResponse struct {
	Status
	Token string `json:"token"`
}

// APIGetToken -- get token api.
func APIGetToken(url string, uid string, vcode string, masterPubKey string) string {
	rsp := &TokenResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/token", url)

	req := &proto.TokenRequest{
		UID:          uid,
		VCode:        vcode,
		MasterPubKey: masterPubKey,
	}
	httpRsp, err := proto.NewRequest().Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	token := &proto.TokenResponse{}
	if err := httpRsp.Json(token); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.Token = token.Token
	return marshal(rsp)
}
