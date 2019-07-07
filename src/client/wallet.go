// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package client

import (
	"fmt"

	"library"

	"github.com/xandout/gorpl/action"
)

func walletBalanceAction(cli *Client) *action.Action {
	return action.New("balance", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"all_balance(satoshi)",
			"unconfirmed_balance(satoshi)",
		}

		// Check.
		if cli.token == "" {
			pprintError("token.is.null", "gettoken [vcode]")
			return nil, nil
		}

		// Balance.
		{
			rsp := &library.WalletBalanceResponse{}
			body := library.APIWalletBalance(cli.apiurl, cli.token)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}

			rows = append(rows, []string{fmt.Sprintf("%v", rsp.AllBalance), fmt.Sprintf("%v", rsp.UnconfirmedBalance)})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}

func walletNewAddressAction(cli *Client) *action.Action {
	return action.New("newaddress", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"address",
			"postion",
		}

		// Check.
		if cli.token == "" {
			pprintError("token.is.null", "gettoken [vcode]")
			return nil, nil
		}

		// New address.
		{
			rsp := &library.EcdsaAddressResponse{}
			body := library.APIEcdsaNewAddress(cli.apiurl, cli.token)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}

			rows = append(rows, []string{rsp.Address, fmt.Sprintf("%v", rsp.Pos)})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}
