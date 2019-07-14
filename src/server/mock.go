package server

import (
	"io/ioutil"
	"net/http/httptest"
	"os"

	"xlog"
)

const (
	mockUID             = "10086"
	mockToken           = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaWQiOiIiLCJtcGsiOiJ0cHViRDZOelZia3JZaFo0WHhoZWF1dXNicVpCQlJoVUFwU01OekJiTU1WSkJlR0plUlBwQVFRRWh4RWZDZUxmbVV5ZXQzRlhYeWJBb1doSjN1WmU0ZlF2cWdWQ2Q4VVBLWDhzUDRxQVhLRUhaR2siLCJuZXQiOiJ0ZXN0bmV0IiwidCI6MTU2MzAzMzc0MCwidWlkIjoiMTAwODYifQ.3gqSJBP4uXJ-P3DNlsZjSBrCCCM_Kp4vAZBJ7I7lc78"
	mockSvrMasterPrvKey = "tprv8ZgxMBicQKsPfNhXDHV93ummM6rEzTmxHf96Mk3FnpgoaoNYPjfSCZyHFnFQnQDLAiMNsvJqEtvjCkvo5P3CPRHQx5GcZxPqRHy31q2oWXD"
	mockCliMasterPrvKey = "tprv8ZgxMBicQKsPeVfrhGFHCRu4cQBY1VFSogap4qSzmNTuow93Y1aeXTco2Vdw41VLUvPC4e3X1ZF9uoJEeRbUpLR4DqtzvLd3AQnQobNaGA4"
	mockCliMasterPubKey = "tpubD6NzVbkrYhZ4XxheauusbqZBBRhUApSMNzBbMMVJBeGJeRPpAQQEhxEfCeLfmUyet3FXXybAoWhJ3uZe4fQvqgVCd8UPKX8sP4qAXKEHZGk"
)

func MockConfig() *Config {
	conf := DefaultConfig()
	conf.DisableVCode = true
	conf.DataDir = "/tmp/tss"

	return conf
}

func MockServer() *httptest.Server {
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	conf := MockConfig()

	os.MkdirAll(conf.DataDir, os.ModePerm)
	os.RemoveAll(conf.DataDir + "/*")
	err := ioutil.WriteFile(conf.DataDir+"/10086.json", []byte(mock10086Json), 0644)
	if err != nil {
		panic(err)
	}
	router := NewAPIRouter(log, conf)

	// Set mock chain.
	router.handler.wdb.setChain(newMockChain(log))
	return httptest.NewServer(router)
}

type mockChain struct {
	log *xlog.Log
}

func newMockChain(log *xlog.Log) *mockChain {
	return &mockChain{log: log}
}

func (c *mockChain) GetUTXO(address string) ([]Unspent, error) {
	var unspents []Unspent

	if address == "mnBETqvxTqcFRSLnR3w2Tpe9Qu58EasQgU" {
		unspent1 := Unspent{
			Txid:         "0f8c5cdf448acb82969193452ac4bb7010c0890ceb96fa5e8c332378654459df",
			Vout:         0,
			Value:        93266,
			Confirmed:    true,
			BlockTime:    1562492930,
			BlockHeight:  1567884,
			Scriptpubkey: "76a914490e0eebcc5d462221ea38d00a6aee1238db2a5788ac",
		}
		unspent2 := Unspent{
			Txid:         "2335b1b00d149907e0ce9eb349da87234d2c9bd0dfcc216cb251c3b21d63054a",
			Vout:         1,
			Value:        10000,
			Confirmed:    true,
			BlockTime:    1562492930,
			BlockHeight:  1567884,
			Scriptpubkey: "76a914490e0eebcc5d462221ea38d00a6aee1238db2a5788ac",
		}
		unspents = append(unspents, unspent1, unspent2)
	}
	return unspents, nil
}

func (c *mockChain) PushTx(hex string) (string, error) {
	return "e0c328bd49e9a1c2ef5f7a1c14f0f9893658f5673fb415ceec1125dcd6641993", nil
}

