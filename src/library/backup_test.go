// thresh-wallet
//
// Copyright 2019 by KeyFuse
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
		assert.False(t, rsp.UserExists)
		assert.False(t, rsp.BackupExists)
	}

	// Create.
	{
		body := APIWalletCreate(ts.URL, token, mockMasterPrvKey, mockMasterPubKey)
		rsp := &WalletCreateResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}

	// Check user exists.
	{
		body := APIWalletCheck(ts.URL, token)
		rsp := &WalletCheckResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.True(t, rsp.UserExists)
		assert.False(t, rsp.BackupExists)
	}

	// Backup.
	{
		body := APIWalletBackup(ts.URL, token, "xx", "icloud", mockRSAPrvKey, mockMasterPrvKey)
		rsp := &WalletBackupResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}

	// Check exists.
	{
		body := APIWalletCheck(ts.URL, token)
		rsp := &WalletCheckResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.True(t, rsp.UserExists)
		assert.True(t, rsp.BackupExists)
	}

	// Restore.
	{
		body := APIWalletRestore(ts.URL, token, mockRSAPrvKey)
		rsp := &WalletRestoreResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
		assert.Equal(t, mockMasterPrvKey, rsp.MasterPrvKey)
	}
}
