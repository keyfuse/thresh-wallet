// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package client

import (
	"fmt"
	"strconv"
	"time"

	"library"

	"github.com/xandout/gorpl/action"
)

func walletBalanceAction(cli *Client) *action.Action {
	return action.New("getbalance", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"current_balance",
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

			rows = append(rows, []string{fmt.Sprintf("%v", rsp.CoinValue)})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}

func walletTxsAction(cli *Client) *action.Action {
	return action.New("gettxs", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"txid",
			"direction",
			"value",
			"confirmed",
			"time",
			"height",
		}

		// Check.
		if cli.token == "" {
			pprintError("token.is.null", "gettoken [vcode]")
			return nil, nil
		}

		{
			rsp := &library.WalletTxsResponse{}
			body := library.APIWalletTxs(cli.apiurl, cli.token, 0, 20)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}

			for _, tx := range rsp.Txs {
				txid := tx.Txid
				direction := "received"
				if tx.Value < 0 {
					direction = "sent"
				}
				value := tx.Value
				confirmed := tx.Confirmed
				ts := time.Unix(tx.BlockTime, 0)
				height := tx.BlockHeight
				rows = append(rows, []string{
					fmt.Sprintf("%v", txid),
					fmt.Sprintf("%v", direction),
					fmt.Sprintf("%v", value),
					fmt.Sprintf("%v", confirmed),
					fmt.Sprintf("%s", ts),
					fmt.Sprintf("%v", height),
				})
			}
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}

func walletNewAddressAction(cli *Client) *action.Action {
	return action.New("getnewaddress", func(args ...interface{}) (interface{}, error) {
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

func walletSendFeesAction(cli *Client) *action.Action {
	return action.New("getsendfees", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"toaddress",
			"send_value",
			"sendable_value",
			"fees(sat)",
			"speed(Fast/Normal/Slow)",
		}

		// Check.
		if cli.token == "" {
			pprintError("token.is.null", "gettoken [vcode]")
			return nil, nil
		}

		if len(args) != 2 {
			pprintError("args.invalid", "getfees <address> <amount>")
			return nil, nil
		}

		address := args[0].(string)
		value, err := strconv.ParseUint(args[1].(string), 10, 64)
		if err != nil {
			pprintError("amount.invalid", "getfees <address> <amount>")
			return nil, nil
		}

		{
			rsp := &library.WalletSendFeesResponse{}
			body := library.APIWalletSendFees(cli.apiurl, cli.token, value)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}

			rows = append(rows, []string{
				address,
				fmt.Sprintf("%v", value),
				fmt.Sprintf("%v", rsp.SendableValue),
				fmt.Sprintf("%v", rsp.Fees),
				rsp.FeeMode})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}

func walletSendToAddressAction(cli *Client) *action.Action {
	return action.New("sendtoaddress", func(args ...interface{}) (interface{}, error) {
		var rows [][]string
		columns := []string{
			"toaddress",
			"value(sat)",
			"fees(sat)",
			"txid",
		}

		// Check.
		if cli.token == "" {
			pprintError("token.is.null", "gettoken [vcode]")
			return nil, nil
		}

		if len(args) != 3 {
			pprintError("args.invalid", "sendtoaddress <address> <amount> <fees>")
			return nil, nil
		}

		address := args[0].(string)
		value, err := strconv.ParseUint(args[1].(string), 10, 64)
		if err != nil {
			pprintError("amount.invalid", "sendtoaddress <address> <amount> <fees>")
			return nil, nil
		}

		fees, err := strconv.ParseUint(args[2].(string), 10, 64)
		if err != nil {
			pprintError("fees.invalid", "sendtoaddress <address> <amount> <fees>")
			return nil, nil
		}

		{
			rsp := &library.WalletSendResponse{}
			body := library.APIWalletSend(cli.apiurl, cli.token, cli.net, cli.masterPrvKey, address, value, fees)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}

			rows = append(rows, []string{address, fmt.Sprintf("%v", value), fmt.Sprintf("%v", fees), rsp.TxID})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}

func walletSendAllToAddressAction(cli *Client) *action.Action {
	return action.New("sendalltoaddress", func(args ...interface{}) (interface{}, error) {
		var fees uint64
		var balance uint64
		var sendable uint64

		var rows [][]string
		columns := []string{
			"toaddress",
			"value(sat)",
			"fees(sat)",
			"txid",
		}

		// Check.
		if cli.token == "" {
			pprintError("token.is.null", "gettoken [vcode]")
			return nil, nil
		}

		if len(args) != 1 {
			pprintError("args.invalid", "sendalltoaddress <address>")
			return nil, nil
		}
		address := args[0].(string)

		// Get all balance.
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
			balance = rsp.CoinValue
		}

		{
			rsp := &library.WalletSendFeesResponse{}
			body := library.APIWalletSendFees(cli.apiurl, cli.token, balance)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}
			fees = rsp.Fees
			sendable = rsp.SendableValue
		}

		{
			rsp := &library.WalletSendResponse{}
			body := library.APIWalletSend(cli.apiurl, cli.token, cli.net, cli.masterPrvKey, address, sendable, fees)
			if err := unmarshal(body, rsp); err != nil {
				pprintError(err.Error(), "")
				return nil, nil
			}

			if rsp.Code != 200 {
				pprintError(rsp.Message, "")
				return nil, nil
			}

			rows = append(rows, []string{address, fmt.Sprintf("%v", sendable), fmt.Sprintf("%v", fees), rsp.TxID})
			PrintQueryOutput(columns, rows)
		}
		return nil, nil
	})
}
