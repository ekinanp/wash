package meta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/puppetlabs/wash/plugin"
	"github.com/stretchr/testify/suite"
)

type SchemaTestSuite struct {
	suite.Suite
}

// TODO: Add a test-case for map[string]interface{} here
func (suite *SchemaTestSuite) TestNewSchema() {
	var schema *plugin.JSONSchema
	// TODO: Once https://github.com/alecthomas/jsonschema/issues/40
	// is (properly) resolved, we should dynamically generate the
	// schema from a struct so maintainers can see what our mock looks
	// like. Right now, the (hacky) fix in our jsonschema fork generates
	// duplicate definitions for anonymous structs (and this behavior's
	// unpredictable), so we store the JSON in a fixture. Note that
	// it still generates the right schema, there's just some redundancy
	// in the generated schema.
	suite.readFixture("before_munging", &schema)
	var expected map[string]interface{}
	suite.readFixture("after_munging", &expected)

	_ = newSchema(schema)

	actualBytes, err := json.Marshal(schema)
	if err != nil {
		suite.FailNow("Failed to marshal the munged JSON schema: %v", err)
	}
	var actual map[string]interface{}
	if err := json.Unmarshal(actualBytes, &actual); err != nil {
		suite.FailNow("Failed to unmarshal the munged JSON schema: %v", err)
	}

	suite.Equal(expected, actual)
}

func (suite *SchemaTestSuite) TestIsValidKeySequenceValidKeySequence() {
	var schema *plugin.JSONSchema
	suite.readFixture("before_munging", &schema)
	var expected map[string]interface{}
	suite.readFixture("after_munging", &expected)

	s := newSchema(schema)
	ks := (keySequence{}).
		EndsWithPrimitiveValue().
		AddObject("dcap").
		AddObject("dcp").
		AddObject("dp")

	suite.True(s.IsValidKeySequence(ks))
}

func (suite *SchemaTestSuite) TestIsValidKeySequenceInvalidKeySequence() {
	var schema *plugin.JSONSchema
	suite.readFixture("before_munging", &schema)
	var expected map[string]interface{}
	suite.readFixture("after_munging", &expected)

	s := newSchema(schema)

	// "DP" is the invalid value here with the invalid property
	// "Foo"
	ks := (keySequence{}).
		EndsWithPrimitiveValue().
		AddObject("foo").
		AddObject("dp")

	suite.False(s.IsValidKeySequence(ks))

	// "AP" is a primitive type, so its value must be "null".
	// Here, however, it is an object.
	ks = (keySequence{}).EndsWithObject().AddObject("ap")
	suite.False(s.IsValidKeySequence(ks))
}

func TestSchema(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}

func (suite *SchemaTestSuite) readFixture(name string, v interface{}) {
	filePath := path.Join("testdata", name+".json")
	rawSchema, err := ioutil.ReadFile(filePath)
	if err != nil {
		suite.T().Fatal(fmt.Sprintf("Failed to read %v", filePath))
	}
	if err := json.Unmarshal(rawSchema, v); err != nil {
		suite.T().Fatal(fmt.Sprintf("Failed to unmarshal %v: %v", filePath, err))
	}
}
