package parser

import (
	"github.com/puppetlabs/wash/cmd/internal/find/types"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

// Result represents the result of parsing `wash find`'s
// arguments.
type Result struct {
	Path      string
	Options   types.Options
	Predicate predicate.Entry
}

/*
Parse parses `wash find`'s arguments, returning the result.
`wash find`'s arguments are specified as "[path] [options] [expression]"
*/
func Parse(args []string) (Result, error) {
	var err error
	r := Result{}
	r.Path, args = parsePath(args)
	r.Options, args, err = parseOptions(args)
	if err != nil {
		return r, err
	}
	r.Predicate, err = parseExpression(args)
	return r, err
}
