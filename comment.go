package comments

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func Find(dir string, all bool, finder Finder) error {
	if finder == nil {
		finder = printFinder{}
	}
	dirs := []string{dir}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, path.Join(dir, file.Name()))
		}
	}
	var i int
	for i, dir = range dirs {
		if err := find(dir, finder); err != nil {
			return err
		}
		if i == 0 && !all {
			break
		}
	}
	return nil
}

func find(dir string, finder Finder) error {
	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && finder.Filter(info)
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	for pkgName, pkg := range packages {
		for filename, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				//ast.Print(fileSet, node)
				switch n := node.(type) {
				case *ast.FuncDecl: // func declaration
					comments := FuncComments{
						Package:  pkgName,
						Filepath: filename,
						FuncName: n.Name.Name,
						FuncDecl: n,
					}
					if n.Recv.NumFields() > 0 {
						t := n.Recv.List[0].Type
						switch t.(type) {
						case *ast.Ident: // not pointer call
							comments.Caller = t.(*ast.Ident).Name
						case *ast.StarExpr: // the pointer call
							comments.Caller = fmt.Sprint(t.(*ast.StarExpr).X)
						}
					}
					if err = finder.Func(comments); err != nil {
						return false
					}
				case *ast.GenDecl: // import, constant, type or variable declaration
					for _, spec := range n.Specs {
						switch t := spec.(type) {
						case *ast.TypeSpec:
							if err = finder.Type(TypeComments{
								Package:  pkgName,
								Filepath: filename,
								TypeName: t.Name.Name,
								GenDecl:  n,
							}); err != nil {
								return false
							}
							switch s := t.Type.(type) {
							case *ast.StructType:
								if err = finder.Field(FieldComments{
									Package:  pkgName,
									Filepath: filename,
									TypeName: t.Name.Name,
									Fields:   s.Fields,
								}); err != nil {
									return false
								}
							}
						}
					}
				}
				return true
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
