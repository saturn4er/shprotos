package shprotos

import (
	"encoding/base64"
	"math"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func MarshalMessage(data map[string]interface{}, message *Message) ([]byte, error) {
	res := proto.NewBuffer(nil)
	err := marshalMessage(res, data, message)
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func marshalMessage(buffer *proto.Buffer, data map[string]interface{}, message *Message) error {
	for _, field := range message.GetFields() {
		fieldValue, ok := data[field.GetName()]
		if !ok {
			continue
		}
		switch fld := field.(type) {
		case *NormalField:
			if fld.IsRepeated() {
				if err := marshalMessageNormalRepeatedField(buffer, fieldValue, fld.Type, fld.KeyNumber); err != nil {
					return errors.Wrapf(err, "failed to marshal normal repeated field %s", field.GetName())
				}
			} else {
				if err := marshalMessageNormalField(buffer, fieldValue, fld.Type, fld.KeyNumber); err != nil {
					return errors.Wrapf(err, "failed to marshal normal field %s", field.GetName())
				}
			}
		case *MapField:
			mapValue := reflect.ValueOf(fieldValue)
			iter := mapValue.MapRange()
			for iter.Next() {
				mapBuffer := proto.NewBuffer(nil)
				mapKey := iter.Key().Interface()
				mapValue := iter.Value().Interface()
				if err := buffer.EncodeVarint(messageKeyVarint(fld.KeyNumber, WireTypeLengthDelimited)); err != nil {
					return errors.Wrap(err, "failed to write field key")
				}
				if err := marshalMessageNormalField(mapBuffer, mapKey, fld.Map.KeyType, 1); err != nil {
					return errors.Wrap(err, "failed to marshal map key")
				}
				if err := marshalMessageNormalField(mapBuffer, mapValue, fld.Map.ValueType, 2); err != nil {
					return errors.Wrap(err, "failed to marshal map key")
				}
				if err := buffer.EncodeRawBytes(mapBuffer.Bytes()); err != nil {
					return errors.Wrap(err, "failed to encode map buffer bytes")
				}
			}
		}
	}
	return nil
}

func marshalMessageNormalField(buffer *proto.Buffer, value interface{}, typ Type, keyNumber uint64) error {
	switch typ := typ.(type) {
	case *Scalar:
		switch typ.ScalarName {
		case "string":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeLengthDelimited)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			if err := buffer.EncodeStringBytes(value.(string)); err != nil {
				return errors.Wrap(err, "failed to encode string bytes")
			}
		case "bytes":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeLengthDelimited)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			value, err := base64.StdEncoding.DecodeString(value.(string))
			if err != nil {
				return errors.Wrap(err, "failed to encode base64 string")
			}
			if err := buffer.EncodeRawBytes(value); err != nil {
				return errors.Wrap(err, "failed to encode string bytes")
			}
		case "bool":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeVarint)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			res := uint64(0)
			if value.(bool) {
				res = 1
			}
			if err := buffer.EncodeVarint(res); err != nil {
				return errors.Wrap(err, "failed to encode bool")
			}
		case "sfixed32", "fixed32":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireType32Bit)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			val, err := uint64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeFixed32(val); err != nil {
				return errors.Wrap(err, "failed to encode fixed 32")
			}
		case "sfixed64", "fixed64":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireType64Bit)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			val, err := uint64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeFixed64(val); err != nil {
				return errors.Wrap(err, "failed to encode fixed 64")
			}
		case "sint32":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeVarint)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			val, err := uint64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeZigzag32(val); err != nil {
				return errors.Wrap(err, "failed to encode zigzag 32")
			}
		case "sint64":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeVarint)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			val, err := uint64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeZigzag64(val); err != nil {
				return errors.Wrap(err, "failed to encode zigzag 64")
			}
		case "float":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireType32Bit)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			fl32, err := float32FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeFixed32(uint64(math.Float32bits(fl32))); err != nil {
				return errors.Wrap(err, "failed to encode fixed 32")
			}
		case "double":
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireType64Bit)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			fl64, err := float64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeFixed64(math.Float64bits(fl64)); err != nil {
				return errors.Wrap(err, "failed to encode fixed 64")
			}
		default:
			if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeVarint)); err != nil {
				return errors.Wrap(err, "failed to write field key")
			}
			val, err := uint64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to resolve value")
			}
			if err := buffer.EncodeVarint(val); err != nil {
				return errors.Wrap(err, "failed to encode fixed 64")
			}
		}
	case *Enum:
		if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeVarint)); err != nil {
			return errors.Wrap(err, "failed to write field key")
		}
		res, err := uint64FromInterface(value)
		if err != nil {
			return errors.Wrap(err, "failed to resolve enum value")
		}
		if err := buffer.EncodeVarint(res); err != nil {
			return errors.Wrap(err, "failed to encode bool")
		}
	case *Message:
		if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeLengthDelimited)); err != nil {
			return errors.Wrap(err, "failed to write field key")
		}
		msgData, err := MarshalMessage(value.(map[string]interface{}), typ)
		if err != nil {
			return errors.Wrap(err, "failed to marshal message")
		}
		if err := buffer.EncodeRawBytes(msgData); err != nil {
			return errors.Wrap(err, "failed to encode message")
		}

	}
	return nil
}

