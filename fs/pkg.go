// Copyright 2018 The ctest Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package fs

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// Package ...
type Package struct {
	Dir           string `json:"Dir"`           // directory containing package sources
	ImportPath    string `json:"ImportPath"`    // import path of package in dir
	ImportComment string `json:"ImportComment"` // path in import comment on package statement
	Name          string `json:"Name"`          // package name
	Doc           string `json:"Doc"`           // package documentation string
	Target        string `json:"Target"`        // install path
	Shlib         string `json:"Shlib"`         // the shared library that contains this package (only set when -linkshared)
	Goroot        bool   `json:"Goroot"`        // is this package in the Go root?
	Standard      bool   `json:"Standard"`      // is this package part of the standard Go library?
	Stale         bool   `json:"Stale"`         // would 'go install' do anything for this package?
	StaleReason   string `json:"StaleReason"`   // explanation for Stale==true
	Root          string `json:"Root"`          // Go root or Go path dir containing this package
	ConflictDir   string `json:"ConcflictDir"`  // this directory shadows Dir in $GOPATH
	BinaryOnly    bool   `json:"BinaryOnly"`    // binary-only package: cannot be recompiled from sources

	// Source files
	GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles       []string // .go sources files that import "C"
	IgnoredGoFiles []string // .go sources ignored due to build constraints
	CFiles         []string // .c source files
	CXXFiles       []string // .cc, .cxx and .cpp source files
	MFiles         []string // .m source files
	HFiles         []string // .h, .hh, .hpp and .hxx source files
	FFiles         []string // .f, .F, .for and .f90 Fortran source files
	SFiles         []string // .s source files
	SwigFiles      []string // .swig files
	SwigCXXFiles   []string // .swigcxx files
	SysoFiles      []string // .syso object files to add to archive
	TestGoFiles    []string // _test.go files in package
	XTestGoFiles   []string // _test.go files outside package

	// Cgo directives
	CgoCFLAGS    []string // cgo: flags for C compiler
	CgoCPPFLAGS  []string // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string // cgo: flags for C++ compiler
	CgoFFLAGS    []string // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string // cgo: flags for linker
	CgoPkgConfig []string // cgo: pkg-config names

	// Dependency information
	Imports      []string // import paths used by this package
	Deps         []string // all (recursively) imported dependencies
	TestImports  []string // imports from TestGoFiles
	XTestImports []string // imports from XTestGoFiles

	// Error information
	Incomplete bool            // this package or a dependency has an error
	Error      *PackageError   // error loading package
	DepsErrors []*PackageError // errors loading dependencies
}

// PackageError ...
type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

// ListPackages ...
func ListPackages() ([]string, error) {
	out, err := exec.Command("go", "list", "./...").Output()
	if err != nil {
		return nil, err
	}
	packagesString := string(out)
	packagesParts := strings.Split(packagesString, "\n")
	return packagesParts, nil
}

// GetPackage ...
func GetPackage(importpath string) (*Package, error) {
	out, err := exec.Command("go", "list", "-json", importpath).Output()
	if err != nil {
		return nil, err
	}
	var p *Package
	json.Unmarshal(out, &p)
	return p, nil
}

// GetPackages ...
func GetPackages(imports ...string) ([]*Package, error) {
	var list []*Package
	for _, v := range imports {
		p, err := GetPackage(v)
		if err != nil {
			return list, err
		}
		list = append(list, p)
	}
	return list, nil
}
