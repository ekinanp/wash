package meta

import (
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
)

/*
If metadata predicates are constructed on metadata values, then metadata
schema predicates are constructed on metadata schemas. Thus, one would expect
that metadata schema predicate parsing is symmetric with metadata predicate
parsing, where instead of walking the metadata values, we walk the metadata
schema. Unfortunately, metadata schemas are JSON schemas. Walking a JSON schema
is more complicated than walking a JSON object because there's a lot more rules
associated with a JSON schema than a JSON object. Thus, it is easier to delegate
to a JSON schema validator.

Consider the expression ".key1 .key2 5 -o 6". This reads "return
true if m['key1']['key2'] == 5 OR m['key1'] == 6". The schema predicate
would return true if "m['key1']['key2'] == number OR m['key1'] == number".
Since we don't care about primitive types, this reduces to
"m['key1']['key2'] == primitive_type OR m['key1'] == primitive_type".
If we normalize all primitive types to the "null" type, then our final schema
predicate is "m['key1']['key2'] == null OR m['key1'] == null". Now if we let
LHS = {"KEY1":{"KEY2":null}} represent the LHS' JSON serialization, and
RHS = {"KEY1":null} represent the RHS' JSON serialization, then our schema
predicate would return true iff the JSON schema validator returned true for
the LHS OR if the validator returned true for the RHS.

Generating the JSON object for a key sequence is tricky because unlike metadata
predicates, child nodes need to know the current key sequence. For example,
in the expression ".key1 .key2 5", we want our generated JSON object to be
{"KEY1": {"KEY2": null}}, and we want this object to be generated by the "5"
node since that is where the key sequence ends. Since our schema predicate
consists of validating JSON objects against the metadata schema, and since those
JSON objects are generated from key sequences, this implies that schema predicates
are generated from key sequences. That is why every schemaP is associated with a
key sequence.

NOTE: We'll need to munge the JSON schema to ensure that we get the right validation.
That munging is done by the schema type.
*/

type schemaPredicate interface {
	predicate.Predicate
	updateKS(func(keySequence) keySequence)
}

type schemaP struct {
	ks       keySequence
	p        func(schema) bool
	children []schemaPredicate
}

func newSchemaP() schemaPredicate {
	sp := &schemaP{}
	sp.p = func(s schema) bool {
		return s.IsValidKeySequence(sp.ks)
	}
	return sp
}

// And returns p1 && p2
func (p1 *schemaP) And(p2 predicate.Predicate) predicate.Predicate {
	sp2 := p2.(schemaPredicate)
	return &schemaP{
		p: func(s schema) bool {
			return p1.IsSatisfiedBy(s) && sp2.IsSatisfiedBy(s)
		},
		children: []schemaPredicate{p1, sp2},
	}
}

// Or returns p1 || p2
func (p1 *schemaP) Or(p2 predicate.Predicate) predicate.Predicate {
	sp2 := p2.(schemaPredicate)
	return &schemaP{
		p: func(s schema) bool {
			return p1.IsSatisfiedBy(s) || sp2.IsSatisfiedBy(s)
		},
		children: []schemaPredicate{p1, sp2},
	}
}

func (p1 *schemaP) Negate() predicate.Predicate {
	return &schemaP{
		p: func(s schema) bool {
			return !p1.IsSatisfiedBy(s)
		},
	}
}

func (p1 *schemaP) IsSatisfiedBy(v interface{}) bool {
	s, ok := v.(schema)
	if !ok {
		return false
	}
	return p1.p(s)
}

func (p1 *schemaP) updateKS(updateFunc func(keySequence) keySequence) {
	p1.ks = updateFunc(p1.ks)
	for _, child := range p1.children {
		child.updateKS(updateFunc)
	}
}

// schemaPLeaf is a base class for primitive schemaPs + the "-empty"
// predicate's schemaP. The reason for the separate type is to
// centralize the negation semantics for a schemaPLeaf.
type schemaPLeaf struct {
	schemaPredicate
}

func newSchemaPLeaf() *schemaPLeaf {
	return &schemaPLeaf{
		schemaPredicate: newSchemaP(),
	}
}

/*
"Not(schemaPLeaf) == schemaPLeaf". To see why, consider the negation of
the leaf's predicate counterpart (like "! 5"). Since the predicate's
negation still returns false for a mis-typed value, and since schemaPs
operate at the type-level, both these conditions imply that the type-level
predicate, and hence the schemaP, does not change when the leaf's predicate
counterpart is negated. In our example, "! 5" only returns true for numeric
values. Thus, its corresponding schemaP should still expect a numeric value
(specifically a "primitive value" since primitive types are normalized to
"null"). Similarly, "! -empty" only returns true for arrays/objects, so its
schemaP should still expect an array/object.
*/
func (p1 *schemaPLeaf) Negate() predicate.Predicate {
	return p1
}

// primitiveSchemaP represents a schema predicate for primitive types
type primitiveSchemaP struct {
	*schemaPLeaf
}

func newPrimitiveSchemaP() *primitiveSchemaP {
	psp := &primitiveSchemaP{
		schemaPLeaf: newSchemaPLeaf(),
	}
	psp.updateKS(func(ks keySequence) keySequence {
		return ks.EndsWithPrimitiveValue()
	})
	return psp
}
