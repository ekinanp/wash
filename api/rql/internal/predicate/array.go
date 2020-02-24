package predicate

import (
	"fmt"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/internal/errz"
)

func Array() rql.ValuePredicate {
	return &array{collectionBase{
		ctype:            "array",
		elementPredicate: &arrayElement{p: PE_ValuePredicate()},
	}}
}

type array struct {
	collectionBase
}

var _ = rql.ValuePredicate(&array{})

type arrayElement struct {
	selector interface{}
	p        rql.ValuePredicate
}

func (p *arrayElement) Marshal() interface{} {
	var marshalledSelector interface{}
	switch t := p.selector.(type) {
	case stringSelector:
		marshalledSelector = stringSelectorToStringMap[t]
	case int:
		marshalledSelector = float64(t)
	default:
		// Should never hit this code-path
		panic(fmt.Sprintf("Unknown selector %T", p.selector))
	}
	return []interface{}{marshalledSelector, p.p.Marshal()}
}

func (p *arrayElement) Unmarshal(input interface{}) error {
	array, ok := input.([]interface{})
	formatErrMsg := "element predicate: must be formatted as [<element_selector>, PE ValuePredicate]"
	if !ok || len(array) < 1 {
		return errz.MatchErrorf(formatErrMsg)
	}
	if firstElem, ok := array[0].(string); ok {
		if firstElem != "some" && firstElem != "all" {
			return errz.MatchErrorf(formatErrMsg)
		}
		if firstElem == "some" {
			p.selector = some
		} else {
			p.selector = all
		}
	} else {
		firstElem, ok := array[0].(float64)
		if !ok {
			return errz.MatchErrorf(formatErrMsg)
		}
		if firstElem < 0 {
			return fmt.Errorf("element predicate: array index must be an unsigned integer (> 0)")
		}
		p.selector = int(firstElem)
	}
	if len(array) > 2 {
		return fmt.Errorf(formatErrMsg)
	} else if len(array) < 2 {
		return fmt.Errorf("%v (missing PE ValuePredicate)", formatErrMsg)
	}
	if err := p.p.Unmarshal(array[1]); err != nil {
		return fmt.Errorf("element predicate: error unmarshalling the PE ValuePredicate: %w", err)
	}
	return nil
}

func (p *arrayElement) ValueInDomain(v interface{}) bool {
	array, ok := v.([]interface{})
	if !ok {
		return false
	}
	switch t := p.selector.(type) {
	default:
		// For the "some" and "all" selectors, we'd like to have some sort of
		// type-checking for the array elements without being too restrictive.
		// Since the # of elements (in general) is unknown, it is enough to
		// require at least one array element to be in the value predicate's
		// domain.
		for _, v := range array {
			if p.p.ValueInDomain(v) {
				return true
			}
		}
		return false
	case int:
		if len(array) <= t {
			return false
		}
		return p.p.ValueInDomain(array[t])
	}
}

func (p *arrayElement) EvalValue(v interface{}) bool {
	array := v.([]interface{})
	switch t := p.selector.(type) {
	case int:
		return p.p.EvalValue(array[t])
	case stringSelector:
		switch t {
		case some:
			for _, v := range array {
				if p.p.EvalValue(v) {
					return true
				}
			}
			return false
		case all:
			for _, v := range array {
				if !p.p.EvalValue(v) {
					return false
				}
			}
			return true
		default:
			// Should never hit this code path
			panic(fmt.Sprintf("Unknown string selector %v", t))
		}
	default:
		// Should never hit this code path
		panic(fmt.Sprintf("Unknown selector %T", p.selector))
	}
}

type stringSelector int8

const (
	some stringSelector = iota
	all
)

var stringSelectorToStringMap = map[stringSelector]string{
	some: "some",
	all:  "all",
}

var _ = rql.ValuePredicate(&arrayElement{})
