// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"

	"proto"
	"xlog"
)

// BlockstreamUTXO --
type BlockstreamUTXO struct {
	Txid   string `json:"txid"`
	Vout   uint32 `json:"vout"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight uint32 `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   uint32 `json:"block_time"`
	} `json:"status"`
	Value uint64 `json:"value"`
}

// BlockstreamTx --
type BlockstreamTx struct {
	Txid     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int    `json:"locktime"`
	Vin      []struct {
		Txid    string `json:"txid"`
		Vout    int    `json:"vout"`
		Prevout struct {
			Scriptpubkey        string `json:"scriptpubkey"`
			ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
			ScriptpubkeyAddress string `json:"scriptpubkey_address"`
			ScriptpubkeyType    string `json:"scriptpubkey_type"`
			Value               int64  `json:"value"`
		} `json:"prevout"`
		Scriptsig    string   `json:"scriptsig"`
		ScriptsigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     int64    `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey        string `json:"scriptpubkey"`
		ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
		ScriptpubkeyAddress string `json:"scriptpubkey_address"`
		ScriptpubkeyType    string `json:"scriptpubkey_type"`
		Value               int64  `json:"value"`
	} `json:"vout"`
	Size   int   `json:"size"`
	Weight int   `json:"weight"`
	Fee    int64 `json:"fee"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int64  `json:"block_time"`
	} `json:"status"`
}

// BlockstreamChain --
type BlockstreamChain struct {
	url  string
	log  *xlog.Log
	conf *Config
}

// NewBlockstreamChain -- creates new BlockstreamChain.
func NewBlockstreamChain(log *xlog.Log, conf *Config) Chain {
	var url string
	switch conf.ChainNet {
	case testnet:
		url = "https://blockstream.info/testnet/api"
	case mainnet:
		url = "https://blockstream.info/api"
	}

	return &BlockstreamChain{
		log:  log,
		conf: conf,
		url:  url,
	}
}

// GetUTXO -- used to get all the unspents of this address.
func (c *BlockstreamChain) GetUTXO(address string) ([]Unspent, error) {
	log := c.log

	path := fmt.Sprintf("%s/address/%s/utxo", c.url, address)
	log.Info("chain.blockstream.strart.getutxo:[%v]", path)

	httpRsp, err := proto.NewRequest().Get(path)
	if err != nil {
		return nil, err
	}
	if httpRsp.StatusCode() != 200 {
		return nil, fmt.Errorf("blockstream.get.utxo.rsp.error:%v", httpRsp.StatusCode())
	}

	var bsutxo []BlockstreamUTXO
	if err := httpRsp.Json(&bsutxo); err != nil {
		return nil, err
	}

	var unspents []Unspent
	for _, utxo := range bsutxo {
		path = fmt.Sprintf("%s/tx/%s", c.url, utxo.Txid)
		httpRsp, err = proto.NewRequest().Get(path)
		if err != nil {
			log.Error("blockstream.utxo[%v].get.tx.error:%+v", utxo.Txid, err)
			continue
		}
		if httpRsp.StatusCode() != 200 {
			log.Error("blockstream.utxo[%v].get.tx.error:%v", utxo.Txid, httpRsp.StatusCode())
			continue
		}
		bsTx := &BlockstreamTx{}
		if err := httpRsp.Json(bsTx); err != nil {
			log.Error("blockstream.utxo[%v].unmarsh.tx.error:%+v", utxo.Txid, err)
			continue
		}
		unspent := Unspent{
			Txid:         utxo.Txid,
			Vout:         utxo.Vout,
			Value:        utxo.Value,
			Confirmed:    utxo.Status.Confirmed,
			BlockTime:    utxo.Status.BlockTime,
			BlockHeight:  utxo.Status.BlockHeight,
			Scriptpubkey: bsTx.Vout[utxo.Vout].Scriptpubkey,
		}
		unspents = append(unspents, unspent)
	}

	log.Info("chain.blockstream.end.getutxo[%+v].time:%v", unspents, httpRsp.Cost())
	return unspents, nil
}
