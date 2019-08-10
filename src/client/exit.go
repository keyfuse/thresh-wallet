// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package client

import (
	"fmt"
	"os"

	"github.com/xandout/gorpl/action"
)

func exitAction(cli *Client) *action.Action {
	return action.New("exit", func(args ...interface{}) (interface{}, error) {
		fmt.Println("Bye!")
		os.Exit(0)
		return nil, nil
	})
}
