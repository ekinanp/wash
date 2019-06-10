package plugin

import (
	"reflect"
	"strings"
)

// EntrySchema represents an entry's schema.
type EntrySchema struct {
	Type      string        `json:"type"`
	ShortType string        `json:"short_type"`
	Singleton bool          `json:"singleton"`
	Actions   []string      `json:"actions"`
	Children  []EntrySchema `json:"children"`
	entry     Entry
}

// ChildSchemas is a helper that's used to implement Parent#ChildSchemas
func ChildSchemas(childTemplates ...Entry) []EntrySchema {
	var schemas []EntrySchema
	for _, template := range childTemplates {
		schemas = append(schemas, schema(template, false))
	}
	return schemas
}

// Schema returns the entry's schema. Plugin authors should use plugin.ChildSchemas
// when implementing Parent#ChildSchemas. Using Schema to do this can cause infinite
// recursion if e's children have the same type as e, which can happen if e's e.g.
// a volume directory.
func Schema(e Entry) EntrySchema {
	return schema(e, true)
}

// Common helper for Schema and ChildSchema
func schema(e Entry, includeChildren bool) EntrySchema {
	// TODO: Handle external plugin schemas
	switch e.(type) {
	case *externalPluginRoot:
		return EntrySchema{
			Type: "external-plugin-root",
		}
	case *externalPluginEntry:
		return EntrySchema{
			Type: "external-plugin-entry",
		}
	}

	s := EntrySchema{
		Type:      strings.TrimPrefix(reflect.TypeOf(e).String(), "*"),
		ShortType: e.shortType(),
		Singleton: e.isSingleton(),
		Actions:   SupportedActionsOf(e),
		entry:     e,
	}
	if includeChildren {
		s.fillChildren(make(map[string]bool))
	}
	return s
}

func (s *EntrySchema) fillChildren(visited map[string]bool) {
	if !ListAction().IsSupportedOn(s.entry) {
		return
	}
	if visited[s.Type] {
		// This means that s' children can have s' type, which is
		// true if s is e.g. a volume directory.
		return
	}
	s.Children = s.entry.(Parent).ChildSchemas()
	visited[s.Type] = true
	for i, child := range s.Children {
		child.fillChildren(visited)
		// Need to re-assign because child is not a pointer,
		// so s.Children[i] won't be updated.
		s.Children[i] = child
	}
}