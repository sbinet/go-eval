// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"exp/eval"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

var fset = token.NewFileSet()
var filename = flag.String("f", "", "file to run")

func main() {
	flag.Parse()
	w := eval.NewWorld()
	if *filename != "" {
		data, err := ioutil.ReadFile(*filename)
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		file, err := parser.ParseFile(fset, *filename, data, 0)
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		files := []*ast.File{file}
		code, err := w.CompilePackage(fset, files, "main")
		if err != nil {
			if list, ok := err.(scanner.ErrorList); ok {
				for _, e := range list {
					fmt.Println(e.String())
				}
			} else {
				fmt.Println(err.String())
			}
			os.Exit(1)
		}
		code, err = w.Compile(fset, "main()")
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		_, err = code.Run()
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		os.Exit(0)
	}

	r := bufio.NewReader(os.Stdin)
	for {
		print("; ")
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		code, err := w.Compile(fset, line)
		if err != nil {
			println(err.String())
			continue
		}
		v, err := code.Run()
		if err != nil {
			println(err.String())
			continue
		}
		if v != nil {
			println(v.String())
		}
	}
}
