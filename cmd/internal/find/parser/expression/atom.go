package expression

import "github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
import "github.com/puppetlabs/wash/cmd/internal/find/parser/errz"
import "fmt"

func notOpParser(parser Parser) predicate.Parser {
	return predicate.ToParser(func(tokens []string) (predicate.Predicate, []string, error) {
		notToken := tokens[0]
		if notToken != "!" && notToken != "-not" {
			return nil, nil, errz.NewMatchError("expected ! or -not")
		}
		tokens = tokens[1:]
		if len(tokens) == 0 {
			return nil, nil, fmt.Errorf("%v: no following expression", notToken)
		}
		if tokens[0] == ")" {
			if !parser.insideParens() {
				return nil, nil, fmt.Errorf("): no beginning '('")
			}
			return nil, nil, fmt.Errorf("%v: no following expression", notToken)
		}
		p, tokens, err := parser.atom().Parse(tokens)
		if err != nil && !IsIncompleteOperatorError(err) {
			if errz.IsMatchError(err) {
				err = IncompleteOperatorError{
					fmt.Sprintf("%v: no following expression", notToken),
				}
			}
			return nil, nil, err
		}
		return p.Negate(), tokens, err
	})
}

// Parenthesize returns a predicate parser that only parses parenthesized expressions.
// The expressions themselves are parsed by the given parser. Note that the parser
// returned by Parenthesize mutates the passed-in parser's state.
//
// If Parser#Parse returns an EmptyExpressionError, then Parenthesize also returns an
// EmptyExpressionError,
func Parenthesize(parser Parser) predicate.Parser {
	return predicate.ToParser(func(tokens []string) (predicate.Predicate, []string, error) {
		if len(tokens) == 0 {
			return nil, nil, errz.NewMatchError("expected an '('")
		}
		if tokens[0] == ")" {
			return nil, nil, fmt.Errorf("): no beginning '('")
		}
		if tokens[0] != "(" {
			return nil, nil, errz.NewMatchError("expected an '('")
		}
		tokens = tokens[1:]
		parser.openParens()
		// Save the current evaluation stack. We will restore it after parsing
		// the parenthesized expression
		stack := parser.stack()
		defer func() {
			parser.setStack(stack)
			parser.closeParens()
		}()
		p, tokens, err := parser.Parse(tokens)
		if err != nil && !IsEmptyExpressionError(err) {
			return p, tokens, err
		}
		// err == nil || IsEmptyExpressionError(err)
		if len(tokens) == 0 || tokens[0] != ")" {
			return nil, nil, fmt.Errorf("(: missing closing ')'")
		}
		if IsEmptyExpressionError(err) {
			return nil, nil, fmt.Errorf("(): empty inner expression")
		}
		tokens = tokens[1:]
		return p, tokens, err
	})
}
