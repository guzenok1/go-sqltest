// Package model contains the data model necessary for generating sqlmock implementations.
package model

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"reflect"
)

const (
	ImportPath = "github.com/guzenok/go-sqltest/sqlmockgen/model"
	compiler   = "source"
)

type (
	InitDataFunc func(a, b string) error
	SqlsDictFunc func() ([]string, error)

	// Package is a Go package.
	Package struct {
		SrcDir string
		Name   string
		Data   map[string]struct{}
		Sqls   map[string]struct{}
	}
)

var (
	typeofInitDataFunc types.Type
	typeofSqlsDictFunc types.Type
)

func init() {
	goImporter := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	pkg, err := goImporter.Import(ImportPath)
	if err != nil {
		panic(err)
	}
	scope := pkg.Scope()

	var (
		f1     InitDataFunc
		f1name = reflect.TypeOf(f1).Name()
		f2     SqlsDictFunc
		f2name = reflect.TypeOf(f2).Name()
	)
	typeofInitDataFunc = scope.Lookup(f1name).Type()
	typeofSqlsDictFunc = scope.Lookup(f2name).Type()
}

func Build(path string) (model *Package, err error) {

	goImporter := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	pkg, err := goImporter.Import(path)
	if err != nil {
		return
	}

	model = &Package{
		Name: pkg.Name(),
		Data: make(map[string]struct{}),
		Sqls: make(map[string]struct{}),
	}

	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue
		}

		funcType, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		if types.AssignableTo(typeofInitDataFunc, funcType) {
			model.Data[name] = struct{}{}
			continue
		}

		if types.AssignableTo(typeofSqlsDictFunc, funcType) {
			model.Sqls[name] = struct{}{}
			continue
		}
	}

	return
}

func (pkg *Package) Print(w io.Writer) {
	fmt.Fprintf(w, "package %s\n", pkg.Name)
}
