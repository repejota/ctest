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
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/repejota/ctest/ui"
	log "github.com/sirupsen/logrus"
)

// CTest is the main type of the program
type CTest struct {
	watchPaths []string

	watchExtensions []string

	mu         sync.Mutex
	watchFiles map[string]time.Time
}

// NewCTest creates a new instance
func NewCTest(extensions, paths []string, recursive bool) (*CTest, error) {
	ctest := &CTest{
		watchPaths:      paths,
		watchExtensions: extensions,
		watchFiles:      make(map[string]time.Time),
	}

	// if paths is empty, then use current directory
	if len(paths) == 0 {
		cwd, _ := os.Getwd()
		ctest.watchPaths = []string{cwd}
	}

	log.Infof("Watching %d paths: %q", len(ctest.watchPaths), ctest.watchPaths)

	// if extensions is empty, then use *.go files
	if len(extensions) == 0 {
		ctest.watchExtensions = []string{".go"}
	}

	log.Infof("Watching %d extensions: %q", len(ctest.watchExtensions), ctest.watchExtensions)

	for _, watchPath := range ctest.watchPaths {
		err := ctest.getFilesToWatch(watchPath, recursive)
		if err != nil {
			return ctest, err
		}
	}

	log.Infof("Watching %d files", len(ctest.watchFiles))

	return ctest, nil
}

// getFilesToWatch build the list of files to watch
func (c *CTest) getFilesToWatch(watchPath string, recursive bool) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		path, _ = filepath.Abs(path)
		_, err = os.Stat(path)
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
				log.Debugf("Watching: %s", path)
			}
		}
		return nil
	}

	log.Debugf("Walking: %s", watchPath)
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

// StartUI starts the UI server
func (c *CTest) StartUI() {
	r := mux.NewRouter()
	r.HandleFunc("/", ui.HomeHandler)
	r.HandleFunc("/test", ui.TestHandler)
	r.HandleFunc("/cover", ui.CoverHandler)
	r.HandleFunc("/git", ui.GitHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// handleChanges handles file changes
func (c *CTest) handleChanges() {
	for file, modtime := range c.watchFiles {
		stat, err := os.Stat(file)
		if err != nil {
			log.Errorf("can't get file info: %s", err.Error())
		}
		ntime := stat.ModTime()
		if ntime.Sub(modtime) > 0 {
			log.Debugf("Changed file: %s", file)

			c.mu.Lock()
			c.watchFiles[file] = ntime
			c.mu.Unlock()

			c.RunTests("go", "test", "-v", "./...")
		}
	}
}

// RunTests runs tests
func (c *CTest) RunTests(command string, args ...string) bool {
	cmd := exec.Command(command)
	cmd.Args = args
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("tests failed: %v", err)
		return false
	}
	return true
}
