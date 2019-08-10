// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package proto

// ServerInfoResponse --
type ServerInfoResponse struct {
	ChainNet    string `json:"chainnet"`
	ServerTime  int64  `json:"server_time"`
	EnableVCode bool   `json:"enable_vcode"`
}
