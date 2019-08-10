// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"testing"

	"server"

	"github.com/stretchr/testify/assert"
)

func TestWalletBackup(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	mobile := "10096"
	// Token.
	{
		body := APIGetToken(ts.URL, mobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	// Check.
	{
		body := APIWalletCheck(ts.URL, token)
		rsp := &WalletCheckResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.False(t, rsp.WalletExists)
		assert.False(t, rsp.BackupExists)
	}

	// Create.
	{
		body := APIWalletCreate(ts.URL, token, mockMasterPrvKey)
		rsp := &WalletCreateResponse{}
		unmarshal(body, rsp)

		t.Logf("create.rsp:%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}

	// Check user exists.
	{
		body := APIWalletCheck(ts.URL, token)
		rsp := &WalletCheckResponse{}
		unmarshal(body, rsp)

		t.Logf("check.rsp:%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.True(t, rsp.WalletExists)
		assert.False(t, rsp.BackupExists)
	}

	// Backup.
	{
		body := APIWalletBackup(ts.URL, token, "xx", "icloud", mockRSAPrvKey, mockMasterPrvKey)
		rsp := &WalletBackupResponse{}
		unmarshal(body, rsp)

		t.Logf("backup.rsp:%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}

	// Check exists.
	{
		body := APIWalletCheck(ts.URL, token)
		rsp := &WalletCheckResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.True(t, rsp.WalletExists)
		assert.True(t, rsp.BackupExists)
	}

	// Restore.
	{
		body := APIWalletRestore(ts.URL, token, mockRSAPrvKey)
		rsp := &WalletRestoreResponse{}
		unmarshal(body, rsp)

		t.Logf("restore.rsp:%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.Equal(t, mockMasterPrvKey, rsp.MasterPrvKey)
	}

	// Verify.
	{
		body := APIWalletBackupVerify(ts.URL, token, mockRSAPrvKey)
		rsp := &WalletBackupVerifyResponse{}
		unmarshal(body, rsp)

		t.Logf("verify.rsp:%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.True(t, rsp.VerifyPassed)
	}
}
