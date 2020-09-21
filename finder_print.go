package comments

import (
	"fmt"
	"os"
)

const (
	funcComments = `
THE FUNC DOC
Package: %s
File: %s
FuncName: %s
Caller: %s
Docs: 
%s
`
	typeComments = `
THE TYPE DOC
Package: %s 
File: %s
Type: %s
Docs:
%s
`
	fieldComments = `
THE FIELD DOC
Package: %s
File: %s
Type: %s
Docs: 
%s
`
)

type printFinder struct{}

func (printFinder) Filter(os.FileInfo) bool { return true }

func (printFinder) Func(comments FuncComments) error {
	fmt.Println(fmt.Sprintf(funcComments, comments.Package, comments.Filepath, comments.FuncName, comments.Caller, comments.FuncDecl.Doc.Text()))
	return nil
}

func (printFinder) Type(comments TypeComments) error {
	fmt.Println(fmt.Sprintf(typeComments, comments.Package, comments.Filepath, comments.TypeName, comments.GenDecl.Doc.Text()))
	return nil
}

func (printFinder) Field(comments FieldComments) error {
	docs := ""
	for _, field := range comments.Fields.List {
		docs += field.Names[0].Name + ":"
		if field.Doc.Text() != "" {
			docs += field.Doc.Text()
		} else {
			docs += "\n"
		}
	}
	fmt.Println(fmt.Sprintf(fieldComments, comments.Package, comments.Filepath, comments.TypeName, docs))
	return nil
}
