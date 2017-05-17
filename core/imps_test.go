package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetImports(t *testing.T) {
	Convey("Test GetImports", t, func() {
		_, err := GetImports("go.uber.org/zap",
			NameContains(true, ".git"),
			NameContains(true, "_test.go"),
			HasSuffix(false, ".go"))
		So(err, ShouldBeNil)
	})
}
