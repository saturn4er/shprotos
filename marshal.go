package shprotos

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"

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
	return unmarshalMessageBytesToMap(proto.NewBuffer(data), msg)
}
func unmarshalMessageBytesToMap(buffer *proto.Buffer, msg *Message) (map[string]interface{}, error) {
	result := make(map[string]interface{})
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
			switch typ := messageField.GetType().(type) {
			case *Enum:
				value, err := buffer.DecodeVarint()
				if err != nil {
					return nil, errors.Wrap(err, "failed to decode varint")
				}
				result[messageField.GetName()] = value
			case *Scalar:
				if typ.ScalarName == "bytes" {
					return nil, errors.New("can't assign varint to bytes field")
				}

				res, err := unmarshaScalar(buffer, typ)
				if err != nil {
					return nil, errors.Wrap(err, "failed to unmarshal varint scalar")
				}
				result[messageField.GetName()] = res
			}
		case WireType64Bit:
			res, err := unmarshaScalar(buffer, messageField.GetType().(*Scalar))
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal varint scalar")
			}
			result[messageField.GetName()] = res
		case WireTypeLengthDelimited:
			data, err := buffer.DecodeRawBytes(false)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode raw bytes")
			}

			switch typ := messageField.GetType().(type) {
			case *Map:
				if _, ok := result[messageField.GetName()]; !ok {
					result[messageField.GetName()] = make(map[interface{}]interface{})
				}
				resultMap := result[messageField.GetName()].(map[interface{}]interface{})

				mapBuffer := proto.NewBuffer(data)
				mapKeyKey, err := mapBuffer.DecodeVarint()
				if err != nil {
					return nil, errors.Wrap(err, "failed to decode map key key")
				}
				mapKeyWireType := mapKeyKey & 7

				var mapKey interface{}
				var mapValue interface{}
				switch mapKeyWireType {
				case WireTypeVarint, WireType64Bit, WireType32Bit:
					res, err := unmarshaScalar(mapBuffer, typ.KeyType.(*Scalar))
					if err != nil {
						return nil, errors.Wrap(err, "failed to unmarshal varint scalar")
					}
					mapKey = res
				case WireTypeLengthDelimited:
					str, err := mapBuffer.DecodeStringBytes()
					if err != nil {
						return nil, errors.Wrap(err, "failed to decode string bytes")
					}
					mapKey = str
				}

				mapValueKey, err := mapBuffer.DecodeVarint()
				if err != nil {
					return nil, errors.Wrap(err, "failed to decode map value key")
				}
				mapValueWireType := mapValueKey & 7
				switch mapValueWireType {
				case WireTypeVarint, WireType32Bit, WireType64Bit:
					res, err := unmarshaScalar(mapBuffer, typ.KeyType.(*Scalar))
					if err != nil {
						return nil, errors.Wrap(err, "failed to unmarshal varint scalar")
					}
					mapValue = res
				case WireTypeLengthDelimited:
					valueBytes, err := mapBuffer.DecodeRawBytes(false)
					if err != nil {
						return nil, errors.Wrap(err, "failed to decode map value raw bytes")
					}
					switch valueType := typ.ValueType.(type) {
					case *Message:
						msgValue, err := unmarshalMessageBytesToMap(proto.NewBuffer(valueBytes), valueType)
						if err != nil {
							return nil, errors.WithStack(err)
						}
						mapValue = msgValue
					}
				}
				resultMap[mapKey] = mapValue
			case *Message:
				msgValue, err := UnmarshalMessageBytesToMap(data, typ)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				if messageField.IsRepeated() {
					if _, ok := result[messageField.GetName()]; !ok {
						result[messageField.GetName()] = []interface{}{}
					}
					result[messageField.GetName()] = append(result[messageField.GetName()].([]interface{}), msgValue)
				} else {
					result[messageField.GetName()] = msgValue
				}
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
							elem, err := unmarshaScalar(buff, typ)
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
				var values []int32
				for {
					elem, err := buff.DecodeVarint()
					if err != nil {
						if err == io.ErrUnexpectedEOF {
							break
						}
						return nil, errors.Wrap(err, "failed to get varint")
					}
					values = append(values, int32(elem))
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
			res, err := unmarshaScalar(buffer, messageField.GetType().(*Scalar))
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal fixed 32 scalar")

			}
			result[messageField.GetName()] = res
		}
	}
	return result, nil
}

func unmarshaScalar(buffer *proto.Buffer, scalar *Scalar) (interface{}, error) {
	switch scalar.ScalarName {
	case "sfixed32":
		value, err := buffer.DecodeFixed32()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode zigzag32")
		}
		return int32(value), nil
	case "sfixed64":
		value, err := buffer.DecodeFixed64()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode zigzag64")
		}
		return int64(value), nil
	case "sint32":
		value, err := buffer.DecodeZigzag32()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode zigzag32")
		}
		return int32(value), nil
	case "sint64":
		value, err := buffer.DecodeZigzag64()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode zigzag64")
		}
		return int64(value), nil
	case "fixed32":
		value, err := buffer.DecodeFixed32()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode fixed 32")
		}
		return uint32(value), nil
	case "fixed64":
		value, err := buffer.DecodeFixed64()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode fixed 64")
		}
		return uint64(value), nil
	case "int32", "int64", "uint32", "uint64", "bool", "double", "float":
		value, err := buffer.DecodeVarint()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "failed to decode varint")
		}
		switch scalar.ScalarName {
		case "int32":
			return int32(value), nil
		case "int64":
			return int64(value), nil
		case "uint32":
			return uint32(value), nil
		case "uint64":
			return uint64(value), nil
		case "bool":
			return value == 1, nil
		case "double":
			return math.Float64frombits(value), nil
		case "float":
			return math.Float32frombits(uint32(value)), nil
		}
		return int32(value), nil
	}
	return nil, errors.Errorf("unknown scalar type: %s", scalar.ScalarName)
}