func marshalMessageNormalRepeatedField(buffer *proto.Buffer, value interface{}, typ Type, keyNumber uint64) error {
	switch typ := typ.(type) {
	case *Scalar, *Message:
		values := value.([]interface{})
		for _, value := range values {
			if err := marshalMessageNormalField(buffer, value, typ, keyNumber); err != nil {
				return errors.Wrap(err, "failed to marshal message normal field")
			}
		}
	case *Enum:
		fieldBuffer := proto.NewBuffer(nil)
		values := value.([]interface{})
		for _, value := range values {
			enumValue, err := uint64FromInterface(value)
			if err != nil {
				return errors.Wrap(err, "failed to get enum value")
			}
			if err := fieldBuffer.EncodeVarint(enumValue); err != nil {
				return errors.Wrap(err, "failed to encode enum array value")
			}
		}
		if err := buffer.EncodeVarint(messageKeyVarint(keyNumber, WireTypeLengthDelimited)); err != nil {
			return errors.Wrap(err, "failed to write field key")
		}
		if err := buffer.EncodeRawBytes(fieldBuffer.Bytes()); err != nil {
			return errors.Wrap(err, "failed to encode enum array field")
		}

	}
	return nil
}

func messageKeyVarint(fieldNum uint64, wireType uint64) uint64 {
	return uint64((fieldNum << 3) | wireType)
}

func uint64FromInterface(val interface{}) (uint64, error) {
	switch val := val.(type) {
	case int:
		return uint64(val), nil
	case uint8:
		return uint64(val), nil
	case uint16:
		return uint64(val), nil
	case uint32:
		return uint64(val), nil
	case uint64:
		return val, nil
	case int8:
		return uint64(val), nil
	case int16:
		return uint64(val), nil
	case int32:
		return uint64(val), nil
	case int64:
		return uint64(val), nil
	case float32:
		return uint64(val), nil
	case float64:
		return uint64(val), nil
	case string:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse int from string")
		}
		return uint64(v), nil
	}
	return 0, errors.Errorf("can't convert %T to uint64", val)
}
func float64FromInterface(val interface{}) (float64, error) {
	switch val := val.(type) {
	case int:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case string:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse float64 from string")
		}
		return v, nil
	}
	return 0, errors.Errorf("can't convert %T to uint64", val)
}
func float32FromInterface(val interface{}) (float32, error) {
	switch val := val.(type) {
	case int:
		return float32(val), nil
	case uint8:
		return float32(val), nil
	case uint16:
		return float32(val), nil
	case uint32:
		return float32(val), nil
	case uint64:
		return float32(val), nil
	case int8:
		return float32(val), nil
	case int16:
		return float32(val), nil
	case int32:
		return float32(val), nil
	case int64:
		return float32(val), nil
	case float32:
		return val, nil
	case float64:
		return float32(val), nil
	case string:
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse float32 from string")
		}
		return float32(v), nil
	}
	return 0, errors.Errorf("can't convert %T to uint64", val)
}
