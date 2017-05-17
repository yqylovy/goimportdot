package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPkgRegex(t *testing.T) {
	Convey("Test PkgRegex", t, func() {
		imps := map[string]StrSet{
			"a/b/c": NewStrSet(
				"test/subt",
			),
		}
		filter := PkgWildcardFilter(true, "test*")
		imps = filter(imps)
		So(imps["a/b/c"], ShouldBeEmpty)
	})
}

func TestParsePkgWildcardStr(t *testing.T) {
	Convey("Test ParsePkgWildcardStr", t, func() {
		str := "w:a*,*b;b:c"
		fs, err := ParsePkgWildcardStr(str)
		So(err, ShouldBeNil)
		So(len(fs), ShouldEqual, 2)
	})
}
