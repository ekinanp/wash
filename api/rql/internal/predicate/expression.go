package predicate

import (
	"time"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/internal"
	"github.com/puppetlabs/wash/api/rql/internal/predicate/expression"
	"github.com/shopspring/decimal"
)

// PE_UnsignedNumericPredicate returns a node representing a predicate expression
// (PE) of UnsignedNumericPredicates
func PE_UnsignedNumericPredicate() rql.NumericPredicate {
	return expression.New("UnsignedNumericPredicate", func() rql.ASTNode {
		return UnsignedNumeric("", decimal.Decimal{})
	}).(rql.NumericPredicate)
}

func PE_ValuePredicate() rql.ValuePredicate {
	return expression.New("ValuePredicate", func() rql.ASTNode {
		return internal.NewNonterminalNode(
			Object(),
			Array(),
			Null(),
			Boolean(false),
			NumericValue("", decimal.Decimal{}),
			TimeValue("", time.Time{}),
			StringValue(),
		)
	}).(rql.ValuePredicate)
}
