// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package for checking the prohibition of a direct call to the os.Exit function in the main file.
package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer - implements type interface for multicheck.
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  Doc,
	Run:  run,
}

// Const text docs.
const Doc = `
	check os.Exit call in main
	//
	package main
	import (
		....
		"os"
		....
	)
	func main() {
		....
		os.Exit(1)
		....
	}
`

// Looking for a function call.
func dedect(pass *analysis.Pass, e *ast.CallExpr) {
	selectorExpr, ok := e.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return
	}

	if ident.Name == "os" && selectorExpr.Sel.Name == "Exit" {
		pass.Reportf(ident.NamePos, "calling os.Exit in main function")
	}
}

// Start Analyzer.
func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if f.Name.Name != "main" {
			continue
		}

		for _, d := range f.Decls {
			funcDecl, ok := d.(*ast.FuncDecl)
			if ok && funcDecl.Name.Name == "main" {
				ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
					callExpr, callExprOk := node.(*ast.CallExpr)
					if callExprOk {
						dedect(pass, callExpr)
					}
					return true
				})
			}
		}
	}

	return nil, nil
}
