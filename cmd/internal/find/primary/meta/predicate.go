package meta

import (
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

/*
Predicate => ObjectPredicate |
             ArrayPredicate  |
             PrimitivePredicate
*/
func parsePredicate(tokens []string) (predicate.Predicate, []string, error) {
	cp := &predicate.CompositeParser{
		MatchErrMsg: "expected either a primitive, object, or array predicate",
		Parsers: []predicate.Parser{
			predicate.ToParser(parseObjectPredicate),
			predicate.ToParser(parseArrayPredicate),
			predicate.ToParser(parsePrimitivePredicate),
		},
	}
	return cp.Parse(tokens)
}

// Predicate is a wrapper to predicate.Predicate. It is only used for
// extracting the schemaP (to avoid a type switch).
type Predicate interface {
	predicate.Predicate
	schemaP() schemaPredicate
}

// genericPredicate represents a `meta` primary predicate "base" class.
// Child classes should only override the Negate method. Here's why:
//   * Some `meta` primary predicates perform type validation, returning
//     false for a mis-typed value. genericPredicate#Negate is a strict
//     negation, so it will return true for a mis-typed value. This is bad.
//
//   * Some of the more complicated predicates require additional negation
//     semantics. For example, ObjectPredicate returns false if the key does
//     not exist. A negated ObjectPredicate should also return false for this
//     case.
//
// Both of these issues are resolved if the child class overrides Negate.
type genericPredicate struct {
	P       func(interface{}) bool
	SchemaP schemaPredicate
}

func genericP(p func(interface{}) bool) genericPredicate {
	return genericPredicate{
		P: p,
		// Default to primitiveSchemaP since it's the common case.
		SchemaP: newPrimitiveSchemaP(),
	}
}

// And returns p1 && p2
func (p1 genericPredicate) And(p2 predicate.Predicate) predicate.Predicate {
	// ep2 => extracted p2
	ep2 := p2.(Predicate)
	return genericPredicate{
		P: func(v interface{}) bool {
			return p1.P(v) && p2.IsSatisfiedBy(v)
		},
		SchemaP: p1.SchemaP.And(ep2.schemaP()).(schemaPredicate),
	}
}

// Or returns p1 || p2
func (p1 genericPredicate) Or(p2 predicate.Predicate) predicate.Predicate {
	ep2 := p2.(Predicate)
	return genericPredicate{
		P: func(v interface{}) bool {
			return p1.P(v) || p2.IsSatisfiedBy(v)
		},
		SchemaP: p1.SchemaP.Or(ep2.schemaP()).(schemaPredicate),
	}
}

// Negate returns Not(p1)
func (p1 genericPredicate) Negate() predicate.Predicate {
	return genericPredicate{
		P: func(v interface{}) bool {
			return !p1.P(v)
		},
		SchemaP: p1.SchemaP.Negate().(schemaPredicate),
	}
}

// IsSatisfiedBy returns true if v satisfies the predicate, false otherwise
func (p1 genericPredicate) IsSatisfiedBy(v interface{}) bool {
	return p1.P(v)
}

func (p1 genericPredicate) schemaP() schemaPredicate {
	return p1.SchemaP
}
