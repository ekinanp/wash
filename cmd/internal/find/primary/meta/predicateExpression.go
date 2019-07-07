package meta

import (
	"fmt"

	"github.com/puppetlabs/wash/cmd/internal/find/parser/expression"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

type predicateExpressionParser struct {
	expression.Parser
	isInnerExpression bool
}

func newPredicateExpressionParser(isInnerExpression bool) predicate.Parser {
	p := &predicateExpressionParser{
		Parser:            expression.NewParser(toPredicateParser(parsePredicate), &predicateAnd{}, &predicateOr{}),
		isInnerExpression: isInnerExpression,
	}
	if isInnerExpression {
		// Inner expressions are parenthesized so that they do not conflict with
		// the top-level predicate expression parser.
		return expression.Parenthesize(p)
	}
	return p
}

// This needs to return predicate.Predicate in order for it to use expression.Parenthesize
func (parser *predicateExpressionParser) Parse(tokens []string) (predicate.Predicate, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, expression.NewEmptyExpressionError("expected a predicate expression")
	}
	p, tks, err := parser.Parser.Parse(tokens)
	if err != nil && !expression.IsIncompleteOperatorError(err) {
		tkErr, ok := err.(expression.UnknownTokenError)
		if !ok {
			// We have a syntax error or an EmptyExpressionError. The latter's possible if
			// parser is an inner predicate expression and tokens is ")".
			return nil, tks, err
		}
		if p == nil || parser.isInnerExpression {
			// possible via something like "-size + 1"
			return nil, tks, fmt.Errorf("unknown predicate %v", tkErr.Token)
		}
		// This predicate expression parser is the top level predicate expression parser. This
		// means that we've finished parsing the `meta` primary's predicate expression.
		err = nil
	}
	return p, tks, err
}
