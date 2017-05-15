package core

import (
	"os"
	"strings"
)

type FileFilter struct {
	IsBlack bool
	Func    func(fp string, info os.FileInfo, err error) bool
}

func NameContains(isblack bool, str string) FileFilter {
	return FileFilter{
		IsBlack: isblack,
		Func: func(fp string, _ os.FileInfo, _ error) bool {
			return strings.Contains(fp, str)
		},
	}
}
func HasSuffix(isblack bool, suffixs ...string) FileFilter {
	return FileFilter{
		IsBlack: isblack,
		Func: func(fp string, _ os.FileInfo, _ error) bool {
			for _, sf := range suffixs {
				if sf == "" {
					continue
				}
				if sf[0] != '.' {
					sf = "." + sf
				}
				if strings.HasSuffix(fp, sf) {
					return true
				}
			}
			return false
		},
	}
}
