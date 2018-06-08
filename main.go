package main

import (
	"fmt"
	"flag"
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

func main() {
	fmt.Printf("Hello GoVal\n\n")

	// define flags
	recursive := flag.Bool(RECURSIVE, true, "recursive")
	addDesc := flag.Bool(DESCRIPTION, true, "addDesc")
	flag.Parse()

	var goVal = &GoVal{0, 0, 0, 0, 0, 0, 0}
	var goDataMap = [] GoData{}

	parseDir(".", *recursive, *addDesc, goVal, &goDataMap)

	tabwriterBasic(goVal)
	tabwriterAdditional(*addDesc, goDataMap)
}

func tabwriterBasic(goVal *GoVal) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.TabIndent)
	fmt.Fprintln(w, "Packages\tFiles\tFunctions\tInternal\tExported\tNo docs\tWith docs\t")
	fmt.Fprintf(w, "%d\t%d\t%d\t%d\t%d\t%d\t%d\t", goVal.Packages, goVal.Files, goVal.Functions, goVal.Internal, goVal.Exported, goVal.NoDocs, goVal.WithDocs)
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	w.Flush()
}

func tabwriterAdditional(addDesc bool, goDataMap [] GoData) {
	if addDesc {
		wx := new(tabwriter.Writer)
		wx.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.TabIndent)
		fmt.Fprintln(wx, "Package\tFile\tFunction\tLn. Start\tLn. End\tLines\t")

		for i := range (goDataMap) {
			goData := goDataMap[i]
			fmt.Fprintf(wx, "%s\t%s\t%s\t%d\t%d\t%d\t", goData.Package, goData.File, goData.Function, goData.Start, goData.End, goData.Lines)
			fmt.Fprintln(wx)
		}
		fmt.Fprintln(wx)
		wx.Flush()
	}
}

func parseDir(dir string, recursive bool, addDesc bool, goVal *GoVal, goDataMap *[] GoData) {
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
		goVal.Packages++

		funcs := []*ast.FuncDecl{}

		for file, f := range pkg.Files {
			goVal.Files++

			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					funcs = append(funcs, fn)

					if addDesc {
						var goData = &GoData{}
						goData.Package = pkg.Name
						goData.File = file
						goData.Function = fn.Name.Name
						goData.Start = fset.Position(fn.Pos()).Line
						goData.End = fset.Position(fn.End()).Line
						goData.Lines = (fset.Position(fn.End()).Line - fset.Position(fn.Pos()).Line) + 1

						*goDataMap = append(*goDataMap, *goData)
					}

					if fn.Doc.Text() == "" {
						goVal.NoDocs++
					} else {
						goVal.WithDocs++
					}

					if fn.Name.IsExported() {
						goVal.Exported++
					} else {
						goVal.Internal++
					}
					goVal.Functions = goVal.Exported + goVal.Internal
				}
			}
		}
	}

	if recursive {
		dirs, err := dirFile.Readdir(-1)
		if err != nil {
			util.ShowError(err)
		}
		for _, info := range dirs {
			if info.IsDir() {
				parseDir(filepath.Join(dir, info.Name()), recursive, addDesc, goVal, goDataMap)
			}
		}
	}
}
