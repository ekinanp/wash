package plugin

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ekinanp/jsonschema"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

const registrySchemaLabel = "mountpoint"

// TypeID returns the entry's type ID. It is needed by the API,
// so plugin authors should ignore this.
func TypeID(e Entry) string {
	pluginName := pluginName(e)
	rawTypeID := RawTypeID(e)
	if pluginName == "" {
		// e is the plugin registry
		return rawTypeID
	}
	return namespace(pluginName, rawTypeID)
}

func schema(e Entry) (*EntrySchema, error) {
	switch t := e.(type) {
	case externalPlugin:
		graph, err := t.SchemaGraph()
		if err != nil {
			return nil, err
		}
		if graph == nil {
			return nil, nil
		}
		s := NewEntrySchema(e, "foo")
		s.graph = graph
		entrySchemaV, _ := s.graph.Get(TypeID(e))
		// Nodes in the graph can only set properties on entrySchema, so only copy that.
		s.entrySchema = entrySchemaV.(EntrySchema).entrySchema
		return s, nil
	case *Registry:
		// The plugin registry includes external plugins, whose schema call can
		// error. Thus, it needs to be treated differently from the other core
		// plugins. Here, the logic is to merge each root's schema graph with the
		// registry's schema graph.
		schema := NewEntrySchema(t, registrySchemaLabel).
			IsSingleton().
			SetDescription(registryDescription)
		schema.graph = linkedhashmap.New()
		schema.graph.Put(TypeID(t), &schema.entrySchema)
		for _, root := range t.pluginRoots {
			childSchema, err := Schema(root)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve the %v plugin's schema: %v", root.eb().name, err)
			}
			if childSchema == nil {
				// The plugin doesn't have a schema, which means it's an external plugin.
				// Create a schema for root so that `stree <mountpoint>` can still display
				// it.
				childSchema = NewEntrySchema(root, CName(root))
			}
			childSchema.IsSingleton()
			schema.Children = append(schema.Children, TypeID(childSchema.entry))
			childGraph := childSchema.graph
			if childGraph == nil {
				// This is a core-plugin
				childGraph = linkedhashmap.New()
				childSchema.fill(childGraph)
			}
			childGraph.Each(func(key interface{}, value interface{}) {
				schema.graph.Put(key, value)
			})
		}
		return schema, nil
	default:
		// e is a core-plugin
		return e.Schema(), nil
	}
}

type entrySchema struct {
	Label                 string         `json:"label"`
	Description           string         `json:"description,omitempty"`
	Singleton             bool           `json:"singleton"`
	Signals               []SignalSchema `json:"signals,omitempty"`
	Actions               []string       `json:"actions"`
	PartialMetadataSchema *JSONSchema    `json:"partial_metadata_schema"`
	MetadataSchema        *JSONSchema    `json:"metadata_schema"`
	Children              []string       `json:"children"`
}

// EntrySchema represents an entry's schema. Use plugin.NewEntrySchema
// to create instances of these objects.
type EntrySchema struct {
	// This pattern's a nice way of making JSON marshalling/unmarshalling
	// easy without having to export these fields via the godocs. The latter
	// is good because plugin authors should use the builders when setting them
	// (so that we maintain a consistent API for e.g. metadata schemas).
	//
	// This pattern was obtained from https://stackoverflow.com/a/11129474
	entrySchema
	partialMetadataSchemaObj interface{}
	metadataSchemaObj        interface{}
	// Store the entry so that we can compute its type ID and, if the entry's
	// a core plugin entry, enumerate its child schemas when marshaling its
	// schema.
	entry Entry
	// graph is set by external plugins
	graph *linkedhashmap.Map
}

// NewEntrySchema returns a new EntrySchema object with the specified label.
func NewEntrySchema(e Entry, label string) *EntrySchema {
	if len(label) == 0 {
		panic("plugin.NewEntrySchema called with an empty label")
	}
	s := &EntrySchema{
		entrySchema: entrySchema{
			Label:   label,
			Actions: SupportedActionsOf(e),
		},
		// The partial metadata's empty by default
		partialMetadataSchemaObj: struct{}{},
		entry:                    e,
	}
	return s
}

