package goodconfig

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGoodConfig(t *testing.T) {

	Config := NewConfig()
	err := Config.Parse("./example/config/application.ini")

	Convey("Parse example/application.ini without errors", t, func(){
			So(err, ShouldEqual, nil)
		})

	Convey("Sections", t, func() {
			Convey("Should parse development and production sections", func(){
					So(len(Config.Sections), ShouldEqual, 2)
				})
		})
	Convey("Values", t, func() {
			Convey("Int", func() {
					Convey("Should return Int type", func() {
							So(Config.Section("production").Get("intValue").ToInt(), ShouldHaveSameTypeAs, 1)
						})
					Convey("Should return the correct value", func() {
							So(Config.Section("production").Get("intValue").ToInt(), ShouldEqual, 33)
						})
				})
			Convey("String", func() {
					Convey("Should return String type", func() {
							So(Config.Section("production").Get("stringValue").ToString(), ShouldHaveSameTypeAs, "string")
						})
					Convey("Should return the correct value", func() {
							So(Config.Section("production").Get("stringValue").ToString(), ShouldEqual, "string me!")
						})
				})
			Convey("Bool", func() {
					Convey("Should return Bool type", func() {
							So(Config.Section("production").Get("intValue").ToBool(), ShouldHaveSameTypeAs, true)
						})
					Convey("Should return the correct value", func() {
							So(Config.Section("production").Get("boolValue").ToBool(), ShouldEqual, true)
							So(Config.Section("development").Get("boolValue").ToBool(), ShouldEqual, false)
						})
				})
			Convey("Map", func() {
					Convey("Should return Map type", func() {
							So(Config.Section("production").Get("parent").ToMap(), ShouldHaveSameTypeAs, make(map[string]Record))
						})
				})
			Convey("Nesting", func() {
					Convey("Should return nested value", func(){
							So(Config.Section("production").Get("parent").Get("child").Get("subchild").ToString(), ShouldEqual, "hello")
						})

					Convey("Should store the errors while getting the nested value", func() {
							So(Config.Section("production").Get("parent").Get("childdd"), ShouldHaveSameTypeAs, &Record{})
							So(Config.HasErrors(), ShouldEqual, true)
						})
					Convey("Should correctly process the defaultSection", func() {
							Config.SetDefaultSection("development")
							So(Config.Get("boolValue").ToBool(), ShouldEqual, false)
						})
				})
		})
	Convey("Inheritance", t, func() {
			Convey("Should keep the production value", func() {
					So(Config.Section("production").Get("parent").Get("child").Get("subchild").ToString(), ShouldEqual, "hello")
				})
			Convey("Should override the development value", func() {
					So(Config.Section("development").Get("parent").Get("child").Get("subchild").ToString(), ShouldEqual, "bye")
				})
			Convey("Should inherit the unchanged value", func() {
					So(Config.Section("development").Get("parent").Get("child").Get("subchild_2").ToString(), ShouldEqual, "hello_2")
				})
		})

	Convey("Recursive fill", t, func() {
			Convey("Should fill the simple record", func() {
					r1 := Record{Value:444}
					r2 := Record{}

					r2.Fill(r1)

					So(r2.Value.(int), ShouldEqual, 444)
				})
		})
}
