package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/ekinanp/jsonschema"
	"github.com/stretchr/testify/suite"
)

type MetadataSchemaTestSuite struct {
	suite.Suite
}

type Schema struct {
	A int
	B int
	C []int
	D struct {
		D_A string
		D_B bool
		D_C struct {
			D_C_A int
			D_C_B int
		}
	}
	E []struct {
		E_A int
		E_B int
	}
}

func (suite *MetadataSchemaTestSuite) TestNewMetadataSchema() {
	r := jsonschema.Reflector{
		ExpandedStruct: true,
	}
	schema := r.Reflect(mockSchema{})
	_ = NewMetadataSchema(schema)
	expected := suite.readFixture("after_munging")

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

func TestMetadataSchema(t *testing.T) {
	suite.Run(t, new(MetadataSchemaTestSuite))
}

func (suite *MetadataSchemaTestSuite) readFixture(name string) map[string]interface{} {
	filePath := path.Join("testdata", "metadataSchema", name+".json")
	rawSchema, err := ioutil.ReadFile(filePath)
	if err != nil {
		suite.T().Fatal(fmt.Sprintf("Failed to read %v", filePath))
	}
	var mp map[string]interface{}
	if err := json.Unmarshal(rawSchema, &mp); err != nil {
		suite.T().Fatal(fmt.Sprintf("Failed to unmarshal %v: %v", filePath, err))
	}
	return mp
}

type mockSchema struct {
	Ap int
	Bp int
	Cp []int
	Dp struct {
		DAp string
		DBp bool
		DCp struct {
			DCAp int
			DCBp int
		}
	}
	Ep []struct {
		EAp int
		EBp int
	}
}
