package meta

import (
	"fmt"

	"github.com/puppetlabs/wash/cmd/internal/find/params"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/errz"
	"github.com/puppetlabs/wash/cmd/internal/find/parser/predicate"
	"github.com/puppetlabs/wash/cmd/internal/find/primary/numeric"
	"github.com/puppetlabs/wash/munge"
)

// TimePredicate => (+|-)? Duration
// Duration      => numeric.DurationRegex | '{' numeric.DurationRegex '}'
func parseTimePredicate(tokens []string) (predicate.Predicate, []string, error) {
	if params.ReferenceTime.IsZero() {
		panic("meta.parseTimePredicate called without setting params.ReferenceTime")
	}
	if len(tokens) == 0 {
		return nil, nil, errz.NewMatchError("expected a +, -, or a digit")
	}
	token := tokens[0]
	p, parserID, err := numeric.ParsePredicate(
		token,
		numeric.ParseDuration,
		numeric.Bracket(numeric.ParseDuration),
	)
	if err != nil {
		if errz.IsMatchError(err) {
			msg := fmt.Sprintf("expected a duration but got %v", token)
			return nil, nil, errz.NewMatchError(msg)
		}
		// err is a parse error, so return it.
		return nil, nil, err
	}
	subFromReferenceTime := true
	if parserID == 1 {
		// User passed-in something like +{1h}. This means they want to
		// base the predicate off of 'timeV - StartTime' instead of
		// 'StartTime - timeV'.
		subFromReferenceTime = false
	}
	return timeP(subFromReferenceTime, p), tokens[1:], nil
}

func timeP(subFromReferenceTime bool, p numeric.Predicate) predicate.Predicate {
	return &timePredicate{
		genericPredicate: func(v interface{}) bool {
			timeV, err := munge.ToTime(v)
			if err != nil {
				return false
			}
			var diff int64
			if subFromReferenceTime {
				diff = int64(params.ReferenceTime.Sub(timeV))
			} else {
				diff = int64(timeV.Sub(params.ReferenceTime))
			}
			if diff < 0 {
				// Time predicates query either the past or the future, but not both.
				// For example, +1h means "more than one hour ago", which queries the
				// past. diff < 0 there means timeV is from the future, so the query
				// doesn't make sense. Similarly, +{1h} means "more than one hour from now",
				// which queries the future. diff < 0 there means timeV is from the past,
				// so the query doesn't make sense. Thus, diff < 0 is a time-mismatch.
				// Time-mismatches always return false.
				return false
			}
			return p(diff)
		},
		subFromReferenceTime: subFromReferenceTime,
		p: p,
	}
}

type timePredicate struct {
	genericPredicate
	subFromReferenceTime bool
	p numeric.Predicate
}

func (tp *timePredicate) Negate() predicate.Predicate {
	return timeP(tp.subFromReferenceTime, tp.p.Negate().(numeric.Predicate))
}
