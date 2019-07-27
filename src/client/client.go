// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package client

import (
	"fmt"

	"library"

	"github.com/xandout/gorpl"
)

type Client struct {
	net          string
	uid          string
	token        string
	apiurl       string
	rsaPrvKey    string
	rsaPubKey    string
	masterPrvKey string
	masterPubKey string
}

func NewClient(apiurl string, uid string, chainnet string, masterPrvKey string) *Client {
	net := "testnet"
	mkpubkey := ""
	mkprvkey := masterPrvKey

	// Chainnet.
	if chainnet == "mainnet" {
		net = "mainnet"
	}

	// Vcode.
	{
		body := library.APIGetVCode(apiurl, uid)
		rsp := &library.VcodeResponse{}
		if err := unmarshal(body, rsp); err != nil {
			panic(err)
		}
		if rsp.Code != 200 {
			panic(rsp.Message)
		}
	}

	// Key.
	{
		if mkprvkey != "" {
			body := library.GetMasterPubKey(net, mkprvkey)
			rsp := &library.MasterPubKeyResponse{}
			if err := unmarshal(body, rsp); err != nil {
				panic(err)
			}
			if rsp.Code != 200 {
				panic(rsp.Message)
			}
			mkpubkey = rsp.MasterPubKey
		}
	}

	// Help action.
	helpAction(nil).Action()
	return &Client{
		net:          net,
		uid:          uid,
		apiurl:       apiurl,
		masterPrvKey: mkprvkey,
		masterPubKey: mkpubkey,
	}
}

func (cli *Client) Start() {
	f := gorpl.New("")
	f.RL.SetPrompt(fmt.Sprintf("threshwallet@%s> ", cli.net))
	f.AddAction(*exitAction(cli))
	f.AddAction(*helpAction(cli))
	f.AddAction(*dumpKeyAction(cli))
	f.AddAction(*tokenAction(cli))
	f.AddAction(*walletCheckAction(cli))
	f.AddAction(*walletCreateAction(cli))
	f.AddAction(*walletBackupAction(cli))
	f.AddAction(*walletRecoverAction(cli))
	f.AddAction(*walletBalanceAction(cli))
	f.AddAction(*walletTxsAction(cli))
	f.AddAction(*walletNewAddressAction(cli))
	f.AddAction(*walletSendFeesAction(cli))
	f.AddAction(*walletSendToAddressAction(cli))
	f.AddAction(*walletSendAllToAddressAction(cli))
	f.Start()
}
