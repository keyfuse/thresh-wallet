// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package client

import (
	"library"

	"github.com/xandout/gorpl/action"
)

func tokenAction(cli *Client) *action.Action {
	return action.New("gettoken", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"status",
		}

		if len(args) != 1 {
			pprintError("args.invalid", "gettoken [vcode], example:gettoken 666666")
			return nil, nil
		}

		// Get token.
		{
			rsp := &library.TokenResponse{}
			body := library.APIGetToken(cli.apiurl, cli.uid, args[0].(string), cli.masterPubKey)
			if err := unmarshal(body, rsp); err != nil {
				rows = append(rows, []string{err.Error()})
				PrintQueryOutput(columns, rows)
				return nil, nil
			}

			if rsp.Code != 200 {
				rows = append(rows, []string{rsp.Message})
				PrintQueryOutput(columns, rows)
				return nil, nil
			}
			cli.token = rsp.Token
			rows = append(rows, []string{"OK"})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}
