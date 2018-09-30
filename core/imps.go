package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GetImports(pkg string, filters ...FileFilter) (pkgimports map[string]StrSet, err error) {
	fullpath := ""

	goPath := os.Getenv("GOPATH")
	gopaths := filepath.SplitList(goPath)
	for _, gp := range gopaths {
		fp := filepath.Join(gp, "src", pkg)
		if _, err := os.Stat(fp); err == nil {
			fullpath = fp
		}
	}
	if fullpath == "" {
		err = fmt.Errorf("Can not find package [%s] in GOPATH [%s]", pkg, goPath)
		return
	}
	pkgimports = make(map[string]StrSet)
	filepath.Walk(fullpath, func(fp string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, filter := range filters {
			if !filter.IsBlack {
				continue
			}
			if filter.Func(fp, info, err) {
				return nil
			}
		}
		for _, filter := range filters {
			if filter.IsBlack {
				continue
			}
			if !filter.Func(fp, info, err) {
				return nil
			}
		}
		pkg := PkgOfFile(fp)
		if _, ok := pkgimports[pkg]; !ok {
			pkgimports[pkg] = NewStrSet()
		}
		ss, err := ParseGoImport(fp)
		if err != nil {
			// TODO: better err
			panic(err)
		}
		pkgimports[pkg].Merge(ss)
		return nil
	})
	return
}

func WriteDot(pkgimports map[string]StrSet, writer io.Writer) (err error) {
	nodes := NewStrSet()
	edges := [][2]string{}
	for pkg, imps := range pkgimports {
		nodes.Put(pkg)
		for imp := range imps {
			nodes.Put(imp)
			edges = append(edges, [2]string{pkg, imp})
		}
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("digraph G {\n")
	for _, edge := range edges {
		buf.WriteString(fmt.Sprintf(`"%s"->"%s";`, edge[0], edge[1]))
		buf.WriteByte('\n')
	}
	for pkg, _ := range nodes {
		buf.WriteString(fmt.Sprintf(`"%s";`, pkg))
		buf.WriteByte('\n')
	}
	buf.WriteString("}\n")
	_, err = writer.Write(buf.Bytes())
	return
}
