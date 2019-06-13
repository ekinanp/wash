package types

import (
	mapset "github.com/deckarep/golang-set"
	apitypes "github.com/puppetlabs/wash/api/types"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
	"github.com/puppetlabs/wash/plugin"
)

// Entry represents an Entry as interpreted by `wash find`
type Entry struct {
	apitypes.Entry
	NormalizedPath string
	Metadata       plugin.JSONObject
}

// NewEntry constructs a new `wash find` entry
func NewEntry(e apitypes.Entry, normalizedPath string) Entry {
	return Entry{
		Entry:          e,
		NormalizedPath: normalizedPath,
		Metadata:       e.Attributes.Meta(),
	}
}

// EntryPredicate represents a predicate on a Wash entry.
type EntryPredicate struct {
	P func(Entry) bool
	// The RequiredActions are evaluated separately so that we can
	// optimize the walker's search. Otherwise, there is no way to
	// get this information from "P" alone.
	RequiredActions mapset.Set
}

// And returns p1 && p2
func (p1 EntryPredicate) And(p2 predicate.Predicate) predicate.Predicate {
	ep2 := p2.(EntryPredicate)
	return EntryPredicate{
		P: func(e Entry) bool {
			return p1.P(e) && ep2.P(e)
		},
		RequiredActions: p1.RequiredActions.Intersect(ep2.RequiredActions),
	}
}

// Or returns p1 || p2
func (p1 EntryPredicate) Or(p2 predicate.Predicate) predicate.Predicate {
	ep2 := p2.(EntryPredicate)
	return EntryPredicate{
		P: func(e Entry) bool {
			return p1.P(e) || ep2.P(e)
		},
		RequiredActions: p1.RequiredActions.Union(ep2.RequiredActions),
	}
}

// Negate returns Not(p1)
func (p1 EntryPredicate) Negate() predicate.Predicate {
	return EntryPredicate{
		P: func(e Entry) bool {
			return !p1.P(e)
		},
		RequiredActions: actionsSet.Difference(p1.RequiredActions),
	}
}

// IsSatisfiedBy returns true if v satisfies the predicate, false otherwise
func (p1 EntryPredicate) IsSatisfiedBy(v interface{}) bool {
	entry, ok := v.(Entry)
	if !ok {
		return false
	}
	return p1.IsSatisfiedByEntry(entry)
}

// IsSatisfiedByEntry returns true if e satisfies the predicate, false otherwise
func (p1 EntryPredicate) IsSatisfiedByEntry(e Entry) bool {
	return p1.P(e) && p1.RequiredActions.IsSubset(toSet(e.Actions))
}

// EntryPredicateParser parses Entry predicates
type EntryPredicateParser func(tokens []string) (EntryPredicate, []string, error)

// Parse parses an EntryPredicate from the given input.
func (parser EntryPredicateParser) Parse(tokens []string) (predicate.Predicate, []string, error) {
	return parser(tokens)
}
