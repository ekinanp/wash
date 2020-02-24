package predicate

import (
	"fmt"
	"strings"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/internal/errz"
)

func Object() rql.ValuePredicate {
	return &object{collectionBase{
		ctype:            "object",
		elementPredicate: &objectElement{p: PE_ValuePredicate()},
	}}
}

type object struct {
	collectionBase
}

var _ = rql.ValuePredicate(&object{})

type objectElement struct {
	key string
	p   rql.ValuePredicate
}

func (p *objectElement) Marshal() interface{} {
	return []interface{}{[]interface{}{"key", p.key}, p.p.Marshal()}
}

func (p *objectElement) Unmarshal(input interface{}) error {
	array, ok := input.([]interface{})
	formatErrMsg := "element predicate: must be formatted as [<element_selector>, PE ValuePredicate]"
	if !ok || len(array) < 1 {
		return errz.MatchErrorf(formatErrMsg)
	}
	keySelector, ok := array[0].([]interface{})
	if !ok || len(keySelector) < 1 || keySelector[0] != "key" {
		return errz.MatchErrorf(formatErrMsg)
	}
	if len(keySelector) > 2 {
		return fmt.Errorf(formatErrMsg)
	}
	if len(keySelector) < 2 {
		return fmt.Errorf("element predicate: missing the key")
	}
	key, ok := keySelector[1].(string)
	if !ok {
		return fmt.Errorf("element predicate: key must be a string, not %T", keySelector[1])
	}
	p.key = key
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

func (p *objectElement) ValueInDomain(v interface{}) bool {
	obj, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	key, ok := p.findMatchingKey(obj)
	return ok && p.p.ValueInDomain(obj[key])
}

func (p *objectElement) EvalValue(v interface{}) bool {
	obj := v.(map[string]interface{})
	k, _ := p.findMatchingKey(obj)
	return p.p.EvalValue(obj[k])
}

func (p *objectElement) findMatchingKey(mp map[string]interface{}) (string, bool) {
	upcasedKey := strings.ToUpper(p.key)
	for k := range mp {
		if strings.ToUpper(k) == upcasedKey {
			return k, true
		}
	}
	return "", false
}

var _ = rql.ValuePredicate(&objectElement{})
