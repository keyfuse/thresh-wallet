// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

// BackupData --
type BackupData struct {
	DeviceID         string `json:"device_id"`
	CloudService     string `json:"cloud_service"`
	CloudServiceName string `json:"cloud_service_name"`
	EncryptedKey     string `json:"encrypted_key"`
	EncryptedKeyHash string `json:"encrypted_key_hash"`
	BackupTime       uint32 `json:"backup_time"`
}

// Backup --
type Backup struct {
	Time  uint32     `json:"time"`
	Email string     `json:"email"`
	Data  BackupData `json:"data"`
}
