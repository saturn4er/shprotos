package shprotos

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	full "github.com/saturn4er/shprotos/testdata"
	"github.com/stretchr/testify/require"
)

func TestMarshalMessage(t *testing.T) {
	parser := Parser{}
	parsedFile, err := parser.Parse("./testdata/full.proto", nil, nil)
	require.NoError(t, err)

	protoMessage := &full.ComplexMessage{
		Enum:           full.ComplexMessage_VALUE_A,
		ScalarInt32:    -10,
		ScalarInt64:    -20,
		ScalarUint32:   10,
		ScalarUint64:   20,
		ScalarSint32:   -50,
		ScalarSint64:   -500,
		ScalarFixed32:  1000,
		ScalarFixed64:  10000,
		ScalarSfixed32: -100000,
		ScalarSfixed64: -1000000,
		ScalarDouble:   -2.5,
		ScalarFloat:    -3.5,
		ScalarBool:     true,
		ScalarString:   "some another hello world",
		Message: &full.ComplexMessage_SimpleMessage{
			SomeField:  300,
			SomeField2: "hello",
		},
		Bytes: []byte("hello world"),
		MapEnum: map[int32]full.ComplexMessage_SimpleEnum{
			1: full.ComplexMessage_VALUE_A,
			2: full.ComplexMessage_VALUE_B,
			3: full.ComplexMessage_UNSPECIFIED,
		},
		MapScalar: map[int32]int32{
			1: 2,
			2: 3,
			4: 5,
		},
		MapMsg: map[string]*full.ComplexMessage_SimpleMessage{
			"hello": {
				SomeField:  1,
				SomeField2: "hello",
			},
			"world": {
				SomeField:  2,
				SomeField2: "world",
			},
		},
		MapBytes:  nil,
		MapString: nil,
		REnum: []full.ComplexMessage_SimpleEnum{
			full.ComplexMessage_VALUE_B7,
			full.ComplexMessage_VALUE_B8,
		},
		RScalar: []int32{1, 2, 3, 4},
		RMsg: []*full.ComplexMessage_SimpleMessage{
			{
				SomeField:  1,
				SomeField2: "hello",
			},
			{
				SomeField:  2,
				SomeField2: "world",
			},
		},
		RBytes: [][]byte{
			[]byte("hello"),
			[]byte("world"),
		},
		Oneof: &full.ComplexMessage_OneofScalar{OneofScalar: 100},
	}

	msgDesc, ok := parsedFile.Message(TypeName{"ComplexMessage"})
	require.True(t, ok)

	msg := `{
	"bytes": "aGVsbG8gd29ybGQ=",
	"enum": 1,
	"map_enum": {
		"1": 1,
		"2": 2,
		"3": 0
	},
	"map_msg": {
		"hello": {
			"some_field": 1,
			"some_field2": "hello"
		},
		"world": {
			"some_field": 2,
			"some_field2": "world"
		}
	},
	"map_scalar": {
		"1": 2,
		"2": 3,
		"4": 5
	},
	"message": {
		"some_field": 300,
		"some_field2": "hello"
	},
	"oneof_scalar": 100,
	"r_bytes": [
		"aGVsbG8=",
		"d29ybGQ="
	],
	"r_enum": [
		9,
		1000
	],
	"r_msg": [{
			"some_field": 1,
			"some_field2": "hello"
		},
		{
			"some_field": 2,
			"some_field2": "world"
		}
	],
	"r_scalar": [
		1,
		2,
		3,
		4
	],
	"scalar_bool": true,
	"scalar_double": -2.5,
	"scalar_fixed32": 1000,
	"scalar_fixed64": 10000,
	"scalar_float": -3.5,
	"scalar_int32": -10,
	"scalar_int64": -20,
	"scalar_sfixed32": -100000,
	"scalar_sfixed64": -1000000,
	"scalar_sint32": -50,
	"scalar_sint64": -500,
	"scalar_string": "some another hello world",
	"scalar_uint32": 10,
	"scalar_uint64": 20
}`
	var mapa = make(map[string]interface{})
	err = json.Unmarshal([]byte(msg), &mapa)
	require.NoError(t, err)

	res, err := MarshalMessage(mapa, msgDesc)
	require.NoError(t, err)

	spew.Dump(UnmarshalMessage(res, msgDesc))

	resultMsg := &full.ComplexMessage{}
	err = proto.Unmarshal(res, resultMsg)
	require.NoError(t, err)

	require.Equal(t, protoMessage, resultMsg)

}
