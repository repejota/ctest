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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// CTest is the main type of the program
type CTest struct {
	watchPaths []string

	watchExtensions []string

	mu         sync.Mutex
	watchFiles map[string]time.Time
}

// NewCTest creates a new instance
func NewCTest(extensions, paths []string) (*CTest, error) {
	ctest := &CTest{
		watchPaths:      paths,
		watchExtensions: extensions,
		watchFiles:      make(map[string]time.Time),
	}

	// if paths is empty, then use current directory
	if len(paths) == 0 {
		ctest.watchPaths = []string{GetCurrentDirectory()}
	}

	log.Printf("Watching %d paths %q", len(ctest.watchPaths), ctest.watchPaths)

	// if extensions is empty, then use *.go files
	if len(extensions) == 0 {
		ctest.watchExtensions = []string{".go"}
	}

	log.Printf("Watching %d extensions %q", len(ctest.watchExtensions), ctest.watchExtensions)

	for _, watchPath := range ctest.watchPaths {
		err := ctest.getFilesToWatch(watchPath, true)
		if err != nil {
			return ctest, err
		}
	}

	log.Printf("Watching %d files for changes", len(ctest.watchFiles))

	return ctest, nil
}

// getFilesToWatch build the list of files to watch for changes
func (c *CTest) getFilesToWatch(watchPath string, recursive bool) error {
	log.Printf("Walking %s recursively", watchPath)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
		if info.IsDir() && path != watchPath && !recursive {
			return filepath.SkipDir
		}
		for _, extension := range c.watchExtensions {
			if filepath.Ext(path) == extension {
				c.mu.Lock()
				c.watchFiles[path] = info.ModTime()
				c.mu.Unlock()
				log.Printf("Watching %s", path)
			}
		}
		return nil
	}
	err := filepath.Walk(watchPath, walkFunc)
	if err != nil {
		return err
	}
	return nil
}

// Start starts the main loop
func (c *CTest) Start() {
	for {
		c.handleChanges()
		time.Sleep(time.Duration(1 * time.Second)) // 1 second delay
	}
}

// handleChanges handles file changes
func (c *CTest) handleChanges() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for file, modtime := range c.watchFiles {
		stat, err := os.Stat(file)
		if err != nil {
			log.Printf("Changed file %s and got an error %s", file, err.Error())
		}
		ntime := stat.ModTime()
		if ntime.Sub(modtime) > 0 {
			c.watchFiles[file] = ntime
			log.Printf("Changed file %s", file)
			// execute tests
			c.RunTests()
		}
	}
}

// RunTests runs tests
func (c *CTest) RunTests() bool {
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	fmt.Println("[STDOUT]", cmd.Stdout)
	fmt.Println("[STDERR]", cmd.Stderr)
	if err != nil {
		return false
	}
	return true
}
