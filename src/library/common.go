// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"encoding/json"
)

const (
	TestNet = "testnet"
	MainNet = "mainnet"
)

// Status --
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func marshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func unmarshal(data string, t interface{}) error {
	if err := json.Unmarshal([]byte(data), t); err != nil {
		return err
	}
	return nil
}
