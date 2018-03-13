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

package ctest

import (
	"log"
)

// CTest is the main type of the program
type CTest struct {
	watchPaths      []string
	watchExtensions []string
	watchFiles      []string
}

// NewCTest creates a new instance
func NewCTest(paths []string) *CTest {
	ctest := &CTest{
		watchPaths:      paths,
		watchExtensions: make([]string, 0),
		watchFiles:      make([]string, 0),
	}

	// if paths is empty, then use current directory
	if len(paths) == 0 {
		ctest.watchPaths = []string{GetCurrentDirectory()}
	}

	log.Println("Watching paths", ctest.watchPaths)

	//ctest.watchExtensions = []string{"*.go"}

	log.Println("Watching extensions", ctest.watchExtensions)

	log.Println("Watching files", ctest.watchFiles)

	return ctest
}

// Start starts the program main loop
func (c *CTest) Start() {

}
