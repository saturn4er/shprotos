package shprotos

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testFileInfo(file *File) *File {
	var Int32Type = &Scalar{file: file, ScalarName: "int32"}

	var StringType = &Scalar{file: file, ScalarName: "string"}

	var RootMessage = file.Messages[0]

	var RootMessage2 = file.Messages[4]

	var RootEnum = file.Enums[0]

	var EmptyMessage = file.Messages[2]

	var NestedMessage = file.Messages[1]

	var NestedEnum = file.Enums[1]

	var NestedNestedEnum = file.Enums[2]

	var MessageWithEmpty = file.Messages[3]

	var CommonCommonEnum = file.Imports[0].Enums[0]

	var CommonCommonMessage = file.Imports[0].Messages[0]

	var ParentScopeEnum = file.Imports[2].Enums[0]

	var Proto2Message = file.Imports[1].Messages[0]

	return &File{
		Services: []*Service{
			{
				Name:          "ServiceExample",
				QuotedComment: `"Service, which do smth"`,
				Methods: []*Method{
					{Name: "getQueryMethod", InputMessage: RootMessage, OutputMessage: RootMessage, QuotedComment: `""`},
					{Name: "mutationMethod", InputMessage: RootMessage2, OutputMessage: NestedMessage, QuotedComment: `"rpc comment"`},
					{Name: "EmptyMsgs", InputMessage: EmptyMessage, OutputMessage: EmptyMessage, QuotedComment: `""`},
					{Name: "MsgsWithEpmty", InputMessage: MessageWithEmpty, OutputMessage: MessageWithEmpty, QuotedComment: `""`},
				},
			},
		},
		Messages: []*Message{
			{
				file:          file,
				Name:          "RootMessage",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage"},
				NormalFields: []*NormalField{
					{Name: "r_msg", Type: NestedMessage, Repeated: true, QuotedComment: `"repeated Message"`},
					{Name: "r_scalar", Type: Int32Type, Repeated: true, QuotedComment: `"repeated Scalar"`},
					{Name: "r_enum", Type: RootEnum, Repeated: true, QuotedComment: `"repeated Enum"`},
					{Name: "r_empty_msg", Type: EmptyMessage, Repeated: true, QuotedComment: `"repeated empty message"`},
					{Name: "n_r_enum", Type: CommonCommonEnum, QuotedComment: `"non-repeated Enum"`},
					{Name: "n_r_scalar", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
					{Name: "n_r_msg", Type: CommonCommonMessage, QuotedComment: `"non-repeated Message"`},
					{Name: "scalar_from_context", Type: Int32Type, QuotedComment: `"field from context"`},
					{Name: "n_r_empty_msg", Type: EmptyMessage, QuotedComment: `"non-repeated empty message field"`},
					{Name: "leading_dot", Type: CommonCommonMessage, QuotedComment: `"leading dot in type name"`},
					{Name: "parent_scope", Type: ParentScopeEnum, QuotedComment: `"parent scope"`},
					{Name: "proto2message", Type: Proto2Message, QuotedComment: `""`},
				},
				OneOffs: []*OneOf{
					{Name: "enum_first_oneoff", Fields: []*NormalField{
						{Name: "e_f_o_e", Type: CommonCommonEnum, QuotedComment: `""`},
						{Name: "e_f_o_s", Type: Int32Type, QuotedComment: `""`},
						{Name: "e_f_o_m", Type: CommonCommonMessage, QuotedComment: `""`},
						{Name: "e_f_o_em", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
					}},
					{Name: "scalar_first_oneoff", Fields: []*NormalField{
						{Name: "s_f_o_s", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
						{Name: "s_f_o_e", Type: RootEnum, QuotedComment: `"non-repeated Enum"`},
						{Name: "s_f_o_mes", Type: RootMessage2, QuotedComment: `"non-repeated Message"`},
						{Name: "s_f_o_m", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
					}},
					{Name: "message_first_oneoff", Fields: []*NormalField{
						{Name: "m_f_o_m", Type: RootMessage2, QuotedComment: `"non-repeated Message"`},
						{Name: "m_f_o_s", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
						{Name: "m_f_o_e", Type: RootEnum, QuotedComment: `"non-repeated Enum"`},
						{Name: "m_f_o_em", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
					}},
					{Name: "empty_first_oneoff", Fields: []*NormalField{
						{Name: "em_f_o_em", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
						{Name: "em_f_o_s", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
						{Name: "em_f_o_en", Type: RootEnum, QuotedComment: `"non-repeated Enum"`},
						{Name: "em_f_o_m", Type: RootMessage2, QuotedComment: `"non-repeated Message"`},
					}},
				},
				MapFields: []*MapField{
					{
						Name:          "map_enum",
						QuotedComment: `"enum_map\n Map with enum value"`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   Int32Type,
							ValueType: NestedEnum,
							file:      file,
						},
					},
					{
						Name:          "map_scalar",
						QuotedComment: `"scalar map\n Map with scalar value"`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   Int32Type,
							ValueType: Int32Type,
							file:      file,
						},
					},
					{
						Name:          "map_msg",
						QuotedComment: `"Map with Message value"`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   StringType,
							ValueType: NestedMessage,
							file:      file,
						},
					},
					{
						Name:          "ctx_map",
						QuotedComment: `""`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   StringType,
							ValueType: NestedMessage,
							file:      file,
						},
					},
					{
						Name:          "ctx_map_enum",
						QuotedComment: `""`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   StringType,
							ValueType: NestedEnum,
							file:      file,
						},
					},
				},
			},
			{
				file:          file,
				Name:          "NestedMessage",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage", "NestedMessage"},
				NormalFields: []*NormalField{
					{Name: "sub_r_enum", Type: NestedEnum, Repeated: true, QuotedComment: `"repeated Enum"`},
					{Name: "sub_sub_r_enum", Type: NestedNestedEnum, Repeated: true, QuotedComment: `"repeated Enum"`},
				},
			},
			{
				file:          file,
				Name:          "Empty",
				QuotedComment: `""`,
				TypeName:      TypeName{"Empty"},
			},
			{
				file:          file,
				Name:          "MessageWithEmpty",
				QuotedComment: `""`,
				TypeName:      TypeName{"MessageWithEmpty"},
				NormalFields: []*NormalField{
					{Name: "empt", Type: EmptyMessage, QuotedComment: `""`},
				},
			},
			{
				file:          file,
				Name:          "RootMessage2",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage2"},
				NormalFields: []*NormalField{
					{Name: "some_field", Type: Int32Type, QuotedComment: `""`},
				},
			},
		},
		Enums: []*Enum{
			{
				file:          file,
				Name:          "RootEnum",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootEnum"},
				Values: []*EnumValue{
					{Name: "RootEnumVal0", Value: 0, QuotedComment: `""`},
					{Name: "RootEnumVal1", Value: 1, QuotedComment: `""`},
					{Name: "RootEnumVal2", Value: 2, QuotedComment: `"It's a RootEnumVal2"`},
				},
			},
			{
				file:          file,
				Name:          "NestedEnum",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage", "NestedEnum"},
				Values: []*EnumValue{
					{Name: "NestedEnumVal0", Value: 0, QuotedComment: `""`},
					{Name: "NestedEnumVal1", Value: 1, QuotedComment: `""`},
				},
			},
			{
				file:          file,
				Name:          "NestedNestedEnum",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage", "NestedMessage", "NestedNestedEnum"},
				Values: []*EnumValue{
					{Name: "NestedNestedEnumVal0", Value: 0, QuotedComment: `""`},
					{Name: "NestedNestedEnumVal1", Value: 1, QuotedComment: `""`},
					{Name: "NestedNestedEnumVal2", Value: 2, QuotedComment: `""`},
					{Name: "NestedNestedEnumVal3", Value: 3, QuotedComment: `""`},
				},
			},
		},
	}
}

func TestParser_Parse(t *testing.T) {
	Convey("Test Parser.Parse", t, func(c C) {
		parser := Parser{}
		test, err := parser.Parse("../../../../testdata/test.proto", nil, []string{"../../../../testdata"})
		c.So(err, ShouldBeNil)
		c.So(test, ShouldNotBeNil)
		test2, err := parser.Parse("../../../../testdata/test2.proto", nil, []string{"../../../../testdata"})
		c.So(err, ShouldBeNil)
		c.So(test2, ShouldNotBeNil)
		c.So(test, ShouldNotEqual, test2)

		c.Convey("Imports should be the same", func(c C) {
			c.So(len(test.Imports), ShouldEqual, 3)
			c.So(len(test2.Imports), ShouldEqual, 1)
			c.So(test.Imports[0], ShouldEqual, test2.Imports[0])
		})
		c.Convey("If we trying to parse same File, it should return pointer to parsed one", func(c C) {
			test22, err := parser.Parse("../../../../testdata/test2.proto", nil, []string{"../../../../testdata"})
			c.So(err, ShouldBeNil)
			c.So(test22, ShouldEqual, test2)
		})
		f := testFileInfo(test)

		c.Convey("test.proto Should contains valid enums", func(c C) {
			c.So(test.Enums, ShouldHaveLength, len(f.Enums))
			for i, enum := range test.Enums {
				validEnum := f.Enums[i]
				c.Convey("Should contain "+validEnum.Name, func(c C) {
					c.So(enum.File, ShouldEqual, validEnum.File)
					c.So(enum.Name, ShouldEqual, validEnum.Name)
					c.So(enum, ShouldEqual, enum)
					c.So(enum.File(), ShouldEqual, test)
					c.So(enum.TypeName, ShouldResemble, validEnum.TypeName)
					c.So(enum.QuotedComment, ShouldEqual, validEnum.QuotedComment)
					c.Convey(validEnum.Name+" enum should contains valid values", func(c C) {
						c.So(enum.Values, ShouldHaveLength, len(validEnum.Values))
						for i, value := range enum.Values {
							validValue := validEnum.Values[i]
							c.Convey(validEnum.Name+" enum should contains valid "+validValue.Name+" value", func(c C) {
								c.So(value.Name, ShouldEqual, validValue.Name)
								c.So(value.Value, ShouldEqual, validValue.Value)
								c.So(value.QuotedComment, ShouldEqual, validValue.QuotedComment)
							})
						}
					})
				})
			}
		})

		c.Convey("test.proto Should contains valid messages", func(c C) {
			c.So(test.Messages, ShouldHaveLength, len(f.Messages))
			for i, msg := range test.Messages {
				validMsg := f.Messages[i]
				c.Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+" message ", func(c C) {
					c.So(msg.File, ShouldEqual, validMsg.File)
					c.So(msg.Name, ShouldEqual, validMsg.Name)
					c.So(msg, ShouldEqual, msg)
					c.So(msg.File(), ShouldEqual, test)
					c.So(msg.TypeName, ShouldResemble, validMsg.TypeName)
					c.So(msg.QuotedComment, ShouldEqual, validMsg.QuotedComment)
					c.So(msg.NormalFields, ShouldHaveLength, len(validMsg.NormalFields))
					for i, fld := range msg.NormalFields {
						validFld := validMsg.NormalFields[i]
						c.Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validFld.Name+" field", func(c C) {
							c.So(fld.Name, ShouldEqual, validFld.Name)
							c.So(fld.Repeated, ShouldEqual, validFld.Repeated)
							c.So(fld.QuotedComment, ShouldEqual, validFld.QuotedComment)
							CompareTypes(c, fld.Type, validFld.Type)
						})
					}
					c.So(msg.MapFields, ShouldHaveLength, len(validMsg.MapFields))
					for i, fld := range msg.MapFields {
						validFld := validMsg.MapFields[i]
						c.Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validFld.Name+" field", func(c C) {
							c.So(fld.Name, ShouldEqual, validFld.Name)
							c.So(fld.QuotedComment, ShouldEqual, validFld.QuotedComment)
							CompareTypes(c, fld.Map, validFld.Map)
						})
					}
					c.So(msg.OneOffs, ShouldHaveLength, len(validMsg.OneOffs))
					for i, oneOf := range msg.OneOffs {
						validOneOf := validMsg.OneOffs[i]
						c.Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validOneOf.Name+" one of", func(c C) {
							c.So(oneOf.Name, ShouldEqual, validOneOf.Name)
							c.So(oneOf.Fields, ShouldHaveLength, len(validOneOf.Fields))
							for i, fld := range oneOf.Fields {
								validFld := validOneOf.Fields[i]
								c.Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validOneOf.Name+"."+validFld.Name+" one of field", func(c C) {
									c.So(fld.Name, ShouldEqual, validFld.Name)
									c.So(fld.QuotedComment, ShouldEqual, validFld.QuotedComment)
									CompareTypes(c, fld.Type, validFld.Type)
								})
							}

						})
					}
				})

			}
		})
		c.Convey("test.proto Should contain valid services", func(c C) {
			c.So(test.Services, ShouldHaveLength, len(f.Services))
			for i, srv := range test.Services {
				validSrv := f.Services[i]
				c.Convey("Should have valid parsed "+validSrv.Name+" service ", func(c C) {
					c.So(srv.Name, ShouldEqual, validSrv.Name)
					c.So(srv.QuotedComment, ShouldEqual, validSrv.QuotedComment)
					c.Convey(validSrv.Name+" should contains valid methods", func(c C) {
						c.So(srv.Methods, ShouldHaveLength, len(validSrv.Methods))
						for i, method := range srv.Methods {
							validMethod := validSrv.Methods[i]
							c.Convey(validSrv.Name+" should contains valid "+validMethod.Name+" method", func(c C) {
								c.So(method.Name, ShouldEqual, validMethod.Name)
								c.So(method.QuotedComment, ShouldEqual, validMethod.QuotedComment)
								c.Convey(validSrv.Name+"."+validMethod.Name+" should have valid input message type", func(c C) {
									CompareTypes(c, method.InputMessage, validMethod.InputMessage)
								})
								c.Convey(validSrv.Name+"."+validMethod.Name+" should have valid output message type", func(c C) {
									CompareTypes(c, method.OutputMessage, validMethod.OutputMessage)
								})
							})
						}
					})
				})
			}
		})
	})
}

func CompareTypes(c C, t1, t2 Type) {
	c.So(t1, ShouldNotBeNil)
	c.So(t2, ShouldNotBeNil)

	switch protoType := t1.(type) {
	case *Scalar:
		c.So(protoType.ScalarName, ShouldEqual, t2.(*Scalar).ScalarName)
	case *Message:
		c.So(t1, ShouldEqual, t2)
		c.So(t1.File(), ShouldEqual, t2.File())
	case *Enum:
		c.So(t1, ShouldEqual, t2)
		c.So(t1.File(), ShouldEqual, t2.File())
	case *Map:
		c.So(t1.(*Map).Message, ShouldEqual, t2.(*Map).Message)
		CompareTypes(c, t1.(*Map).KeyType, t2.(*Map).KeyType)
		CompareTypes(c, t1.(*Map).ValueType, t2.(*Map).ValueType)
		c.So(t1.File(), ShouldEqual, t2.File())
	default:
		panic("Undefined type")
	}
}
