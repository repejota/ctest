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

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Version is the release number
	//
	// This number is the latest tag from the git repository.
	Version string
	// Build is the build string
	//
	// This string is the branch name and the commit hash (short format)
	Build string
)

func main() {
	versionPtr := flag.Bool("version", false, "Show version information")
	flag.Parse()
	if *versionPtr {
		fmt.Println("ctest : Version", Version, "Build", Build)
		os.Exit(0)
	}
}
