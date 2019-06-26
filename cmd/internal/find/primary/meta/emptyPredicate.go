package meta

import (
	"github.com/puppetlabs/wash/cmd/internal/find/parser/errz"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

// emptyPredicate => -empty
func parseEmptyPredicate(tokens []string) (predicate.Predicate, []string, error) {
	if len(tokens) == 0 || tokens[0] != "-empty" {
		return nil, nil, errz.NewMatchError("expected '-empty'")
	}
	return emptyP(false), tokens[1:], nil
}

func emptyP(negated bool) predicate.Predicate {
	return &emptyPredicate{
		genericPredicate: genericP(func(v interface{}) bool {
			switch t := v.(type) {
			case map[string]interface{}:
				if negated {
					return len(t) > 0
				}
				return len(t) == 0
			case []interface{}:
				if negated {
					return len(t) > 0
				}
				return len(t) == 0
			default:
				return false
			}
		}),
		negated: negated,
	}
}

type emptyPredicate struct {
	genericPredicate
	negated bool
}

func (p *emptyPredicate) Negate() predicate.Predicate {
	return emptyP(!p.negated)
}

type emptyPredicateSchemaP struct {
	*schemaPLeaf
	P schemaPredicate
}

func newEmptyPredicateSchemaP() *emptyPredicateSchemaP {
	esp := &emptyPredicateSchemaP{
		schemaPLeaf: newSchemaPLeaf(),
	}

	// An empty predicate's schemaP returns true iff the value's
	// an array OR an object.
	objP := newSchemaP()
	objP.updateKS(func(ks keySequence) keySequence {
		return ks.EndsWithObject()
	})
	arrayP := newSchemaP()
	arrayP.updateKS(func(ks keySequence) keySequence {
		return ks.EndsWithArray()
	})
	esp.P = objP.Or(arrayP).(schemaPredicate)

	return esp
}

func (p1 *emptyPredicateSchemaP) IsSatisfiedBy(v interface{}) bool {
	return p1.P.IsSatisfiedBy(v)
}
