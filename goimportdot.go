package main

import (
	"flag"
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/yqylovy/goimportdot/core"
)

func convertpath(pkgimports map[string]core.StrSet) map[string]core.StrSet {
	if filepath.Separator != '\\' {
		return pkgimports
	}
	pkgimportsnew := make(map[string]core.StrSet)
	for pkgName, deps := range pkgimports {
		newPkgName := strings.Replace(pkgName, "\\", "/", -1)
		newdep := make(core.StrSet)
		pkgimportsnew[newPkgName] = newdep
		for dep, _ := range deps {
			newdepName := strings.Replace(dep, "\\", "/", -1)
			newdep[newdepName] = deps[dep]
		}
	}
	return pkgimportsnew
}
func main() {
	var ignoreGit = true
	var ignoreTest = true
	var onlySelfPkg = true

	var packageName = ""
	var root = ""
	var filters = ""

	var level = -1

	flag.BoolVar(&ignoreGit, "ignoregit", ignoreGit, "ignore files in git")
	flag.BoolVar(&ignoreTest, "ignoretest", ignoreTest, "ignore test files")
	flag.BoolVar(&onlySelfPkg, "only", onlySelfPkg, "only to draw the input package")
	flag.StringVar(&filters, "filter", "", "filter to (ignore/only include) package match wildcard,example: -filter=w:a*,*b;b:c means only include package start with a and ends with b, ignore package named c")
	flag.StringVar(&root, "root", root, "only draw package with the graph start from root")
	flag.IntVar(&level, "level", level, "show how many level , -1 for all")
	flag.StringVar(&packageName, "pkg", packageName, "the package to draw")
	flag.Parse()

	if packageName == "" {
		fmt.Println("You must specify the packge name with -pkg ")
		return
	}

	fileFilter := []core.FileFilter{
		core.HasSuffix(false, ".go"),
	}
	if ignoreGit {
		fileFilter = append(fileFilter, core.NameContains(true, ".git"))
	}
	if ignoreTest {
		fileFilter = append(fileFilter, core.NameContains(true, "_test.go"))
	}

	pkgAndImports, err := core.GetImports(packageName, fileFilter...)
	if err != nil {
		panic(err)
	}
	pkgAndImports = convertpath(pkgAndImports)
	pkgFilters := []core.PkgFilter{}
	if onlySelfPkg {
		pkgFilters = append(pkgFilters, core.PkgWildcardFilter(false, packageName+"*"))
	}
	if root != "" {
		pkgFilters = append(pkgFilters, core.RootFilter(root))
	}
	moreFilters, err := core.ParsePkgWildcardStr(filters)
	if err != nil {
		fmt.Printf("No right filter [%s], please check!", filters)
		return
	}
	pkgFilters = append(pkgFilters, moreFilters...)

	if level >= 0 {
		pkgFilters = append(pkgFilters, core.PkgLevelFilter(level))
	}

	for _, f := range pkgFilters {
		pkgAndImports = f(pkgAndImports)
	}
	core.WriteDot(pkgAndImports, os.Stdout)
}
