package main

import (
	"fmt"
	"os"

	"github.com/xandout/gorpl"

	"github.com/xandout/gorpl/action"
)

var mode = "table"

func main() {

	exitAction := action.New("exit", func(args ...interface{}) (interface{}, error) {
		fmt.Println("Bye!")
		os.Exit(0)
		return nil, nil
	})
	modeAction := action.New("mode", func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			fmt.Printf("Current mode is %s\n", mode)
		}
		return "", nil
	})
	csvAction := action.New("csv", func(args ...interface{}) (interface{}, error) {
		mode = "csv"
		fmt.Printf("Mode set to %s\n", mode)

		return "", nil
	})
	tableAction := action.New("table", func(args ...interface{}) (interface{}, error) {
		mode = "table"
		fmt.Printf("Mode set to %s\n", mode)
		return "", nil
	})
	csvChild := action.New("csvChild", func(args ...interface{}) (interface{}, error) {
		fmt.Println("csvChild!")
		fmt.Println(args)
		return nil, nil
	})
	csvChildChild := action.New("csvChildChild", func(args ...interface{}) (interface{}, error) {
		fmt.Println("csvChildChild!")
		fmt.Println(args)
		return nil, nil
	})
	csvChild.AddChild(csvChildChild)
	csvAction.AddChild(csvChild)

	modeAction.AddChild(csvAction)
	modeAction.AddChild(tableAction)

	f := gorpl.New(";")
	f.AddAction(*modeAction)
	f.AddAction(*exitAction)
	f.Start()
}
