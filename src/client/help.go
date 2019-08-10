// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package client

import (
	"github.com/xandout/gorpl/action"
)

func helpAction(cli *Client) *action.Action {
	return action.New("help", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"commands",
			"usage",
			"example",
		}
		rows = append(rows, []string{"help", "help", "help"})
		rows = append(rows, []string{"dumpkey", "dumpkey", "help"})
		rows = append(rows, []string{"gettoken", "gettoken <vcode>", "gettoken 666888"})
		rows = append(rows, []string{"checkwallet", "checkwallet", "checkwallet"})
		rows = append(rows, []string{"createwallet", "createwallet", "createwallet"})
		rows = append(rows, []string{"backupwallet", "backupwallet", "backupwallet"})
		rows = append(rows, []string{"recoverwallet", "recoverwallet", "recoverwallet"})
		rows = append(rows, []string{"getbalance", "getbalance", "getbalance"})
		rows = append(rows, []string{"gettxs", "gettxs", "gettxs"})
		rows = append(rows, []string{"getaddresses", "getaddresses", "getaddresses"})
		rows = append(rows, []string{"getnewaddress", "getnewaddress", "getnewaddress"})
		rows = append(rows, []string{"getsendfees", "getsendfees <address> <value>", "getsendfees tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw 10000"})
		rows = append(rows, []string{"sendtoaddress", "sendtoaddress <address> <value> <fees>", "sendtoaddress tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw 10000 1000"})
		rows = append(rows, []string{"sendalltoaddress", "sendalltoaddress <address>", "sendalltoaddress tb1qsdp08c4uua6ya865mmxvsqeqlv3gzp2lv5jtsw"})
		PrintQueryOutput(columns, rows)
		return nil, nil
	})
}
