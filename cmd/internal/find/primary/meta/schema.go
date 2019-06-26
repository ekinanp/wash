package meta

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ekinanp/jsonschema"
	"github.com/puppetlabs/wash/plugin"
	"github.com/xeipuuv/gojsonschema"
)

// schema represents an entry metadata schema. It is a wrapper for
// plugin.JSONSchema.
type schema struct {
	loader gojsonschema.JSONLoader
}

func newSchema(s *plugin.JSONSchema) schema {
	// Given a meta primary key sequence like ".key1.key2 5", the schema predicate's
	// job is to check that the value "m['key1']['key2']" is possible. To do this
	// correctly, we need to munge s a bit before returning our schema instance. That's
	// what this code does.
	var mungeType func(t *jsonschema.Type)
	var mungeProperties func(map[string]*jsonschema.Type) map[string]*jsonschema.Type
	mungeProperties = func(properties map[string]*jsonschema.Type) map[string]*jsonschema.Type {
		// The meta primary searches for the first matching key where a matching key is the
		// first key s.t. upcase(matching_key) == upcase(key). Thus, all property names need
		// to be capitalized.
		upcasedProperties := make(map[string]*jsonschema.Type)
		for property, schema := range properties {
			if _, ok := upcasedProperties[property]; ok {
				continue
			}
			mungeType(schema)
			p := strings.ToUpper(property)
			upcasedProperties[p] = schema
		}
		return upcasedProperties
	}
	mungeType = func(t *jsonschema.Type) {
		if t == nil || len(t.Ref) > 0 {
			return
		}

		mungeType(t.Not)

		switch t.Type {
		case "array":
			mungeType(t.Items)
			mungeType(t.AdditionalItems)
			for _, items := range [][]*jsonschema.Type{t.AllOf, t.AnyOf, t.OneOf} {
				for _, schema := range items {
					mungeType(schema)
				}
			}
		case "object":
			// Metadata schemas should be simple, so we shouldn't have to worry
			// about dependencies (for now).
			t.Dependencies = nil
			t.Properties = mungeProperties(t.Properties)
			t.PatternProperties = mungeProperties(t.PatternProperties)
			if t.Required != nil || t.MinProperties >= 1 {
				// This ensures that the schema predicate (correctly) returns false for
				// something like "-m -empty". It is _very_ important to enforce this.
				// Otherwise since an entry's default metadata is an empty object, "-m -empty"
				// will trigger a bunch of API requests because its schema predicate will
				// be satisfied by e.g. "docker/containers", "aws/<profile>/resources", etc.
				t.MinProperties = 1
			}
			t.Required = nil
		default:
			// We've hit a primitive type. Normalize it by setting it to "null".
			t.Type = "null"
		}
	}

	mungeType(s.Type)
	for _, schema := range s.Definitions {
		mungeType(schema)
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		msg := fmt.Sprintf("Failed to marshal the metadata schema: %v", err)
		panic(msg)
	}
	return schema{
		gojsonschema.NewBytesLoader(bytes),
	}
}

// IsValidKeySequence returns true if ks is a valid key sequence in s, false
// otherwise.
func (s schema) IsValidKeySequence(ks keySequence) bool {
	r, err := gojsonschema.Validate(s.loader, gojsonschema.NewGoLoader(ks.toJSON()))
	if err != nil {
		msg := fmt.Sprintf("s.Validate: gojsonschema.Validate: returned an unexpected error: %v", err)
		panic(msg)
	}
	return r.Valid()
}
