// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

type Backup struct {
	Time             int64  `json:"time"`
	Email            string `json:"email"`
	DeviceID         string `json:"deviceid"`
	CloudService     string `json:"cloud_service"`
	EncryptedPrvKey  string `json:"encrypted_prvkey"`
	EncryptionPubKey string `json:"encryption_pubkey"`
}
