// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package client

import (
	"github.com/xandout/gorpl/action"
)

func dumpKeyAction(cli *Client) *action.Action {
	return action.New("dumpkey", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"masterprvkey(local)",
		}
		rows = append(rows, []string{cli.masterPrvKey})
		PrintQueryOutput(columns, rows)
		return nil, nil
	})
}
