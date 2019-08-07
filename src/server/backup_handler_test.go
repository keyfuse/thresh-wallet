// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"testing"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	"proto"

	"github.com/stretchr/testify/assert"
)

func TestBackupVCodeHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		t.Logf("rsp:%+v", rsp)
	}
}

func TestBackupStoreHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	var err error
	var vcode string
	var pubkeypem string
	var signature string
	var brokensignature string
	var prvkey *rsa.PrivateKey

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	// RSA.
	{
		prvkey, err = rsa.GenerateKey(rand.Reader, 1024)
		assert.Nil(t, err)
		pubkeypem = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&prvkey.PublicKey)}))

		digest, err := hex.DecodeString(vcode)
		assert.Nil(t, err)

		sig, err := rsa.SignPKCS1v15(rand.Reader, prvkey, crypto.SHA256, digest)
		assert.Nil(t, err)

		signature = fmt.Sprintf("%x", sig)
		brokensignature = signature[2:]
	}

	// vcode 400.
	{
		req := &proto.BackupStoreRequest{
			VCode:            "xx",
			Signature:        signature,
			EncryptionPubKey: pubkeypem,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/store", req)
		assert.Nil(t, err)
		assert.Equal(t, 400, httpRsp.StatusCode())
	}

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	//  signature 400.
	{
		req := &proto.BackupStoreRequest{
			VCode:            vcode,
			Signature:        brokensignature,
			EncryptionPubKey: pubkeypem,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/store", req)
		assert.Nil(t, err)
		assert.Equal(t, 400, httpRsp.StatusCode())
	}

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	// RSA.
	{
		digest, err := hex.DecodeString(vcode)
		assert.Nil(t, err)

		sig, err := rsa.SignPKCS1v15(rand.Reader, prvkey, crypto.SHA256, digest)
		assert.Nil(t, err)

		signature = fmt.Sprintf("%x", sig)
	}

	//  200 ok.
	{
		req := &proto.BackupStoreRequest{
			VCode:            vcode,
			Signature:        signature,
			EncryptedPrvKey:  "fake",
			EncryptionPubKey: pubkeypem,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/store", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	//  exists 400.
	{
		req := &proto.BackupStoreRequest{
			VCode:            vcode,
			Signature:        signature,
			EncryptedPrvKey:  "fake",
			EncryptionPubKey: pubkeypem,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/store", req)
		assert.Nil(t, err)
		assert.Equal(t, 400, httpRsp.StatusCode())
	}
}

func TestBackupRestoreHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	var err error
	var vcode string
	var signature string
	var pubkeypem string
	var brokensignature string
	var prvkey *rsa.PrivateKey

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	// RSA.
	{
		prvkey, err = rsa.GenerateKey(rand.Reader, 1024)
		assert.Nil(t, err)
		pubkeypem = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&prvkey.PublicKey)}))

		digest, err := hex.DecodeString(vcode)
		assert.Nil(t, err)

		sig, err := rsa.SignPKCS1v15(rand.Reader, prvkey, crypto.SHA256, digest)
		assert.Nil(t, err)

		signature = fmt.Sprintf("%x", sig)
		brokensignature = signature[2:]
	}

	//  store 200 ok.
	{
		req := &proto.BackupStoreRequest{
			VCode:            vcode,
			Signature:        signature,
			EncryptedPrvKey:  "fake",
			EncryptionPubKey: pubkeypem,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/store", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}

	// vcode 400.
	{
		req := &proto.BackupRestoreRequest{
			VCode:     "xx",
			Signature: signature,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/restore", req)
		assert.Nil(t, err)
		assert.Equal(t, 400, httpRsp.StatusCode())
	}

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	//  signature 400.
	{
		req := &proto.BackupRestoreRequest{
			VCode:     vcode,
			Signature: brokensignature,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/restore", req)
		assert.Nil(t, err)
		assert.Equal(t, 400, httpRsp.StatusCode())
	}

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	// RSA.
	{
		digest, err := hex.DecodeString(vcode)
		assert.Nil(t, err)

		sig, err := rsa.SignPKCS1v15(rand.Reader, prvkey, crypto.SHA256, digest)
		assert.Nil(t, err)

		signature = fmt.Sprintf("%x", sig)
	}

	//  200 ok.
	{
		req := &proto.BackupRestoreRequest{
			VCode:     vcode,
			Signature: signature,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/restore", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestBackupVerifyHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	var err error
	var vcode string
	var signature string
	var pubkeypem string
	var pubkeyhash string
	var prvkey *rsa.PrivateKey

	// VCode.
	{
		req := &proto.BackupVCodeRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVCodeResponse{}
		httpRsp.Json(rsp)
		vcode = rsp.VCode
	}

	// RSA.
	{
		prvkey, err = rsa.GenerateKey(rand.Reader, 1024)
		assert.Nil(t, err)
		pubkeypem = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&prvkey.PublicKey)}))
		pubkeyhash = fmt.Sprintf("%x", sha256.Sum256([]byte(pubkeypem)))

		digest, err := hex.DecodeString(vcode)
		assert.Nil(t, err)

		sig, err := rsa.SignPKCS1v15(rand.Reader, prvkey, crypto.SHA256, digest)
		assert.Nil(t, err)

		signature = fmt.Sprintf("%x", sig)
	}

	//  store 200 ok.
	{
		req := &proto.BackupStoreRequest{
			VCode:            vcode,
			Signature:        signature,
			EncryptedPrvKey:  "fake",
			EncryptionPubKey: pubkeypem,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/store", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}

	//  200 ok.
	{
		req := &proto.BackupVerifyRequest{
			EncryptionPubKeyHash: pubkeyhash,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/verify", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVerifyResponse{}
		httpRsp.Json(rsp)
		assert.True(t, rsp.VerifyPassed)
	}

	//  200 but not passed.
	{
		req := &proto.BackupVerifyRequest{
			EncryptionPubKeyHash: pubkeyhash[2:],
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/backup/verify", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.BackupVerifyResponse{}
		httpRsp.Json(rsp)
		assert.False(t, rsp.VerifyPassed)
	}
}
