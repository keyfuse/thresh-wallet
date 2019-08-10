// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"net/http"
)

// HelloResponse --
type HelloResponse struct {
	Status
	Hello string `json:"hello"`
}

// APIHello -- test for IOS/Android static.
func APIHello(msg string) string {
	rsp := &HelloResponse{}
	rsp.Code = http.StatusOK
	rsp.Hello = "yes, hello " + msg
	return marshal(rsp)
}
