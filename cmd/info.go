package cmd

import (
	"sync"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	cmdutil "github.com/puppetlabs/wash/cmd/util"
	"github.com/spf13/cobra"
	goyaml "gopkg.in/yaml.v2"
)

func infoCommand() *cobra.Command {
	use, aliases := generateShellAlias("info")
	infoCmd := &cobra.Command{
		Use:     use + " <path> [<path>]...",
		Aliases: aliases,
		Short:   "Wash's version of 'stat'. Prints the entries' info at the specified paths",
		RunE:    toRunE(infoMain),
	}
	infoCmd.Flags().StringP("output", "o", "yaml", "Set the output format (json, yaml, or text)")
	return infoCmd
}

func infoMain(cmd *cobra.Command, args []string) exitCode {
	paths := args
	if len(paths) == 0 {
		paths = []string{"."}
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		panic(err.Error())
	}

	marshaller, err := cmdutil.NewMarshaller(output)
	if err != nil {
		cmdutil.ErrPrintf(err.Error())
		return exitCode{1}
	}

	conn := cmdutil.NewClient()

	// Use a sorted map so that we can control how the information's
	// displayed.
	infoMap := infoResultMap{}

	// Fetch the data
	ec := 0
	var infoMapMux sync.Mutex
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			entry, err := conn.Info(path)
			if err != nil {
				ec = 1
				cmdutil.SafeErrPrintf("%v: %v\n", path, err)
				return
			}

			entryMap := orderedMap{linkedhashmap.New()}
			entryMap.Put("Name", entry.Name)
			entryMap.Put("CName", entry.CName)
			entryMap.Put("Actions", entry.Actions)
			entryMap.Put("Attributes", entry.Attributes.ToMap(false))

			infoMapMux.Lock()
			infoMap[path] = entryMap
			infoMapMux.Unlock()
		}(path)
	}
	wg.Wait()

	// Marshal the results
	var result interface{} = infoMap
	if len(paths) == 1 {
		// For a single path, it is enough to print the info object
		result = infoMap[paths[0]]
	}
	marshalledResult, err := marshaller.Marshal(result)
	if err != nil {
		cmdutil.ErrPrintf("error marshalling the info results: %v\n", err)
	} else {
		cmdutil.Print(marshalledResult)
	}

	// Return the exit code
	return exitCode{ec}
}

// This wrapped type's here to implement MarshalYAML because the (default)
// JSON unmarshalling causes the orderedMap values to be marshalled as a
// map[string]interface. This loses their insertion order.
type infoResultMap map[string]orderedMap

func (mp infoResultMap) MarshalYAML() ([]byte, error) {
	var yamlMap goyaml.MapSlice
	for path, infoResult := range mp {
		yamlMap = append(yamlMap, goyaml.MapItem{
			Key:   path,
			Value: infoResult.toMapSlice(),
		})
	}
	return goyaml.Marshal(yamlMap)
}

// This wrapped type's here because linkedhashmap doesn't implement the
// json.Marshaler and yaml.Marshaler interfaces.
type orderedMap struct {
	*linkedhashmap.Map
}

func (mp orderedMap) MarshalJSON() ([]byte, error) {
	return mp.ToJSON()
}

// We implement MarshalYAML to preserve each key's ordering.
func (mp orderedMap) MarshalYAML() ([]byte, error) {
	return goyaml.Marshal(mp.toMapSlice())
}

func (mp orderedMap) toMapSlice() goyaml.MapSlice {
	var yamlMap goyaml.MapSlice
	mp.Each(func(key interface{}, value interface{}) {
		yamlMap = append(yamlMap, goyaml.MapItem{
			Key:   key,
			Value: value,
		})
	})
	return yamlMap
}
