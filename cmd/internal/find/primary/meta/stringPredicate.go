package meta

import (
	"fmt"

	"github.com/puppetlabs/wash/cmd/internal/find/parser/errz"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

// StringPredicate => [^-].*
func parseStringPredicate(tokens []string) (predicate.Generic, []string, error) {
	if len(tokens) == 0 || len(tokens[0]) == 0 {
		return nil, nil, errz.NewMatchError("expected a nonempty string")
	}
	token := tokens[0]
	if token[0] == '-' {
		// We impose this restriction to avoid conflicting with the
		// other primaries
		msg := fmt.Sprintf("%v begins with a '-'", token)
		return nil, nil, errz.NewMatchError(msg)
	}
	p := func(v interface{}) bool {
		strV, ok := v.(string)
		if !ok {
			return false
		}
		return strV == token
	}
	return p, tokens[1:], nil
}
