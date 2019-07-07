package gorpl

import (
	"fmt"
	"log"
	"strings"

	"github.com/xandout/gorpl/action"

	"github.com/chzyer/readline"
)

// Recursive function used to build PrefixCompleter
// I barely understand this...it breaks any time I touch it
func wrFunc(actions []action.Action) []readline.PrefixCompleterInterface {
	var pci []readline.PrefixCompleterInterface
	for _, a := range actions {

		if len(a.Children) > 0 {
			pci = append(pci, readline.PcItem(a.Name, wrFunc(a.Children)...))
		} else {
			pci = append(pci, readline.PcItem(a.Name))
		}
	}
	return pci
}

// Used to build PrefixCompleter...calls wrFunc
func (r *Repl) walkActions() {

	var pcis []readline.PrefixCompleterInterface
	for _, aa := range r.Actions {
		if len(aa.Children) > 0 {
			pcis = append(pcis, readline.PcItem(aa.Name, wrFunc(aa.Children)...))
		} else {
			pcis = append(pcis, readline.PcItem(aa.Name))
		}
	}

	r.RL.Config.AutoComplete = readline.NewPrefixCompleter(pcis...)
}

// Repl houses all of our config data
type Repl struct {
	RL         *readline.Instance
	Actions    []action.Action
	Default    action.Action
	Prefix     string
	Terminator string
}

// New sets up the Repl
func New(term string) Repl {
	r, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            "/tmp/repl.history",
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return Repl{
		RL:         r,
		Terminator: term,
	}
}

// AddAction registers a named function, Action
func (r *Repl) AddAction(action action.Action) {
	r.Actions = append(r.Actions, action)
}

// Start runs the Repl as configured
func (r *Repl) Start() error {

	defer r.RL.Close()
	// Setup readline completer
	r.walkActions()

	// History
	var cmdHist []string
	var lines []string

	// If user does not supply default Action then we just echo
	if r.Default.Action == nil {
		r.Default = action.Action{
			Action: func(args ...interface{}) (interface{}, error) {
				fmt.Println(args)
				return nil, nil
			},
		}
	}

REPL_LOOP:
	for {
		line, err := r.RL.Readline()
		if err != nil {
			log.Fatal(err)
		}
		// Remove extra whitespace
		line = strings.TrimSpace(line)

		cmd := strings.Split(line, " ")

		// Empty input
		if len(line) == 0 {
			continue
		}

		cmdHist = append(cmdHist, line)

		// Is this a built in command?
		for _, a := range r.Actions {
			if a.Name == cmd[0] {
				// Must be registered.  Now we need to work out children

				args := make([]interface{}, len(cmd[1:]))
				for i, v := range cmd[1:] {
					args[i] = v
				}

				// Determine if child
				var childRunner func(parentAction action.Action, children []action.Action, newArgs []interface{})
				childRunner = func(parentAction action.Action, children []action.Action, newArgs []interface{}) {
					mapped := make(map[string]action.Action, len(children))
					for _, a := range children {
						mapped[a.Name] = a
					}
					if len(newArgs) > 0 {
						firstarg := newArgs[0].(string)
						if act, ok := mapped[firstarg]; ok {
							childRunner(act, act.Children, newArgs[1:])
						} else {
							parentAction.Action(newArgs...)
						}
					} else {
						parentAction.Action()
					}
				}
				childRunner(a, a.Children, args)
				//r.RL.SetPrompt("> ")
				r.RL.SaveHistory(line)
				continue REPL_LOOP
			}
		}

		// NOT a registered Action, treat as string with multiline and pass to default.
		if !strings.HasSuffix(line, r.Terminator) {
			lines = append(lines, line)
			r.RL.SetPrompt(">>> ")
			r.RL.SaveHistory(line)
			continue
		} else {
			args := make([]interface{}, 1)
			lines = append(lines, line)
			//r.RL.SetPrompt("> ")
			r.RL.SaveHistory(line)
			args[0] = strings.Join(lines, " ")
			ret, err := r.Default.Action(args...)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println(ret)
			}
			lines = []string{}
			continue
		}
	}
}
