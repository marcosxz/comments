package comments

import (
	"go/ast"
	"os"
)

type Finder interface {
	Filter(os.FileInfo) bool
	Func(FuncComments) error
	Type(TypeComments) error
	Field(FieldComments) error
}

type FuncComments struct {
	Package  string
	Filepath string
	FuncName string
	Caller   string
	FuncDecl *ast.FuncDecl
}

type TypeComments struct {
	Package  string
	Filepath string
	TypeName string
	GenDecl  *ast.GenDecl
}

type FieldComments struct {
	Package  string
	Filepath string
	TypeName string
	Fields   *ast.FieldList
}
