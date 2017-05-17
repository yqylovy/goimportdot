package core

import (
	"regexp"

	"strings"
)

type PkgFilter func(map[string]StrSet) map[string]StrSet

func RootFilter(root string) PkgFilter {
	return func(imps map[string]StrSet) (ret map[string]StrSet) {
		ret = make(map[string]StrSet)
		cur := []string{root}
		for len(cur) > 0 {
			newcur := NewStrSet()
			for _, pkg := range cur {
				if pkgimp, ok := imps[pkg]; ok {
					ret[pkg] = pkgimp
					newcur.Merge(pkgimp)
				}
			}
			cur = newcur.Array()
		}
		return ret
	}
}

func PkgWildcardFilter(isBlack bool, pkgs ...string) PkgFilter {
	regs := []*regexp.Regexp{}
	for _, pkg := range pkgs {
		rgp := regexp.MustCompile("^" + strings.Replace(pkg, "*", ".*", -1) + "$")
		regs = append(regs, rgp)
	}
	return func(imps map[string]StrSet) (ret map[string]StrSet) {
		ret = make(map[string]StrSet)
	BIG:
		for pkg, imps := range imps {
			for _, rgp := range regs {
				if isBlack == rgp.MatchString(pkg) {
					continue BIG
				}
			}
			for k := range imps {
				for _, rgp := range regs {
					if isBlack == rgp.MatchString(k) {
						imps.Del(k)
					}
				}
			}
			ret[pkg] = imps
		}
		return ret
	}
}

func ParsePkgWildcardStr(str string) (fs []PkgFilter, err error) {
	if str == "" {
		return
	}
	strArr := strings.Split(str, ";")
	for _, str := range strArr {
		str = strings.TrimSpace(str)
		wb_pkgs := strings.SplitN(str, ":", 2)
		fs = append(fs, PkgWildcardFilter(wb_pkgs[0] == "b", strings.Split(wb_pkgs[1], ",")...))
	}
	return
}

func PkgLevelFilter(level int) PkgFilter {
	return func(imps map[string]StrSet) (ret map[string]StrSet) {
		if level < 0 {
			return imps
		}
		// find the root
		// which are not pointed to
		allTarget := NewStrSet()
		for _, targets := range imps {
			allTarget.Merge(targets)
		}
		levelMap := map[string]int{}
		for pkg := range imps {
			if !allTarget.Contains(pkg) {
				levelMap[pkg] = 0
			}
		}
		for i := 0; i < level; i++ {
			nextLevel := NewStrSet()
			for pkg, pkgLevel := range levelMap {
				if pkgLevel != i {
					continue
				}
				for target := range imps[pkg] {
					nextLevel.Put(target)
				}
			}
			for next := range nextLevel {
				levelMap[next] = i + 1
			}
		}

		ret = make(map[string]StrSet, len(levelMap))
		for pkg, lvl := range levelMap {
			ret[pkg] = imps[pkg]
			if lvl == level {
				continue
			}
		}
		return ret
	}
}
