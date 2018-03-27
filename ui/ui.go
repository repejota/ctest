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

package ui

import (
	"html/template"
	"log"
	"net/http"

	"github.com/repejota/ctest/git"
)

// HomeHandler ...
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Packages []string
	}
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, data)
}

// TestHandler ...
func TestHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Packages []string
	}
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("templates/test.html"))
	tmpl.Execute(w, data)
}

// CoverHandler ...
func CoverHandler(w http.ResponseWriter, r *http.Request) {
	packageImportsStrings, err := git.ListPackages()
	if err != nil {
		log.Println(err)
	}
	packages, err := git.GetPackages(packageImportsStrings...)
	if err != nil {
		log.Println(err)
	}
	var data struct {
		Packages []*git.Package
	}
	data.Packages = packages
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("templates/cover.html"))
	tmpl.Execute(w, data)
}

// GitHandler ...
func GitHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Packages []string
	}
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("templates/git.html"))
	tmpl.Execute(w, data)
}
