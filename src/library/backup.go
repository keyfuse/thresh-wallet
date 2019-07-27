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

// WalletBackupResponse --
type WalletBackupResponse struct {
	Status
}

// APIWalletBackup --
func APIWalletBackup(url string, token string, deviceID string, cloudService string, rsaPrvKey string, masterPrvKey string) string {
	var vcode string
	var signature string
	var encryptedPrvKey string
	var encryptionPubKey string

	rsp := &WalletBackupResponse{}
	rsp.Code = http.StatusOK

	// vcode.
	{
		path := fmt.Sprintf("%s/api/backup/vcode", url)
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		ret := &proto.BackupVCodeResponse{}
		if err := httpRsp.Json(ret); err != nil {
			rsp.Code = httpRsp.StatusCode()
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		vcode = ret.VCode
	}

	// Signature.
	{
		body := RSASign(vcode, rsaPrvKey)
		ret := RSASignResponse{}
		if err := unmarshal(body, &ret); err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		signature = ret.Signature
	}

	// Encrypt.
	{
		body := RSAEncrypt(masterPrvKey, rsaPrvKey)
		ret := RSAEncryptResponse{}
		if err := unmarshal(body, &ret); err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		encryptedPrvKey = ret.CipherText
	}

	// RSA PubKey.
	{
		body := GetRSAPubKey(rsaPrvKey)
		ret := RSAPubKeyResponse{}
		if err := unmarshal(body, &ret); err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		encryptionPubKey = ret.PubKey
	}

	// Backup.
	{
		path := fmt.Sprintf("%s/api/backup/store", url)
		req := &proto.BackupStoreRequest{
			VCode:            vcode,
			DeviceID:         deviceID,
			Signature:        signature,
			CloudService:     cloudService,
			EncryptedPrvKey:  encryptedPrvKey,
			EncryptionPubKey: encryptionPubKey,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		ret := &proto.BackupStoreResponse{}
		if err := httpRsp.Json(ret); err != nil {
			rsp.Code = httpRsp.StatusCode()
			rsp.Message = err.Error()
			return marshal(rsp)
		}
	}
	return marshal(rsp)
}

// WalletRestoreResponse --
type WalletRestoreResponse struct {
	Status
	Time         int64  `json:"time"`
	MasterPrvKey string `json:"masterprvkey"`
}

// APIWalletRestore -- used to restore the backup from the server.
func APIWalletRestore(url string, token string, rsaPrvKey string) string {
	var vcode string
	var signature string
	var encryptedPrvKey string

	rsp := &WalletRestoreResponse{}
	rsp.Code = http.StatusOK

	// vcode.
	{
		path := fmt.Sprintf("%s/api/backup/vcode", url)
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		ret := &proto.BackupVCodeResponse{}
		if err := httpRsp.Json(ret); err != nil {
			rsp.Code = httpRsp.StatusCode()
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		vcode = ret.VCode
	}

	// Signature.
	{
		body := RSASign(vcode, rsaPrvKey)
		ret := RSASignResponse{}
		if err := unmarshal(body, &ret); err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		signature = ret.Signature
	}

	// Restore.
	{
		path := fmt.Sprintf("%s/api/backup/restore", url)
		req := &proto.BackupRestoreRequest{
			VCode:     vcode,
			Signature: signature,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		ret := &proto.BackupRestoreResponse{}
		if err := httpRsp.Json(ret); err != nil {
			rsp.Code = httpRsp.StatusCode()
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		rsp.Time = ret.Time
		encryptedPrvKey = ret.EncryptedPrvKey
	}

	// Decrypt the masterprvkey.
	{
		body := RSADecrypt(encryptedPrvKey, rsaPrvKey)
		ret := RSADecryptResponse{}
		if err := unmarshal(body, &ret); err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		rsp.MasterPrvKey = ret.PlainText
	}
	return marshal(rsp)
}
