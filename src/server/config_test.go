// thresh-wallet
//
// Copyright 2019 by KeyFuse
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
