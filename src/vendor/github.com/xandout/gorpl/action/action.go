package action

import (
	"log"
	"strings"
)

// Action is a named func
type Action struct {
	Name     string
	Action   func(args ...interface{}) (interface{}, error)
	Children []Action
}

// AddChild appends a child Action to a.Children
func (a *Action) AddChild(child *Action) *Action {
	a.Children = append(a.Children, *child)
	return child
}

// New creates a root level Action for Repl
func New(cmd string, action func(args ...interface{}) (interface{}, error)) *Action {
	if strings.Count(cmd, " ") > 0 {
		log.Println("cmd cannot contain spaces")
		return &Action{}
	}

	return &Action{Name: cmd, Action: action}

}
