package primary

import (
	"fmt"
	"strings"

	"github.com/puppetlabs/wash/cmd/internal/find/types"
)

// Action is the action primary
//
// actionPrimary => <action>
//nolint
var Action = Parser.add(&Primary{
	Description: "Returns true if the entry supports action",
	name:        "action",
	args:        "action",
	parseFunc: func(tokens []string) (types.EntryPredicate, []string, error) {
		if len(tokens) == 0 {
			return nil, nil, fmt.Errorf("requires additional arguments")
		}
		action := tokens[0]
		if !types.IsAction(action) {
			// User's querying an invalid action, so return an error.
			validActions := strings.Join(types.Actions(), ", ")
			return nil, nil, fmt.Errorf("%v is an invalid action. Valid actions are %v", tokens[0], validActions)
		}
		p := func(e types.Entry) bool {
			return e.Supports(action)
		}
		return p, tokens[1:], nil
	},
})
