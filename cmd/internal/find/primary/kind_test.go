package primary

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/puppetlabs/wash/cmd/internal/find/types"
	"github.com/stretchr/testify/suite"
)

type KindPrimaryTestSuite struct {
	primaryTestSuite
}

func (s *KindPrimaryTestSuite) TestErrors() {
	s.RETC("", "requires additional arguments")
	s.RETC("[a", "invalid glob: unexpected end of input")
}

func (s *KindPrimaryTestSuite) TestValidInput() {
	// Test the entry predicate
	s.RTC("a", "", types.Entry{})
	// Test the main schema predicate
	s.RSTC("dock*container", "", "docker/containers/container", "docker/containers/container/volume")
}

func (s *KindPrimaryTestSuite) TestKindP() {
	g, err := glob.Compile("dock*container")
	if s.NoError(err) {
		p := kindP(g, false)

		// Test the entry predicate
		s.True(p.P(types.Entry{}))

		// Test the schema predicate
		schema := &types.EntrySchema{}
		schema.SetPathsToNode([]string{"docker/containers/container"})
		s.True(p.SchemaP().P(schema))
		schema.SetPathsToNode([]string{
			"docker/containers/container/fs/dir/file",
			"docker/containers/container/fs/file",
			"docker/volumes/volume/dir/file",
			"docker/volumes/volume/file",
		})
		s.False(p.SchemaP().P(schema))

		// Ensure that the predicate requires entry schemas
		s.True(p.SchemaRequired())
	}
}

func (s *KindPrimaryTestSuite) TestKindP_Negate() {
	g, err := glob.Compile("dock*container")
	if s.NoError(err) {
		p := kindP(g, false).Negate().(types.EntryPredicate)

		// Test the entry predicate
		s.True(p.P(types.Entry{}))

		// Test the schema predicate
		schema := &types.EntrySchema{}
		schema.SetPathsToNode([]string{"docker/containers/container"})
		s.False(p.SchemaP().P(schema))
		schema.SetPathsToNode([]string{
			"docker/containers/container/fs/dir/file",
			"docker/containers/container/fs/file",
			"docker/volumes/volume/dir/file",
			"docker/volumes/volume/file",
		})
		s.True(p.SchemaP().P(schema))

		// Ensure that the predicate still requires entry schemas
		s.True(p.SchemaRequired())
	}
}

func TestKindPrimary(t *testing.T) {
	s := new(KindPrimaryTestSuite)
	s.Parser = Kind
	s.SchemaPParser = types.EntryPredicateParser(Kind.parseFunc).ToSchemaPParser()
	s.ConstructEntry = func(v interface{}) types.Entry {
		return v.(types.Entry)
	}
	s.ConstructEntrySchema = func(v interface{}) *types.EntrySchema {
		s := &types.EntrySchema{}
		switch t := v.(type) {
		case []string:
			s.SetPathsToNode(t)
		case string:
			s.SetPathsToNode([]string{t})
		default:
			panic("The kind primary's tests must take an array of paths (strings) as satisfying values")
		}
		return s
	}
	suite.Run(t, s)
}
