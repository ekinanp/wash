package ast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/puppetlabs/wash/api/rql"
	"github.com/puppetlabs/wash/api/rql/ast/asttest"
	"github.com/stretchr/testify/suite"
)

// This file contains real-world meta primary test cases
//
// TODO: Once EntrySchemaInDomain's implemented, add schema predicate
// tests

type MetaPrimaryRealWorldTestSuite struct {
	asttest.Suite
	e rql.Entry
}

// TTC => TrueTestCase
func (s *MetaPrimaryRealWorldTestSuite) TTC(key interface{}, predicate interface{}) {
	q := Query()
	s.MUM(q, s.A("meta", s.A("object", s.A(s.A("key", key), predicate))))
	s.Suite.EETTC(q, s.e)
}

// FTC => FalseTestCase
func (s *MetaPrimaryRealWorldTestSuite) FTC(key interface{}, predicate interface{}) {
	q := Query()
	s.MUM(q, s.A("meta", s.A("object", s.A(s.A("key", key), predicate))))
	s.Suite.EEFTC(q, s.e)
}

func (s *MetaPrimaryRealWorldTestSuite) TestMetaPrimary_ValuePredicates() {
	s.TTC("architecture", s.A("string", s.A("=", "x86_64")))
	// False b/c x86_64 != x86_6
	s.FTC("architecture", s.A("string", s.A("=", "x86_6")))
	// False b/c it is negation
	s.FTC("architecture", s.A("NOT", s.A("string", s.A("=", "x86_64"))))
	// False b/c "architecture" is not a Boolean value
	s.FTC("architecture", true)
	s.FTC("architecture", s.A("NOT", true))

	s.TTC("blockDeviceMappings",
		s.A("array",
			s.A("some",
				s.A("object",
					s.A(s.A("key", "deviceName"),
						s.A("string", s.A("=", "/dev/sda1"))),
				),
			),
		),
	)
	// True b/c of negation (negating the "/dev/sda" part)
	s.TTC("blockDeviceMappings",
		s.A("NOT",
			s.A("array",
				s.A("some",
					s.A("object",
						s.A(s.A("key", "deviceName"),
							s.A("string", s.A("=", "/dev/sda"))),
					),
				),
			),
		),
	)
	// True b/c of negation (negating the "/dev/sda" part)
	s.TTC("blockDeviceMappings",
		s.A("array",
			s.A("some",
				s.A("NOT",
					s.A("object",
						s.A(s.A("key", "deviceName"),
							s.A("string", s.A("=", "/dev/sda"))),
					),
				),
			),
		),
	)
	// False b/c "deviceNam" is not a valid key
	s.FTC("blockDeviceMappings",
		s.A("array",
			s.A("some",
				s.A("object",
					s.A(s.A("key", "deviceNam"),
						s.A("string", s.A("=", "/dev/sda1"))),
				),
			),
		),
	)
	// Also false b/c "deviceNam" is not a valid key
	s.FTC("blockDeviceMappings",
		s.A("array",
			s.A("some",
				s.A("NOT",
					s.A("object",
						s.A(s.A("key", "deviceNam"),
							s.A("string", s.A("=", "/dev/sda"))),
					),
				),
			),
		),
	)

	s.TTC("cpuOptions",
		s.A("object",
			s.A(s.A("key", "coreCount"),
				s.A("number", s.A("=", "4")),
			),
		),
	)
	// False b/c "coreCount" is not a valid string value
	s.FTC("cpuOptions",
		s.A("object",
			s.A(s.A("key", "coreCount"),
				s.A("string", s.A("=", "4")),
			),
		),
	)

	s.TTC("tags",
		s.A("array",
			s.A("some",
				s.A("AND",
					s.A("object",
						s.A(s.A("key", "key"),
							s.A("string", s.A("=", "termination_date"))),
					),
					s.A("object",
						s.A(s.A("key", "value"),
							s.A("time", s.A("<", "2017-08-07T13:55:25.680464+00:00"))),
					),
				),
			),
		),
	)

	/*
		What's the difference between these two? Fuck, maybe the custom
		negation shit I did in find was the way to go. At least it was
		easy to think about. You just reduced stuff and it worked, right huh?

		Domains:
			* First one -- array with @ least one element that's obj w/ "key" and "value" keys
			* Second one -- array with @ least one element that's obj w/ "key" OR "value" keys

		If they are semantically equal, shouldn't they have the same domains? Am I wrong on this?
		I guess it all depends on how I reason about this stuff. Maybe the ValueInDomain shit
		was extremely stupid.

		Maybe it's simpler to just do a non-reduced expression, i.e. what the user meant?

		Is Domain even necessary? 'kind' primary lets me specify what I'm filtering, why do I
		need something else?
	*/
	s.TTC("tags",
		s.A("NOT",
			s.A("array",
				s.A("some",
					s.A("AND",
						s.A("object",
							s.A(s.A("key", "key"),
								s.A("string", s.A("=", "termination_date"))),
						),
						s.A("object",
							s.A(s.A("key", "value"),
								s.A("time", s.A("<", "2017-08-07T13:55:25.680464+00:00"))),
						),
					),
				),
			),
		),
	)
	s.TTC("tags",
		s.A("array",
			s.A("all",
				s.A("NOT",
					s.A("AND",
						s.A("object",
							s.A(s.A("key", "key"),
								s.A("string", s.A("=", "termination_date"))),
						),
						s.A("object",
							s.A(s.A("key", "value"),
								s.A("time", s.A("<", "2017-08-07T13:55:25.680464+00:00"))),
						),
					),
				),
			),
		),
	)

	// False b/c only some of the tags fit the schema
	s.FTC("tags",
		s.A("array",
			s.A("all",
				s.A("AND",
					s.A("object",
						s.A(s.A("key", "key"),
							s.A("string", s.A("=", "termination_date"))),
					),
					s.A("object",
						s.A(s.A("key", "value"),
							s.A("time", s.A("<", "2017-08-07T13:55:25.680464+00:00"))),
					),
				),
			),
		),
	)

	s.TTC("tags",
		s.A("array",
			s.A("some",
				s.A("OR",
					s.A("object",
						s.A(s.A("key", "key"),
							s.A("string", s.A("=", "foo"))),
					),
					s.A("object",
						s.A(s.A("key", "key"),
							s.A("string", s.A("=", "department"))),
					),
				),
			),
		),
	)
	// False b/c "key" cannot be both "foo" and "department"
	s.FTC("tags",
		s.A("array",
			s.A("some",
				s.A("AND",
					s.A("object",
						s.A(s.A("key", "key"),
							s.A("string", s.A("=", "foo"))),
					),
					s.A("object",
						s.A(s.A("key", "key"),
							s.A("string", s.A("=", "department"))),
					),
				),
			),
		),
	)

	// We've tested enough false cases (false expression value, mis-typed values, etc.)
	// that we can now focus on some more "true" cases.

	s.TTC("elasticGpuAssociations", nil)

	s.TTC("networkInterfaces",
		s.A("array",
			s.A("some",
				s.A("AND",
					s.A("object",
						s.A(s.A("key", "association"),
							s.A("object",
								s.A(s.A("key", "ipOwnerID"),
									s.A("string", s.A("=", "amazon"))),
							),
						),
					),
					s.A("object",
						s.A(s.A("key", "privateIpAddresses"),
							s.A("array",
								s.A("some",
									s.A("object",
										s.A(s.A("key", "association"),
											s.A("object",
												s.A(s.A("key", "ipOwnerID"),
													s.A("string", s.A("=", "amazon"))),
											),
										),
									),
								),
							),
						),
					),
				),
			),
		),
	)

	// Test a value PE that combines primitive values
	s.TTC("tags",
		s.A("array",
			s.A("some",
				s.A("object",
					s.A(s.A("key", "key"),
						s.A("OR",
							s.A("string", s.A("=", "foo")),
							s.A("string", s.A("=", "department")),
						)),
				),
			),
		),
	)
}

func (s *MetaPrimaryRealWorldTestSuite) TestMetaPrimary_NegatedValuePredicates() {
	s.TTC("architecture", s.A("NOT", s.A("string", s.A("=", "x86_6"))))

	s.TTC("blockDeviceMappings",
		s.A("array",
			s.A("some",
				s.A("object",
					s.A(s.A("key", "deviceName"),
						s.A("NOT", s.A("string", s.A("=", "/dev/sda")))),
				),
			),
		),
	)
}

func TestMetaPrimaryRealWorld(t *testing.T) {
	s := new(MetaPrimaryRealWorldTestSuite)

	rawMeta, err := ioutil.ReadFile("testdata/metadata.json")
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to read testdata/metadata.json"))
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMeta, &m); err != nil {
		t.Fatal(fmt.Sprintf("Failed to unmarshal testdata/metadata.json: %v", err))
	}
	s.e.Metadata = m

	suite.Run(t, s)
}
