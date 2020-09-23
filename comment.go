package comments

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func AstFind(dir string, all bool, finder Finder) error {
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
		if err = find(dir, finder); err != nil {
			return err
		}
		if !all && i == 0 {
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
	for _, pkg := range packages {
		for filepath, file := range pkg.Files {
			if err = packageComments(finder, pkg, file, filepath); err != nil {
				return err
			}
			for _, decl := range file.Decls {
				switch x := decl.(type) {
				case *ast.FuncDecl: // func declaration
					if err = funcComments(finder, pkg, file, filepath, x); err != nil {
						return err
					}
				case *ast.GenDecl: // import, constant, type or variable declaration
					if err = typeAndFieldComments(finder, pkg, file, filepath, x); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func packageComments(finder Finder, pkg *ast.Package, file *ast.File, filepath string) error {
	return finder.Package(PackageComments{
		Filepath: filepath,
		Pkg:      pkg,
		File:     file,
	})
}

func funcComments(finder Finder, pkg *ast.Package, file *ast.File, filepath string, decl *ast.FuncDecl) error {
	pkgComments := PackageComments{
		Filepath: filepath,
		Pkg:      pkg,
		File:     file,
	}
	return finder.Func(FuncComments{
		PackageComments: pkgComments,
		Func:            decl,
	})
}

func typeAndFieldComments(finder Finder, pkg *ast.Package, file *ast.File, filepath string, decl *ast.GenDecl) error {
	pkgComments := PackageComments{
		Filepath: filepath,
		Pkg:      pkg,
		File:     file,
	}
	for _, spec := range decl.Specs {
		switch x := spec.(type) {
		case *ast.TypeSpec:
			typeComments := TypeComments{
				PackageComments: pkgComments,
				Decl:            decl,
				Type:            x,
			}
			if err := finder.Type(typeComments); err != nil {
				return err
			}
			switch ex := x.Type.(type) {
			case *ast.StructType:
				fieldComments := FieldComments{
					TypeComments: typeComments,
					Fields:       ex.Fields,
				}
				if err := finder.Field(fieldComments); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
