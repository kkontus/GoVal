package main

import (
	"fmt"
	"os"
	"go/parser"
	"go/token"
	"go/ast"
	"path/filepath"
	"errors"
	. "GoVal/config"
	"GoVal/util"
	"text/tabwriter"
)

var cPkg = 0
var cFiles = 0
var cInternalFn = 0
var cExportedFn = 0
var cNoDocFn = 0
var cWithDocFn = 0

func main() {
	fmt.Printf("Hello GoVal\n\n")

	parseDir(".", true, false)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.TabIndent)
	fmt.Fprintln(w, "Packages\tFiles\tFunctions\tInternal\tExported\tNo docs\tWith docs\t")
	fmt.Fprintf(w, "%d\t%d\t%d\t%d\t%d\t%d\t%d\t", cPkg, cFiles, cInternalFn+cExportedFn, cInternalFn, cExportedFn, cNoDocFn, cWithDocFn)
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	w.Flush()
}

func parseDir(dir string, recursive bool, addDesc bool) {
	dirFile, err := os.Open(dir)
	if err != nil {
		util.ShowError(err)
	}

	defer dirFile.Close()

	info, err := dirFile.Stat()
	if err != nil {
		util.ShowError(err)
	}

	if !info.IsDir() {
		util.ShowError(errors.New("Path is not a valid directory: " + dir))
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, Filter, 0)

	if err != nil {
		util.ShowError(err)
	}

	for _, pkg := range pkgs {
		cPkg++

		funcs := []*ast.FuncDecl{}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 8, '\t', tabwriter.TabIndent)
		if addDesc {
			fmt.Fprintln(w, "Package\tFile\tFunction\tLn. Start\tLn. End\tLines\t")
		}

		for file, f := range pkg.Files {
			cFiles++

			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					funcs = append(funcs, fn)

					if addDesc {
						fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\t", pkg.Name, file, fn.Name.Name, fset.Position(fn.Pos()).Line, fset.Position(fn.End()).Line, (fset.Position(fn.End()).Line-fset.Position(fn.Pos()).Line)+1)
						fmt.Fprintln(w)
					}

					if fn.Doc.Text() == "" {
						cNoDocFn++
					} else {
						cWithDocFn++
					}

					if fn.Name.IsExported() {
						cExportedFn++
					} else {
						cInternalFn++
					}
				}
			}
		}
		if addDesc {
			fmt.Fprintln(w)
			fmt.Fprintln(w)
		}
		w.Flush()
	}

	if recursive {
		dirs, err := dirFile.Readdir(-1)
		if err != nil {
			util.ShowError(err)
		}
		for _, info := range dirs {
			if info.IsDir() {
				parseDir(filepath.Join(dir, info.Name()), recursive, addDesc)
			}
		}
	}
}