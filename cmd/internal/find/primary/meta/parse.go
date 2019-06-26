// Package meta contains all the parsing logic for the `meta` primary
package meta

import (
	"github.com/puppetlabs/wash/cmd/internal/find/types"
)

// The functionality here is tested in primary/meta_test.go

// Parse is the meta primary's parse function.
func Parse(tokens []string) (*types.EntryPredicate, []string, error) {
	p, tokens, err := parseExpression(tokens)
	if err != nil {
		return nil, nil, err
	}
	entryP := types.ToEntryP(func(e types.Entry) bool {
		return p.IsSatisfiedBy(e.Metadata)
	})
	entryP.SchemaP = func(s *types.EntrySchema) bool {
		if s.MetadataSchema == nil {
			// The entry doesn't have a metadata schema, so
			// return true for now.
			//
			// TODO: Should we require metadata schemas? They can
			// be cumbersome to add for dynamic languages like
			// Ruby/Python
			//
			// TODO: Fix plugin.JSONSchema to default to an "empty"
			// object's schema if an entry doesn't have a metadata
			// schema. This makes serialization/de-serialization
			// easier.
			return true
		}
		return p.(Predicate).schemaP().IsSatisfiedBy(newSchema(s.MetadataSchema))
	}
	return entryP, tokens, nil
}