// MarshalJSON marshals the entry's schema to JSON. It takes
// a value receiver so that the entry schema's still marshalled
// when it's referenced as an interface{} object. See
// https://stackoverflow.com/a/21394657 for more details.
//
// Note that UnmarshalJSON is not implemented since that is not
// how plugin.EntrySchema objects are meant to be used.
func (s EntrySchema) MarshalJSON() ([]byte, error) {
	if s.entry == nil {
		// Nodes in the external plugin graph don't use NewEntrySchema, they directly set the
		// undocumented fields of EntrySchema. Since graph and entry won't be set - and this is
		// part of a graph already - directly serialize entrySchema instead of using the graph.
		return json.Marshal(s.entrySchema)
	}

	graph := s.graph
	if graph == nil {
		if _, ok := s.entry.(externalPlugin); ok {
			// We should never hit this code-path because external plugins with
			// unknown schemas will return a nil schema. Thus, EntrySchema#MarshalJSON
			// will never be invoked.
			msg := fmt.Sprintf(
				"s.MarshalJSON: called with a nil graph for external plugin entry %v (type ID %v)",
				CName(s.entry),
				TypeID(s.entry),
			)
			panic(msg)
		}
		// We're marshalling a core plugin entry's schema. Note that the reason
		// we use an ordered map is to ensure that the first key in the marshalled
		// schema corresponds to s.
		graph = linkedhashmap.New()
		s.fill(graph)
	}
	return graph.ToJSON()
}

// SetDescription sets the entry's description.
func (s *EntrySchema) SetDescription(description string) *EntrySchema {
	s.entrySchema.Description = strings.Trim(description, "\n")
	return s
}

// IsSingleton marks the entry as a singleton entry.
func (s *EntrySchema) IsSingleton() *EntrySchema {
	s.entrySchema.Singleton = true
	return s
}

// AddSignal adds the given signal to s' supported signals. See https://puppetlabs.github.io/wash/docs#signal
// for a list of common signal names. You should try to re-use these names if you can.
func (s *EntrySchema) AddSignal(name string, description string) *EntrySchema {
	return s.addSignalSchema(name, "", description)

}

// AddSignalGroup adds the given signal group to s' supported signals
func (s *EntrySchema) AddSignalGroup(name string, regex string, description string) *EntrySchema {
	if len(regex) <= 0 {
		panic("s.AddSignalGroup: received empty regex")
	}
	return s.addSignalSchema(name, regex, description)
}

func (s *EntrySchema) addSignalSchema(name string, regex string, description string) *EntrySchema {
	if len(name) <= 0 {
		panic("s.addSignalSchema: received empty name")
	}
	if len(description) <= 0 {
		panic("s.addSignalSchema: received empty description")
	}
	schema := SignalSchema{
		signalSchema: signalSchema{
			Name:        name,
			Regex:       regex,
			Description: description,
		},
	}
	err := schema.normalize()
	if err != nil {
		msg := fmt.Sprintf("s.addSignalSchema: received invalid regex: %v", err)
		panic(msg)
	}
	s.Signals = append(s.Signals, schema)
	return s
}

// SetPartialMetadataSchema sets the partial metadata's schema. obj is an empty
// struct that will be marshalled into a JSON schema. SetPartialMetadataSchema
// will panic if obj is not a struct.
func (s *EntrySchema) SetPartialMetadataSchema(obj interface{}) *EntrySchema {
	// We need to know if s.entry has any wrapped types in order to correctly
	// compute the schema. However that information is known when s.fill() is
	// called. Thus, we'll set the schema object to obj so s.fill() can properly
	// calculate the schema.
	s.partialMetadataSchemaObj = obj
	return s
}

// SetMetadataSchema sets Entry#Metadata's schema. obj is an empty struct that will be
// marshalled into a JSON schema. SetMetadataSchema will panic if obj is not a struct.
//
// NOTE: Only use SetMetadataSchema if you're overriding Entry#Metadata. Otherwise, use
// SetPartialMetadataSchema.
func (s *EntrySchema) SetMetadataSchema(obj interface{}) *EntrySchema {
	// See the comments in SetPartialMetadataSchema to understand why this line's necessary
	s.metadataSchemaObj = obj
	return s
}

