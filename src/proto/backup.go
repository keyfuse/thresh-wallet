// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package proto

// BackupCheckRequest --
type BackupCheckRequest struct {
}

// BackupCheckResponse --
type BackupCheckResponse struct {
}

// BackupVCodeRequest --
type BackupVCodeRequest struct {
}

// BackupVCodeResponse --
type BackupVCodeResponse struct {
	VCode string `json:"vcode"`
}

// BackupStoreRequest --
type BackupStoreRequest struct {
	Email            string `json:"email"`
	VCode            string `json:"vcode"`
	DeviceID         string `json:"deviceid"`
	Signature        string `json:"signature"`
	CloudService     string `json:"cloud_service"`
	EncryptedPrvKey  string `json:"encrypted_prvkey"`
	EncryptionPubKey string `json:"encryption_pubkey"`
}

// BackupStoreResponse --
type BackupStoreResponse struct {
}

// BackupRestoreRequest --
type BackupRestoreRequest struct {
	VCode     string `json:"vcode"`
	Signature string `json:"signature"`
}

// BackupStoreResponse --
type BackupRestoreResponse struct {
	Time            int64  `json:"time"`
	EncryptedPrvKey string `json:"encrypted_prvkey"`
}

// BackupVerifyRequest --
type BackupVerifyRequest struct {
	EncryptionPubKeyHash string `json:"encryption_pubkey_hash"`
}

// BackupVerifyResponse --
type BackupVerifyResponse struct {
	VerifyPassed    bool  `json:"verify_passed"`
	VerifyTimestamp int64 `json:"verify_timestamp"`
}