var (
	mock10086Json = `{
 "uid": "10086",
 "did": "",
 "backup": null,
 "lastpos": 7,
 "address": {
  "miqi14i2nweWYkcAh49E8Zk6gVAta7ohqJ": {
   "pos": 6,
   "address": "miqi14i2nweWYkcAh49E8Zk6gVAta7ohqJ",
   "unspents": null,
   "balance": {
    "all_balance": 0,
    "unconfirmed_balance": 0
   }
  },
  "mmBRSnFG7o1BX5DaK8Da3xKxvjBh6fzNQq": {
   "pos": 3,
   "address": "mmBRSnFG7o1BX5DaK8Da3xKxvjBh6fzNQq",
   "unspents": null,
   "balance": {
    "all_balance": 0,
    "unconfirmed_balance": 0
   }
  },
  "mnBETqvxTqcFRSLnR3w2Tpe9Qu58EasQgU": {
   "pos": 2,
   "address": "mnBETqvxTqcFRSLnR3w2Tpe9Qu58EasQgU",
   "unspents": [
    {
     "txid": "0f8c5cdf448acb82969193452ac4bb7010c0890ceb96fa5e8c332378654459df",
     "vout": 0,
     "value": 93266,
     "confirmed": true,
     "block_time": 1562492930,
     "block_height": 1567884,
     "Scriptpubkey": "76a914490e0eebcc5d462221ea38d00a6aee1238db2a5788ac"
    },
    {
     "txid": "2335b1b00d149907e0ce9eb349da87234d2c9bd0dfcc216cb251c3b21d63054a",
     "vout": 1,
     "value": 10000,
     "confirmed": true,
     "block_time": 1562492930,
     "block_height": 1567884,
     "Scriptpubkey": "76a914490e0eebcc5d462221ea38d00a6aee1238db2a5788ac"
    }
   ],
   "balance": {
    "all_balance": 103266,
    "unconfirmed_balance": 0
   }
  },
  "msV128vgApMNEFbTUy5wto12ucZNFdtKTA": {
   "pos": 0,
   "address": "msV128vgApMNEFbTUy5wto12ucZNFdtKTA",
   "unspents": null,
   "balance": {
    "all_balance": 0,
    "unconfirmed_balance": 0
   }
  },
  "msYdTCo8sxSWNdgdNUsxMM1ghA44mNaksY": {
   "pos": 4,
   "address": "msYdTCo8sxSWNdgdNUsxMM1ghA44mNaksY",
   "unspents": null,
   "balance": {
    "all_balance": 0,
    "unconfirmed_balance": 0
   }
  },
  "muAK3ufJer1nSUerdf95r5As442DagfBXS": {
   "pos": 5,
   "address": "muAK3ufJer1nSUerdf95r5As442DagfBXS",
   "unspents": null,
   "balance": {
    "all_balance": 0,
    "unconfirmed_balance": 0
   }
  },
  "mv7hzrEL4WYXvMzLawe82Mn82Mm7had4FY": {
   "pos": 1,
   "address": "mv7hzrEL4WYXvMzLawe82Mn82Mm7had4FY",
   "unspents": null,
   "balance": {
    "all_balance": 0,
    "unconfirmed_balance": 0
   }
  }
 },
 "svrmasterprvkey": "tprv8ZgxMBicQKsPfNhXDHV93ummM6rEzTmxHf96Mk3FnpgoaoNYPjfSCZyHFnFQnQDLAiMNsvJqEtvjCkvo5P3CPRHQx5GcZxPqRHy31q2oWXD",
 "climasterpubkey": "tpubD6NzVbkrYhZ4XxheauusbqZBBRhUApSMNzBbMMVJBeGJeRPpAQQEhxEfCeLfmUyet3FXXybAoWhJ3uZe4fQvqgVCd8UPKX8sP4qAXKEHZGk"
}`
)
