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
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/repejota/ctest"
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

func init() {
	log.SetFlags(0)

	t := time.Now()
	tf := t.Format(time.RFC3339)
	prefix := fmt.Sprintf("ctest[%s]: ", tf)
	log.SetPrefix(prefix)

	log.SetOutput(ioutil.Discard)
}

func main() {
	verbosePtr := flag.Bool("verbose", false, "Enable verbose mode")
	versionPtr := flag.Bool("version", false, "Show version information")
	var extension string
	flag.StringVar(&extension, "extension", "*.go", "Extensions to watch")
	flag.Parse()
	if *verbosePtr {
		log.SetOutput(os.Stderr)
	}
	if *versionPtr {
		output := ctest.ShowVersionInfo(Version, Build)
		fmt.Println(output)
		os.Exit(0)
	}

	paths := flag.Args()

	log.Println("------")
	log.Println(extension)
	log.Println("------")

	ctest := ctest.NewCTest(paths)

	ctest.Start()
}
