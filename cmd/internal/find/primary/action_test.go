package primary

import (
	"testing"

	apitypes "github.com/puppetlabs/wash/api/types"
	"github.com/puppetlabs/wash/cmd/internal/find/types"
	"github.com/stretchr/testify/suite"
)

type ActionPrimaryTestSuite struct {
	primaryTestSuite
}

func (s *ActionPrimaryTestSuite) TestErrors() {
	s.RETC("", "requires additional arguments")
	s.RETC("foo", "foo is an invalid action. Valid actions are.*list")
}

func (s *ActionPrimaryTestSuite) TestValidInput_EntryP() {
	s.RTC("list", "", []string{"list"}, []string{"exec"})
	// Test multiple supported actions
	s.RTC("list", "", []string{"read", "stream", "list"}, []string{"read", "stream"})
}

func (s *ActionPrimaryTestSuite) TestValidInput_SchemaP() {
	// Same test cases as EntryP
	s.RSTC("list", "", []string{"list"}, []string{"exec"})
	s.RSTC("list", "", []string{"read", "stream", "list"}, []string{"read", "stream"})
}

func TestActionPrimary(t *testing.T) {
	s := new(ActionPrimaryTestSuite)
	s.Parser = Action
	s.ConstructEntry = func(v interface{}) types.Entry {
		e := types.Entry{}
		e.Actions = v.([]string)
		return e
	}
	s.ConstructEntrySchema = func(v interface{}) *types.EntrySchema {
		s := &types.EntrySchema{
			EntrySchema: &apitypes.EntrySchema{},
		}
		s.Actions = v.([]string)
		return s
	}
	suite.Run(t, s)
}
