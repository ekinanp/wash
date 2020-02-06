package predicate

import (
	"testing"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/ast/asttest"
	"github.com/puppetlabs/wash/api/rql/internal/predicate/expression"
	"github.com/stretchr/testify/suite"
)

type ArrayTestSuite struct {
	CollectionTestSuite
}

func (s *ArrayTestSuite) TestMarshal_ElementPredicate() {
	inputs := []interface{}{
		s.A("array", s.A("some", true)),
		s.A("array", s.A("all", true)),
		s.A("array", s.A(float64(0), true)),
	}
	for _, input := range inputs {
		p := Array()
		s.MUM(p, input)
		s.MTC(p, input)
	}
}

func (s *ArrayTestSuite) TestUnmarshalErrors_ElementPredicate() {
	// Start by testing the match errors
	p := Array().(*array).collectionBase.elementPredicate
	s.UMETC(p, s.A(), "formatted.*<element_selector>.*PE ValuePredicate", true)
	s.UMETC(p, s.A(true), "formatted.*<element_selector>.*PE ValuePredicate", true)
	s.UMETC(p, s.A("foo"), "formatted.*<element_selector>.*PE ValuePredicate", true)

	// Now test the syntax errors
	selectors := []interface{}{
		"some",
		"all",
		float64(1),
	}
	for _, selector := range selectors {
		s.UMETC(p, s.A(selector, s.A("string", s.A("=", "foo")), "bar"), "formatted.*<element_selector>.*PE ValuePredicate", false)
		s.UMETC(p, s.A(selector), "formatted.*<element_selector>.*PE ValuePredicate.*missing.*PE ValuePredicate", false)
	}
	s.UMETC(p, s.A(float64(-10), "foo"), "array.*index.*unsigned.*int", false)
}

func (s *ArrayTestSuite) TestValueInDomain_ElementPredicate() {
	p := Array()

	// Test "some", "all"
	for _, selector := range []interface{}{"some", "all"} {
		s.MUM(p, s.A("array", s.A(selector, true)))
		s.VIDFTC(p, "foo", true, map[string]interface{}{}, []interface{}{}, []interface{}{"foo"}, []interface{}{true, "foo"})
		s.VIDTTC(p, []interface{}{false}, []interface{}{false, false})
	}

	// Test nth
	s.MUM(p, s.A("array", s.A(float64(1), true)))
	s.VIDFTC(p, "foo", true, map[string]interface{}{}, []interface{}{}, []interface{}{"foo"}, []interface{}{"foo", "bar"}, []interface{}{true, "foo"})
	s.VIDTTC(p, []interface{}{"foo", false}, []interface{}{"foo", false, "bar"})
}

func (s *ArrayTestSuite) TestEvalValue_ElementPredicate() {
	p := Array()

	// Test "some"
	s.MUM(p, s.A("array", s.A("some", true)))
	s.EVFTC(p, []interface{}{false}, []interface{}{})
	s.EVTTC(p, []interface{}{true}, []interface{}{false, true})

	// Test "all"
	s.MUM(p, s.A("array", s.A("all", true)))
	s.EVFTC(p, []interface{}{false}, []interface{}{true, false})
	s.EVTTC(p, []interface{}{true}, []interface{}{true, true})

	// Test "n"
	s.MUM(p, s.A("array", s.A(float64(0), true)))
	s.EVFTC(p, []interface{}{"foo", "bar"}, []interface{}{false, true})
	s.EVTTC(p, []interface{}{true}, []interface{}{true, "foo"})
	// Add a case with a non-empty array
	s.MUM(p, s.A("array", s.A(float64(1), true)))
	s.EVFTC(p, []interface{}{true, false})
	s.EVTTC(p, []interface{}{false, true}, []interface{}{"foo", true})
}

func (s *ArrayTestSuite) TestExpression_AtomAndNot_ElementPredicate() {
	expr := expression.New("array", func() rql.ASTNode {
		return Array()
	})

	for _, selector := range []interface{}{"some", "all", float64(0)} {
		s.MUM(expr, s.A("array", s.A(selector, true)))
		s.VIDFTC(expr, "foo")
		s.VIDTTC(expr, []interface{}{false})
		s.EVFTC(expr, "foo", true, map[string]interface{}{}, []interface{}{}, []interface{}{false})
		s.EVTTC(expr, []interface{}{true})
		s.MUM(expr, s.A("NOT", s.A("array", s.A(selector, true))))
		s.VIDFTC(expr, "foo")
		s.VIDTTC(expr, []interface{}{false})
		s.EVTTC(expr, []interface{}{false})
		s.EVFTC(expr, "foo", true, map[string]interface{}{}, []interface{}{}, []interface{}{true})
	}

	// Assert that the unmarshaled atom doesn't implement the other *Predicate
	// interfaces
	s.MUM(expr, s.A("array", s.A("some", true)))
	s.AssertNotImplemented(
		expr,
		asttest.EntryPredicateC,
		asttest.EntrySchemaPredicateC,
		asttest.StringPredicateC,
		asttest.NumericPredicateC,
		asttest.TimePredicateC,
		asttest.ActionPredicateC,
	)
}

func TestArray(t *testing.T) {
	s := new(ArrayTestSuite)
	s.isArray = true
	suite.Run(t, s)
}
