package external

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/puppetlabs/wash/plugin"
)

// pluginRoot represents an external plugin's root.
type pluginRoot struct {
	pluginEntry
}

// Init initializes the external plugin root
func (r *pluginRoot) Init(cfg map[string]interface{}) error {
	if cfg == nil {
		cfg = make(map[string]interface{})
	}
	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshal plugin config %v into JSON: %v", cfg, err)
	}

	// Give external plugins about five-seconds to finish their
	// initialization
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	inv, err := r.script.InvokeAndWait(ctx, "init", nil, string(cfgJSON))
	if err != nil {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out while waiting for init to finish")
		default:
			return err
		}
	}
	var decodedRoot decodedExternalPluginEntry
	if err := json.Unmarshal(inv.Stdout().Bytes(), &decodedRoot); err != nil {
		return newStdoutDecodeErr(
			context.Background(),
			"the plugin root",
			err,
			inv,
			"{}",
		)
	}

	// Fill in required fields with data we already know.
	if decodedRoot.Name == "" {
		decodedRoot.Name = r.Name()
	} else if decodedRoot.Name != r.Name() {
		panic(fmt.Sprintf(`plugin root's name must match the basename (without extension) of %s
it's safe to omit name from the response to 'init'`, r.script.Path()))
	}
	if decodedRoot.Methods == nil {
		decodedRoot.Methods = rawMethods(`"list"`)
	}
	entry, err := decodedRoot.toExternalPluginEntry(context.Background(), false, true)
	if err != nil {
		return err
	}
	if !plugin.ListAction().IsSupportedOn(entry) {
		panic(fmt.Sprintf("plugin root for %s must implement 'list'", r.script.Path()))
	}
	script := r.script
	r.pluginEntry = *entry
	r.pluginEntry.script = script

	// Fill in the schema graph if provided
	if val := r.methods["schema"].tupleValue; val != nil {
		r.schemaGraphs = r.partitionSchemaGraph(val.(*linkedhashmap.Map))
	}

	return nil
}

func (r *pluginRoot) WrappedTypes() plugin.SchemaMap {
	// This only makes sense for core plugins because it is a Go-specific
	// limitation.
	return nil
}

// partitionSchemaGraph partitions graph into a map of <type_id> => <schema_graph>
func (r *pluginRoot) partitionSchemaGraph(graph *linkedhashmap.Map) map[string]*linkedhashmap.Map {
	var populate func(*linkedhashmap.Map, string, map[string]interface{}, map[string]bool)
	populate = func(g *linkedhashmap.Map, typeID string, node map[string]interface{}, visited map[string]bool) {
		if visited[typeID] {
			return
		}
		g.Put(typeID, node)
		visited[typeID] = true
		var nodeSchema plugin.EntrySchema
		_ = nodeSchema.FromMap(node)
		for _, childTypeID := range nodeSchema.Children {
			childNode, ok := graph.Get(childTypeID)
			if !ok {
				msg := fmt.Sprintf("plugin.partitionSchemaGraph: expected child %v to be present in the graph", childTypeID)
				panic(msg)
			}
			populate(g, childTypeID, childNode.(map[string]interface{}), visited)
		}
	}

	schemaGraphs := make(map[string]*linkedhashmap.Map)
	graph.Each(func(key interface{}, value interface{}) {
		g := linkedhashmap.New()
		populate(g, key.(string), value.(map[string]interface{}), make(map[string]bool))
		schemaGraphs[key.(string)] = g
	})

	return schemaGraphs
}

func rawMethods(strings ...string) []json.RawMessage {
	raw := make([]json.RawMessage, len(strings))
	for i, s := range strings {
		raw[i] = []byte(s)
	}
	return raw
}
