package primary

import (
	"testing"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/ast/asttest"
	"github.com/puppetlabs/wash/api/rql/internal/predicate"
	"github.com/puppetlabs/wash/api/rql/internal/predicate/expression"
	"github.com/stretchr/testify/suite"
)

type MetaTestSuite struct {
	asttest.Suite
}

func (s *MetaTestSuite) TestMarshal() {
	p := Meta(predicate.Object())
	input := s.A("meta", s.A("object", s.A(s.A("key", "foo"), true)))
	s.MUM(p, input)
	s.MTC(p, input)
}

func (s *MetaTestSuite) TestUnmarshalErrors() {
	n := Meta(predicate.Object())
	s.UMETC(n, "foo", `meta.*formatted.*"meta".*PE ObjectPredicate`, true)
	s.UMETC(n, s.A("foo", s.A("object", s.A(s.A("key", "foo"), true))), `meta.*formatted.*"meta".*PE ObjectPredicate`, true)
	s.UMETC(n, s.A("meta", "foo", "bar"), `meta.*formatted.*"meta".*PE ObjectPredicate`, false)
	s.UMETC(n, s.A("meta"), `meta.*formatted.*"meta".*PE ObjectPredicate.*missing.*PE ObjectPredicate`, false)
	s.UMETC(n, s.A("meta", s.A("object")), "meta.*PE ObjectPredicate.*element", false)
}

func (s *MetaTestSuite) TestEntryInDomain() {
	p := Meta(predicate.Object())
	s.MUM(p, s.A("meta", s.A("object", s.A(s.A("key", "foo"), true))))
	e := rql.Entry{}
	s.EIDFTC(p, e)
	e.Metadata = map[string]interface{}{"foo": false}
	s.EIDTTC(p, e)
}

func (s *MetaTestSuite) TestEvalEntry() {
	p := Meta(predicate.Object())
	s.MUM(p, s.A("meta", s.A("object", s.A(s.A("key", "foo"), true))))
	e := rql.Entry{}
	e.Metadata = map[string]interface{}{"foo": false}
	s.EEFTC(p, e)
	e.Metadata["foo"] = true
	s.EETTC(p, e)
}

func (s *MetaTestSuite) TestExpression_AtomAndNot() {
	expr := expression.New("meta", func() rql.ASTNode {
		return Meta(predicate.Object())
	})

	s.MUM(expr, s.A("meta", s.A("object", s.A(s.A("key", "foo"), true))))
	e := rql.Entry{}
	e.Metadata = map[string]interface{}{}
	s.EEFTC(expr, e)
	e.Metadata = map[string]interface{}{"foo": false}
	s.EEFTC(expr, e)
	e.Metadata["foo"] = true
	s.EETTC(expr, e)

	s.AssertNotImplemented(
		expr,
		asttest.ValuePredicateC,
		asttest.StringPredicateC,
		asttest.NumericPredicateC,
		asttest.TimePredicateC,
		asttest.ActionPredicateC,
	)

	s.MUM(expr, s.A("NOT", s.A("meta", s.A("object", s.A(s.A("key", "foo"), true)))))
	e.Metadata = map[string]interface{}{}
	s.EEFTC(expr, e)
	e.Metadata = map[string]interface{}{"foo": false}
	s.EETTC(expr, e)
	e.Metadata["foo"] = true
	s.EEFTC(expr, e)
}

func TestMeta(t *testing.T) {
	suite.Run(t, new(MetaTestSuite))
}