func (s *EntrySchema) fill(graph *linkedhashmap.Map) {
	// Fill-in the partial metadata + metadata schemas
	var err error
	if s.partialMetadataSchemaObj != nil {
		s.entrySchema.PartialMetadataSchema, err = s.schemaOf(s.partialMetadataSchemaObj)
		if err != nil {
			s.fillPanicf("bad value passed into SetPartialMetadataSchema: %v", err)
		}
	}
	if s.metadataSchemaObj != nil {
		s.entrySchema.MetadataSchema, err = s.schemaOf(s.metadataSchemaObj)
		if err != nil {
			s.fillPanicf("bad value passed into SetMetadataSchema: %v", err)
		}
	}
	graph.Put(TypeID(s.entry), &s.entrySchema)

	// Fill-in the children
	if !ListAction().IsSupportedOn(s.entry) {
		return
	}
	// "sParent" is read as "s.parent"
	sParent := s.entry.(Parent)
	children := sParent.ChildSchemas()
	if children == nil {
		s.fillPanicf("ChildSchemas() returned nil")
	}
	for _, child := range children {
		if child == nil {
			s.fillPanicf("found a nil child schema")
		}
		// The ID here is meaningless. We only set it so that TypeID can get the
		// plugin name
		child.entry.eb().id = s.entry.eb().id
		childTypeID := TypeID(child.entry)
		s.entrySchema.Children = append(s.Children, childTypeID)
		if _, ok := graph.Get(childTypeID); ok {
			continue
		}
		passAlongWrappedTypes(sParent, child.entry)
		child.fill(graph)
	}
}

// This helper's used by CachedList + EntrySchema#fill(). The reason for
// the helper is because /fs/schema uses repeated calls to CachedList when
// fetching the entry, so we need to pass-along the wrapped types when
// searching for it. However, Parent#ChildSchemas uses empty Entry objects
// that do not go through CachedList (by definition). Thus, the entry found
// by /fs/schema needs to pass its wrapped types along to the children to
// determine their metadata schemas. This is done in s.fill().
func passAlongWrappedTypes(p Parent, child Entry) {
	var wrappedTypes SchemaMap
	if root, ok := child.(HasWrappedTypes); ok {
		wrappedTypes = root.WrappedTypes()
	} else {
		wrappedTypes = p.eb().wrappedTypes
	}
	child.eb().wrappedTypes = wrappedTypes
}

// Helper that wraps the common code shared by the SetMeta*Schema methods
func (s *EntrySchema) schemaOf(obj interface{}) (*JSONSchema, error) {
	typeMappings := make(map[reflect.Type]*jsonschema.Type)
	for t, s := range s.entry.eb().wrappedTypes {
		typeMappings[reflect.TypeOf(t)] = s.Type
	}
	r := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		// Setting this option ensures that the schema's root is obj's
		// schema instead of a reference to a definition containing obj's
		// schema. This way, we can validate that "obj" is a JSON object's
		// schema. Otherwise, the check below will always fail.
		ExpandedStruct: true,
		TypeMappings:   typeMappings,
	}
	schema := r.Reflect(obj)
	if schema.Type.Type != "object" {
		return nil, fmt.Errorf("expected a JSON object but got %v", schema.Type.Type)
	}
	return schema, nil
}

// Helper for s.fill(). We make it a separate method to avoid re-creating
// closures for each recursive s.fill() call.
func (s *EntrySchema) fillPanicf(format string, a ...interface{}) {
	formatStr := fmt.Sprintf("s.fill (%v): %v", TypeID(s.entry), format)
	msg := fmt.Sprintf(formatStr, a...)
	panic(msg)
}

func pluginName(e Entry) string {
	// Using ID(e) will panic if e.id() is empty. The latter's possible
	// via something like "Schema(registry) => TypeID(Root)", where
	// CachedList(registry) was not yet called. This can happen if, for
	// example, the user starts the Wash shell and runs `stree`.
	trimmedID := strings.Trim(e.eb().id, "/")
	if trimmedID == "" {
		switch e.(type) {
		case Root:
			return CName(e)
		case *Registry:
			return ""
		default:
			// e has no ID. This is possible if e's from the apifs package. For now,
			// it is enough to return "__apifs__" here because this is an unlikely
			// edge case.
			//
			// TODO: Panic here once https://github.com/puppetlabs/wash/issues/438
			// is resolved.
			return "__apifs__"
		}
	}
	segments := strings.SplitN(trimmedID, "/", 2)
	return segments[0]
}

// RawTypeID returns e's raw type ID. Plugin authors should ignore this.
func RawTypeID(e Entry) string {
	switch t := e.(type) {
	case externalPlugin:
		rawTypeID := t.RawTypeID()
		if rawTypeID == "" {
			rawTypeID = "unknown"
		}
		return rawTypeID
	default:
		// e is either a core plugin entry or the plugin registry itself
		reflectType := unravelPtr(reflect.TypeOf(e))
		return reflectType.PkgPath() + "/" + reflectType.Name()
	}
}

func namespace(pluginName string, rawTypeID string) string {
	return pluginName + "::" + rawTypeID
}

func unravelPtr(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return unravelPtr(t.Elem())
	}
	return t
}
