package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ekinanp/jsonschema"
	"github.com/puppetlabs/wash/plugin"
	"github.com/xeipuuv/gojsonschema"
)

// MetadataSchema represents an entry metadata schema.
// It is a wrapper for plugin.JSONSchema.
type MetadataSchema struct {
	loader gojsonschema.JSONLoader
}

// NewMetadataSchema creates a new metadata schema object from
// the provided JSON schema.
func NewMetadataSchema(s *jsonschema.Schema) MetadataSchema {
	// The main thing we care about here is that all the values queried
	// by the meta primary exist. Thus, we eliminate required properties
	// and the "type" field for primitive types. We also capitalize property
	// names since the meta primary searches for the first matching key s.t.
	// upcase(key) == upcase(matching key).
	var mungeType func(t *jsonschema.Type)
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
			t.Required = nil
			t.AdditionalProperties = []byte("false")
			properties := make(map[string]*jsonschema.Type)
			for property, schema := range t.Properties {
				if _, ok := properties[property]; ok {
					continue
				}
				mungeType(schema)
				upcasedProperty := strings.ToUpper(property)
				properties[upcasedProperty] = schema
			}
			t.Properties = properties
		default:
			// We've hit a primitive type
			t.Type = ""
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
	return MetadataSchema{
		gojsonschema.NewBytesLoader(bytes),
	}
}

// Validate returns true if obj conforms to s, false otherwise.
func (s MetadataSchema) Validate(obj plugin.JSONObject) bool {
	r, err := gojsonschema.Validate(s.loader, gojsonschema.NewGoLoader(obj))
	if err != nil {
		msg := fmt.Sprintf("s.Validate: gojsonschema.Validate: returned an unexpected error: %v", err)
		panic(msg)
	}
	return r.Valid()
}
