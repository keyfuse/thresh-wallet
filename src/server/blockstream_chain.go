// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/hex"
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
	Version  int64  `json:"version"`
	Locktime int64  `json:"locktime"`
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
		BlockHeight int64  `json:"block_height"`
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
	return unspents, nil
}

// GetTxs -- used to get transactions by address.
func (c *BlockstreamChain) GetTxs(address string) ([]Tx, error) {
	path := fmt.Sprintf("%s/address/%s/txs", c.url, address)

	httpRsp, err := proto.NewRequest().Get(path)
	if err != nil {
		return nil, err
	}
	if httpRsp.StatusCode() != 200 {
		return nil, fmt.Errorf("blockstream.get.txs.rsp.error:%v", httpRsp.StatusCode())
	}

	var bstxs []BlockstreamTx
	if err := httpRsp.Json(&bstxs); err != nil {
		return nil, err
	}

	var txs []Tx
	for _, tx := range bstxs {
		var data string
		var sentValue int64
		var receivedValue int64

		for _, vin := range tx.Vin {
			if vin.Prevout.ScriptpubkeyAddress == address {
				sentValue += vin.Prevout.Value
			}
		}
		for _, vout := range tx.Vout {
			if vout.ScriptpubkeyAddress == address {
				receivedValue += vout.Value
			}
			if vout.ScriptpubkeyType == "op_return" {
				hexstr := vout.Scriptpubkey
				if len(hexstr) > 4 {
					hexstr = hexstr[4:]
				}
				bytes, _ := hex.DecodeString(hexstr)
				data = string(bytes)
			}
		}

		tx := Tx{
			Txid:        tx.Txid,
			Fee:         tx.Fee,
			Data:        data,
			Value:       receivedValue - sentValue,
			Confirmed:   tx.Status.Confirmed,
			BlockTime:   tx.Status.BlockTime,
			BlockHeight: tx.Status.BlockHeight,
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// GetFees -- used to get the fees from the chain mempool.
func (c *BlockstreamChain) GetFees() (map[string]float32, error) {
	path := fmt.Sprintf("%s/fee-estimates", c.url)

	httpRsp, err := proto.NewRequest().Get(path)
	if err != nil {
		return nil, err
	}
	if httpRsp.StatusCode() != 200 {
		return nil, fmt.Errorf("blockstream.get.fees.rsp.error:%v", httpRsp.StatusCode())
	}

	fees := make(map[string]float32)
	if err := httpRsp.Json(&fees); err != nil {
		return nil, err
	}
	return fees, nil
}

func (c *BlockstreamChain) GetTickers() (map[string]Ticker, error) {
	httpRsp, err := proto.NewRequest().Get("https://blockchain.info/ticker")
	if err != nil {
		return nil, err
	}
	if httpRsp.StatusCode() != 200 {
		return nil, fmt.Errorf("blockstream.get.tickers.rsp.error:%v", httpRsp.StatusCode())
	}

	tickers := make(map[string]Ticker)
	if err := httpRsp.Json(&tickers); err != nil {
		return nil, err
	}
	return tickers, nil
}

// GetTxLink -- get the tx web link.
func (c *BlockstreamChain) GetTxLink() string {
	var url string
	conf := c.conf

	switch conf.ChainNet {
	case testnet:
		url = "https://blockstream.info/testnet/tx/%v"
	case mainnet:
		url = "https://blockstream.info/tx/%v"
	}
	return url
}

// PushTx -- used to push tx to the chain.
func (c *BlockstreamChain) PushTx(hex string) (string, error) {
	log := c.log

	path := fmt.Sprintf("%s/tx", c.url)
	log.Info("chain.blockstream.strart.pushtx:[%v].tx:%v", path, hex)

	httpRsp, err := proto.NewRequest().Post(path, hex)
	if err != nil {
		return "", err
	}
	if httpRsp.StatusCode() != 200 {
		return "", fmt.Errorf("blockstream.push.tx.rsp.error:%v", httpRsp.StatusCode())
	}
	txid := httpRsp.Body()

	log.Info("chain.blockstream.end.pushtx.txid:%v", txid)
	return txid, nil
}
