package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yqylovy/goimportdot/core"
)

func main() {
	var ignoreGit = true
	var ignoreTest = true
	var onlySelfPkg = true
	var prefixWhite = ""
	var prefixBlack = ""
	var nameWhite = ""
	var nameBlack = ""
	var root = ""
	var packageName = ""
	flag.BoolVar(&ignoreGit, "ignoregit", ignoreGit, "ignore files in git")
	flag.BoolVar(&ignoreTest, "ignoretest", ignoreTest, "ignore test files")
	flag.BoolVar(&onlySelfPkg, "only", onlySelfPkg, "only to draw the input package")
	flag.StringVar(&prefixWhite, "prefix_white", prefixWhite, "only include package with this prefix,split by ':' ")
	flag.StringVar(&prefixBlack, "prefix_black", prefixBlack, "ignore package with this prefix,split by ':' ")
	flag.StringVar(&nameWhite, "name_white", nameWhite, "only include package with this name,split by ':' ")
	flag.StringVar(&nameBlack, "name_black", nameBlack, "ignore package with this name,split by ':' ")
	flag.StringVar(&root, "root", root, "only draw package with the graph start from root")

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

	pkgFilters := []core.PkgFilter{}
	if onlySelfPkg {
		pkgFilters = append(pkgFilters, core.PkgPrefixFilter(false, packageName))
	}
	if root != "" {
		pkgFilters = append(pkgFilters, core.RootFilter(root))
	}
	if prefixWhite != "" {
		pkgFilters = append(pkgFilters, core.PkgPrefixFilter(false, strings.Split(prefixWhite, ":")...))
	}
	if prefixBlack != "" {
		pkgFilters = append(pkgFilters, core.PkgPrefixFilter(true, strings.Split(prefixBlack, ":")...))
	}
	if nameWhite != "" {
		pkgFilters = append(pkgFilters, core.PkgNameFilter(false, strings.Split(nameWhite, ":")...))
	}
	if nameBlack != "" {
		pkgFilters = append(pkgFilters, core.PkgNameFilter(true, strings.Split(nameBlack, ":")...))
	}

	for _, f := range pkgFilters {
		pkgAndImports = f(pkgAndImports)
	}
	core.WriteDot(pkgAndImports, os.Stdout)
}
