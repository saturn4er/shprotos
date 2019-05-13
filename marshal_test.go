package parser

import (
	"encoding/base64"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/saturn4er/shprotos/testdata"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalMessageBytesToMap(t *testing.T) {
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
		ScalarBool:     false,
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
	}
	data, err := proto.Marshal(protoMessage)
	require.NoError(t, err)

	msg, ok := parsedFile.Message(TypeName{"ComplexMessage"})
	require.True(t, ok)

	res, err := UnmarshalMessageBytesToMap(data, msg)
	require.NoError(t, err)

	spew.Dump(res)

	require.Equal(t, uint64(protoMessage.Enum), res["enum"])

	require.Equal(t, protoMessage.ScalarInt32, res["scalar_int32"])
	require.Equal(t, protoMessage.ScalarInt64, res["scalar_int64"])
	require.Equal(t, protoMessage.ScalarUint32, res["scalar_uint32"])
	require.Equal(t, protoMessage.ScalarUint64, res["scalar_uint64"])
	require.Equal(t, protoMessage.ScalarSint32, res["scalar_sint32"])
	require.Equal(t, protoMessage.ScalarSint64, res["scalar_sint64"])
	require.Equal(t, protoMessage.ScalarFixed32, res["scalar_fixed32"])
	require.Equal(t, protoMessage.ScalarFixed64, res["scalar_fixed64"])
	require.Equal(t, protoMessage.ScalarSfixed32, res["scalar_sfixed32"])
	require.Equal(t, protoMessage.ScalarSfixed64, res["scalar_sfixed64"])
	require.Equal(t, protoMessage.ScalarBool, res["scalar_bool"])
	require.Equal(t, protoMessage.ScalarString, res["scalar_string"])
	require.Equal(t, map[string]interface{}{
		"some_field":  protoMessage.Message.SomeField,
		"some_field2": protoMessage.Message.SomeField2,
	}, res["message"])
	require.Equal(t, base64.StdEncoding.EncodeToString(protoMessage.Bytes), res["bytes"])

	require.Equal(t, map[interface{}]interface{}{
		int32(1): int32(full.ComplexMessage_VALUE_A),
		int32(2): int32(full.ComplexMessage_VALUE_B),
		int32(3): int32(full.ComplexMessage_UNSPECIFIED),
	}, res["map_enum"])
	require.Equal(t, map[interface{}]interface{}{
		int32(1): int32(2),
		int32(2): int32(3),
		int32(4): int32(5),
	}, res["map_scalar"])
	require.Equal(t, map[interface{}]interface{}{
		"hello": map[string]interface{}{
			"some_field":  int32(1),
			"some_field2": "hello",
		},
		"world": map[string]interface{}{
			"some_field":  int32(2),
			"some_field2": "world",
		},
	}, res["map_msg"])

	require.Equal(t, []int32{
		int32(full.ComplexMessage_VALUE_B7),
		int32(full.ComplexMessage_VALUE_B8),
	}, res["r_enum"])
	require.Equal(t, []interface{}{
		int32(1), int32(2), int32(3), int32(4),
	}, res["r_scalar"])
	require.Equal(t, []interface{}{
		map[string]interface{}{
			"some_field":  int32(1),
			"some_field2": "hello",
		},
		map[string]interface{}{
			"some_field":  int32(2),
			"some_field2": "world",
		},
	}, res["r_msg"])
	require.Equal(t, []interface{}{
		"aGVsbG8=",
		"d29ybGQ=",
	}, res["r_bytes"])
}
