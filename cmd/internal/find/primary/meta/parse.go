// Package meta contains all the parsing logic for the `meta` primary
package meta

import (
	"github.com/puppetlabs/wash/cmd/internal/find/types"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

// Parse is the meta primary's parse function.
func Parse(tokens []string) (predicate.Entry, []string, error) {
	p, tokens, err := parseObjectPredicate(tokens)
	if err != nil {
		return nil, nil, err
	}
	return func(e types.Entry) bool {
		mp := map[string]interface{}(e.Attributes.Meta())
		return p(mp)
	}, tokens, nil
}
