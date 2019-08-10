// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	conf := DefaultConfig()
	b, err := json.MarshalIndent(conf, "", "\t")
	assert.Nil(t, err)
	err = ioutil.WriteFile("/tmp/test.json", b, 0644)
	assert.Nil(t, err)

	got, err := LoadConfig("/tmp/test.json")
	assert.Nil(t, err)
	assert.Equal(t, conf, got)
}

/*
{
	"datadir": "./wallet",
	"chainnet": "testnet",
	"endpoint": ":9099",
	"token_secret": "thresh-wallet-demo-token-secret",
	"spv_provider": "blockstream",
	"enable_vcode": true,
	"vcode_expired": 300,
	"wallet_sync_interval_ms": 30000,
	"smtp": {
		"server": "smtp.gmail.com",
		"port": 456,
		"username": "keyfuse",
		"password": "keyfuse",
		"backup_to": "a@gmail.com,b@gmail.com"
	}
}
*/
func TestLoadSmtpConfig(t *testing.T) {
	conf := DefaultConfig()
	conf.Smtp = &SmtpConfig{
		Server:   "smtp.gmail.com",
		Port:     456,
		UserName: "keyfuse",
		Password: "keyfuse",
		BackupTo: "a@gmail.com,b@gmail.com",
	}
	b, err := json.MarshalIndent(conf, "", "\t")
	assert.Nil(t, err)
	err = ioutil.WriteFile("/tmp/test.json", b, 0644)
	assert.Nil(t, err)

	got, err := LoadConfig("/tmp/test.json")
	assert.Nil(t, err)
	assert.Equal(t, conf, got)
}
