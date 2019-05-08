package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMessageMethods(t *testing.T) {
	Convey("Test Message.HaveFields", t, func(c C) {
		c.Convey("Should return true, if there's a normal field", func(c C) {
			c.So(Message{NormalFields: []*NormalField{{}}}.HaveFields(), ShouldBeTrue)
		})
		c.Convey("Should return true, if there's a map field", func(c C) {
			c.So(Message{MapFields: []*MapField{{}}}.HaveFields(), ShouldBeTrue)
		})
		c.Convey("Should return true, if there's a oneof", func(c C) {
			c.So(Message{OneOffs: []*OneOf{{Fields: []*NormalField{{}}}}}.HaveFields(), ShouldBeTrue)
		})
		c.Convey("Should return false, if there no fields", func(c C) {
			c.So(Message{}.HaveFields(), ShouldBeFalse)
		})
	})
	Convey("Test Message.HaveFieldsExcept", t, func(c C) {
		c.Convey("Should return true, if there's a normal field", func(c C) {
			msg := Message{NormalFields: []*NormalField{
				{Name: "a"},
				{Name: "b"},
			}}
			c.So(msg.HaveFieldsExcept("a"), ShouldBeTrue)
		})
		c.Convey("Should return true, if there's a map field", func(c C) {
			msg := Message{MapFields: []*MapField{
				{Name: "a"},
				{Name: "b"},
			}}
			c.So(msg.HaveFieldsExcept("b"), ShouldBeTrue)
		})
		c.Convey("Should return true, if there's a oneof", func(c C) {
			msg := Message{OneOffs: []*OneOf{
				{
					Fields: []*NormalField{
						{Name: "a"},
						{Name: "b"},
						{Name: "c"},
					},
				},
			}}
			c.So(msg.HaveFieldsExcept("b"), ShouldBeTrue)
		})
		c.Convey("Should return false, if there's only excepted normal field", func(c C) {
			msg := Message{NormalFields: []*NormalField{
				{Name: "a"},
			}}
			c.So(msg.HaveFieldsExcept("a"), ShouldBeFalse)
		})
		c.Convey("Should return false, if there's only excepted map field", func(c C) {
			msg := Message{MapFields: []*MapField{
				{Name: "b"},
			}}
			c.So(msg.HaveFieldsExcept("b"), ShouldBeFalse)
		})

		c.Convey("Should return false, if there's only excepted oneof field", func(c C) {
			msg := Message{OneOffs: []*OneOf{
				{
					Fields: []*NormalField{
						{Name: "b"},
					},
				},
			}}
			c.So(msg.HaveFieldsExcept("b"), ShouldBeFalse)
		})
		c.Convey("Should return false, if there's no filelds", func(c C) {
			msg := Message{}
			c.So(msg.HaveFieldsExcept("b"), ShouldBeFalse)
		})
	})
}
