package apitypes

import "github.com/puppetlabs/wash/plugin"

// Entry represents a Wash entry as interpreted by the API.
//
// swagger:response
type Entry struct {
	TypeID     string                 `json:"type_id"`
	Path       string                 `json:"path"`
	Actions    []string               `json:"actions"`
	Name       string                 `json:"name"`
	CName      string                 `json:"cname"`
	Attributes plugin.EntryAttributes `json:"attributes"`
	Metadata   plugin.JSONObject      `json:"metadata"`
}

func NewEntry(e plugin.Entry) Entry {
	return Entry{
		TypeID:     plugin.TypeID(e),
		Name:       plugin.Name(e),
		CName:      plugin.CName(e),
		Actions:    plugin.SupportedActionsOf(e),
		Attributes: plugin.Attributes(e),
		Metadata:   plugin.PartialMetadata(e),
	}
}

// Supports returns true if e supports the given action, false
// otherwise.
func (e *Entry) Supports(action plugin.Action) bool {
	for _, a := range e.Actions {
		if action.Name == a {
			return true
		}
	}
	return false
}
