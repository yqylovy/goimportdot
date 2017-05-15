package core

import "strings"

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
func PkgNameFilter(isBlack bool, pkgs ...string) PkgFilter {
	ignores := NewStrSet()
	for _, pkg := range pkgs {
		ignores.Put(pkg)
	}
	return func(imps map[string]StrSet) (ret map[string]StrSet) {
		ret = make(map[string]StrSet)
		for pkg, imps := range imps {
			if reverse(isBlack, ignores.Contains(pkg)) {
				continue
			}
			for k := range imps {
				if reverse(isBlack, ignores.Contains(k)) {
					imps.Del(k)
				}
			}
			ret[pkg] = imps
		}
		return ret
	}
}

func PkgPrefixFilter(isBlack bool, prefixs ...string) PkgFilter {
	ignores := NewStrSet()
	for _, pkg := range prefixs {
		ignores.Put(pkg)
	}

	return func(imps map[string]StrSet) (ret map[string]StrSet) {
		ret = make(map[string]StrSet)
	BIG:
		for pkg, imps := range imps {
			for str := range ignores {
				if reverse(isBlack, strings.HasPrefix(pkg, str)) {
					continue BIG
				}
			}
			for k := range imps {
				for str := range ignores {
					if reverse(isBlack, strings.HasPrefix(k, str)) {
						imps.Del(k)
					}
				}
			}
			ret[pkg] = imps
		}
		return ret
	}
}

func reverse(isBlack, result bool) bool {
	if isBlack {
		return result
	} else {
		return !result
	}
}
