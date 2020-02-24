package predicate

import (
	"fmt"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/internal"
	"github.com/puppetlabs/wash/api/rql/internal/errz"
	"github.com/puppetlabs/wash/api/rql/internal/matcher"
)

// Common base class for Object/Array predicates

type collectionBase struct {
	rql.ValuePredicate
	// ctype => collectionType
	ctype            string
	elementPredicate rql.ValuePredicate
}

func (p *collectionBase) Marshal() interface{} {
	return []interface{}{p.ctype, p.ValuePredicate.Marshal()}
}

func (p *collectionBase) Unmarshal(input interface{}) error {
	isArraySize := p.ctype == "array"
	nt := internal.NewNonterminalNode(
		// SizePredicate
		&size{isArraySize: isArraySize, p: PE_UnsignedNumericPredicate()},
		// ElementPredicate
		p.elementPredicate,
	)
	nt.SetMatchErrMsg(fmt.Sprintf("expected a size predicate or a %v predicate", p.ctype))
	errMsgPrefix := fmt.Sprintf("must be formatted as [\"%v\", <size_predicate> | <%v_element_predicate>]", p.ctype, p.ctype)
	if !matcher.Array(matcher.Value(p.ctype))(input) {
		return errz.MatchErrorf(errMsgPrefix)
	}
	array := input.([]interface{})
	if len(array) > 2 {
		return fmt.Errorf(errMsgPrefix)
	}
	if len(array) < 2 {
		return fmt.Errorf("%v (missing the predicate part)", errMsgPrefix)
	}
	if err := nt.Unmarshal(array[1]); err != nil {
		return fmt.Errorf("error unmarshalling the %v predicate: %w", p.ctype, err)
	}
	p.ValuePredicate = nt.MatchedNode().(rql.ValuePredicate)
	return nil
}
