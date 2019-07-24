package primary

import (
	"fmt"

	"github.com/gobwas/glob"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
	"github.com/puppetlabs/wash/cmd/internal/find/types"
)

// Kind is the kind primary
//
// kindPrimary => -kind ShellGlob
//nolint
var Kind = Parser.add(&Primary{
	Description: "Returns true if the entry's kind matches the provided glob",
	name:        "kind",
	args:        "glob",
	parseFunc: func(tokens []string) (types.EntryPredicate, []string, error) {
		if len(tokens) == 0 {
			return nil, nil, fmt.Errorf("requires additional arguments")
		}
		g, err := glob.Compile(tokens[0])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid glob: %v", err)
		}
		return kindP(g, false), tokens[1:], nil
	},
})

func kindP(g glob.Glob, negated bool) types.EntryPredicate {
	p := kindPredicate{
		EntryPredicate: types.ToEntryP(func(e types.Entry) bool {
			// kind is a schema predicate, so the entry predicate should
			// always return true
			return true
		}),
		g: g,
	}
	p.SetSchemaP(types.ToEntrySchemaP(func(s *types.EntrySchema) bool {
		for _, path := range s.PathsToNode() {
			if g.Match(path) {
				if negated {
					return false
				}
				return true
			}
		}
		// Here we didn't match anything.
		return negated
	}))
	p.RequireSchema()
	return p
}

// The separate type's necessary to implement proper Negation semantics.
type kindPredicate struct {
	types.EntryPredicate
	g       glob.Glob
	negated bool
}

func (p kindPredicate) Negate() predicate.Predicate {
	return kindP(p.g, !p.negated)
}

const kindDetailedDescription = `
-kind glob

Returns true if the entry's size attribute is n 512-byte blocks,
rounded up to the nearest block. If n is suffixed with a unit,
then the raw size is compared to n scaled as:

c        character (byte)
k        kibibytes (1024 bytes)
M        mebibytes (1024 kibibytes)
G        gibibytes (1024 mebibytes)
T        tebibytes (1024 gibibytes)
P        pebibytes (1024 tebibytes)

If n is prefixed with a +/-, then the comparison returns true if
the size is greater-than/less-than n.

Examples:
  -size 2        Returns true if the entry's size is 2 512-byte blocks,
                 rounded up to the nearest block

  -size +2       Returns true if the entry's size is greater than 2
                 512-byte blocks, rounded up to the nearest block

  -size -2       Returns true if the entry's size is less than 2
                 512-byte blocks, rounded up to the nearest block

  -size +1k      Returns true if the entry's size is greater than 1
                 kibibyte
`
