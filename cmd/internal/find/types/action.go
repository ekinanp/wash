package types

import "github.com/puppetlabs/wash/plugin"
import mapset "github.com/deckarep/golang-set"

// Actions returns a list of all the Wash actions
func Actions() []string {
	return actions
}

// IsAction returns true if the specified action is a Wash
// action, false otherwise
func IsAction(action string) bool {
	_, ok := actionsMap[action]
	return ok
}

var actionsMap = func() map[string]plugin.Action {
	// Avoid the unnecessary copying done by plugin.Actions
	return plugin.Actions()
}()

var actions = func() []string {
	var as []string
	for action := range actionsMap {
		as = append(as, action)
	}
	return as
}()

var actionsSet = func() mapset.Set {
	return toSet(Actions())
}()

func toSet(actions []string) mapset.Set {
	s := mapset.NewThreadUnsafeSet()
	for _, action := range actions {
		s.Add(action)
	}
	return s
}
