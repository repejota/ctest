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
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/repejota/ctest/git"
)

// File ...
type File struct {
	Name    string
	Lines   int
	Package *git.Package
}

// NewFile ...
func NewFile(name string, pkg *git.Package) *File {
	f := &File{
		Name:    name,
		Package: pkg,
	}
	ff, _ := os.Open(f.FullName())
	defer ff.Close()
	f.Lines, _ = lineCounter(ff)
	return f
}

// FullName ...
func (f *File) FullName() string {
	fullname := fmt.Sprintf("%s/%s", f.Package.Dir, f.Name)
	return fullname
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
