// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package proto

// VCodeRequest --
type VCodeRequest struct {
	UID string `json:"uid"`
}

// TokenRequest --
type TokenRequest struct {
	UID          string `json:"uid"`
	VCode        string `json:"vcode"`
	DeviceID     string `json:"deviceid"`
	MasterPubKey string `json:"masterpubkey"`
}

// TokenResponse --
type TokenResponse struct {
	Token string `json:"token"`
}
