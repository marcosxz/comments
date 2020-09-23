package comments

import (
	"fmt"
	"go/ast"
	"os"
	"strings"
)

const (
	funcCommentsTmp = `
THE FUNC DOC
Package: %s
File: %s
FuncName: %s
CallerName: %s
FuncDocs: %s
CallerDocs: %s
`
	typeCommentsTmp = `
THE TYPE DOC
Package: %s 
File: %s
Type: %s
Docs: %s
`
	fieldCommentsTmp = `
THE FIELD DOC
Package: %s
File: %s
Type: %s
Docs: %s
`
	packageCommentsTmp = `
THE PACKAGE DOC
Package: %s
File: %s
Docs: %s
`
)

type printFinder struct{}

func (printFinder) Filter(os.FileInfo) bool { return true }

func (printFinder) Package(comments PackageComments) error {
	fmt.Println(fmt.Sprintf(packageCommentsTmp, comments.Pkg.Name,
		comments.Filepath, strings.TrimSuffix(comments.File.Doc.Text(), "\n")))
	return nil
}

func (printFinder) Func(comments FuncComments) error {
	var callerName string
	var callerDocs string
	if comments.Func.Recv.NumFields() > 0 {
		expr := comments.Func.Recv.List[0].Type
	loop:
		switch x := expr.(type) {
		case *ast.Ident: // not pointer call
			callerName = x.Name
			if x.Obj != nil {
				switch t := x.Obj.Decl.(type) {
				case *ast.TypeSpec:
					switch s := t.Type.(type) {
					case *ast.StructType:
						for _, field := range s.Fields.List {
							if len(field.Names) == 0 {
								continue
							}
							callerDocs += field.Names[0].Name + ":"
							if field.Doc.Text() != "" {
								callerDocs += field.Doc.Text()
							} else {
								callerDocs += "\n"
							}
						}
					}
				}
			}
		case *ast.StarExpr: // the pointer call
			//caller = fmt.Sprint(x.X)
			switch x := x.X.(type) {
			case *ast.Ident:
				expr = x
				goto loop
			case *ast.StarExpr:
				expr = x
				goto loop
			}
		}
	}
	fmt.Println(fmt.Sprintf(funcCommentsTmp,
		comments.Pkg.Name, comments.Filepath, comments.Func.Name.String(), callerName,
		strings.TrimSuffix(comments.Func.Doc.Text(), "\n"), strings.TrimSuffix(callerDocs, "\n")))
	return nil
}

func (printFinder) Type(comments TypeComments) error {
	fmt.Println(fmt.Sprintf(typeCommentsTmp, comments.Pkg.Name, comments.Filepath,
		comments.Type.Name.String(), strings.TrimSuffix(comments.Decl.Doc.Text(), "\n")))
	return nil
}

func (printFinder) Field(comments FieldComments) error {
	docs := ""
	for _, field := range comments.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		docs += field.Names[0].Name + ":"
		if field.Doc.Text() != "" {
			docs += field.Doc.Text()
		} else {
			docs += "\n"
		}
	}
	fmt.Println(fmt.Sprintf(fieldCommentsTmp, comments.Pkg.Name,
		comments.Filepath, comments.Type.Name.String(), strings.TrimSuffix(docs, "\n")))
	return nil
}
