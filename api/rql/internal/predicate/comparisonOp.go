package predicate

type ComparisonOp string

const (
	LT   ComparisonOp = "<"
	LTE               = "<="
	GT                = ">"
	GTE               = ">="
	EQL               = "="
	NEQL              = "!="
)

var comparisonOpMap = map[ComparisonOp]bool{
	LT:   true,
	LTE:  true,
	GT:   true,
	GTE:  true,
	EQL:  true,
	NEQL: true,
}
