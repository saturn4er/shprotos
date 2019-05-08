package parser

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	WireTypeVarint          = 0
	WireType64Bit           = 1
	WireTypeLengthDelimited = 2
	WireTypeStartGroup      = 3 // deprecated
	WireTypeEndGroup        = 4 // deprecated
	WireType32Bit           = 5
)

func UnmarshalMessageBytesToMap(data []byte, msg *Message) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	buffer := proto.NewBuffer(data)
	for {
		key, err := buffer.DecodeVarint()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return result, nil
			}
			return nil, errors.Wrap(err, "failed to get key")
		}
		fieldNum := key >> 3
		messageField, ok := msg.FieldByKeyNumber(int(fieldNum))
		if !ok {
			return nil, errors.Errorf("can't find field %d in message %s", fieldNum, msg.Name)
		}
		switch key & 7 {
		case WireTypeVarint:
			value, err := buffer.DecodeVarint()
			if err != nil {
				panic(err)
			}
			switch typ := messageField.GetType().(type) {
			case *Enum:
				result[messageField.GetName()] = value
			case *Scalar:
				if typ.ScalarName == "bytes" {
					return nil, errors.New("can't assign varint to bytes field")
				}
				result[messageField.GetName()] = value
			}
		case WireType64Bit:
			value, err := buffer.DecodeFixed64()
			if err != nil {
				panic(err)
			}
			fmt.Printf("field %d is equal to %v\n", fieldNum, value)
		case WireTypeLengthDelimited:
			data, err := buffer.DecodeRawBytes(false)
			if err != nil {
				panic(err)
			}

			switch typ := messageField.GetType().(type) {
			case *Map:
			case *Message:
				msgValue, err := UnmarshalMessageBytesToMap(data, typ)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				result[messageField.GetName()] = msgValue
			case *Scalar:
				switch typ.ScalarName {
				case "bytes":
					val := base64.StdEncoding.EncodeToString(data)
					if messageField.IsRepeated() {
						if _, ok := result[messageField.GetName()]; !ok {
							result[messageField.GetName()] = []interface{}{}
						}
						result[messageField.GetName()] = append(result[messageField.GetName()].([]interface{}), string(val))
					} else {
						result[messageField.GetName()] = string(val)
					}
				default:
					if messageField.IsRepeated() {
						buff := proto.NewBuffer(data)
						var values []interface{}
						for {
							elem, err := buff.DecodeVarint()
							if err != nil {
								if err == io.ErrUnexpectedEOF {
									break
								}
								return nil, errors.Wrap(err, "failed to get varint")
							}
							values = append(values, elem)
						}
						result[messageField.GetName()] = values
					} else {
						result[messageField.GetName()] = string(data)
					}
				}

			case *Enum:
				buff := proto.NewBuffer(data)
				var values []uint64
				for {
					elem, err := buff.DecodeVarint()
					if err != nil {
						if err == io.ErrUnexpectedEOF {
							break
						}
						return nil, errors.Wrap(err, "failed to get varint")
					}
					values = append(values, elem)
				}
				result[messageField.GetName()] = values
			default:
				fmt.Printf("unknown length delimited value%T\n", messageField.GetType())
			}
		case WireTypeStartGroup:
			fmt.Println("Unhandled group start")
		case WireTypeEndGroup:
			fmt.Println("Unhandled group end")
		case WireType32Bit:
			fmt.Println("Unhandled 32 bit")
		}
	}
	return result, nil
}
