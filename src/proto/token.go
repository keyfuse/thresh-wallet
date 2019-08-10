// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package proto

// VCodeRequest --
type VCodeRequest struct {
	UID string `json:"uid"`
}

// TokenRequest --
type TokenRequest struct {
	UID      string `json:"uid"`
	VCode    string `json:"vcode"`
	DeviceID string `json:"deviceid"`
}

// TokenResponse --
type TokenResponse struct {
	Token string `json:"token"`
}
