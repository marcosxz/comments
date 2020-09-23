package comments

import (
	"go/ast"
	"os"
)

type Finder interface {
	Filter(os.FileInfo) bool
	Package(PackageComments) error
	Func(FuncComments) error
	Type(TypeComments) error
	Field(FieldComments) error
}

type PackageComments struct {
	Filepath string
	Pkg      *ast.Package
	File     *ast.File
}

type FuncComments struct {
	PackageComments
	Func *ast.FuncDecl
}

type TypeComments struct {
	PackageComments
	Decl *ast.GenDecl
	Type *ast.TypeSpec
}

type FieldComments struct {
	TypeComments
	Fields *ast.FieldList
}
