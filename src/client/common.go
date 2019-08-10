// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/olekukonko/tablewriter"
)

const (
	TestNet = "testnet"
	MainNet = "mainnet"
)

func marshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func unmarshal(data string, t interface{}) error {
	if err := json.Unmarshal([]byte(data), t); err != nil {
		return err
	}
	return nil
}

func expandTabsAndNewLines(s string) string {
	var buf bytes.Buffer
	// 4-wide columns, 1 character minimum width.
	w := tabwriter.NewWriter(&buf, 4, 0, 1, ' ', 0)
	fmt.Fprint(w, strings.Replace(s, "\n", "␤\n", -1))
	_ = w.Flush()
	return buf.String()
}

func PrintQueryOutput(cols []string, allRows [][]string) {
	if len(cols) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetHeader(cols)
	for _, row := range allRows {
		for i, r := range row {
			row[i] = expandTabsAndNewLines(r)
		}
		table.Append(row)
	}
	table.Render()
	nRows := len(allRows)
	fmt.Fprintf(os.Stdout, "(%d rows)\n", nRows)
}

func pprintError(msg string, help string) {
	var rows [][]string
	columns := []string{
		"error",
		"help",
	}
	rows = append(rows, []string{msg, help})
	PrintQueryOutput(columns, rows)
}

func mask(data string) string {
	pad := 16
	size := len(data)
	if size < pad {
		return data
	}

	prefix := data[0:8]
	suffix := data[size-8:]
	m := strings.Repeat("*", size-pad)
	return prefix + m + suffix
}
