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
		Enum:   full.ComplexMessage_VALUE_A,
		Scalar: 100,
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
		MapMsg: map[int32]*full.ComplexMessage_SimpleMessage{
			1: {
				SomeField:  1,
				SomeField2: "hello",
			},
			2: {
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

	require.Equal(t, uint64(full.ComplexMessage_VALUE_A), res["enum"])
	require.Equal(t, uint64(100), res["scalar"])
	require.Equal(t, map[string]interface{}{
		"some_field":  uint64(300),
		"some_field2": "hello",
	}, res["message"])
	require.Equal(t, base64.StdEncoding.EncodeToString([]byte("hello world")), res["bytes"])
	// require.Equal(t, res["map_enum"], )
	// require.Equal(t, res["map_scalar"], )
	// require.Equal(t, res["map_msg"], )
	require.Equal(t, []interface{}{
		uint64(full.ComplexMessage_VALUE_B7),
		uint64(full.ComplexMessage_VALUE_B8),
	}, res["r_enum"])
	// require.Equal(t, res["r_scalar"], )
	// require.Equal(t, res["r_msg"], )
	// require.Equal(t, res["r_bytes"], )
}
