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

package ctest_test

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/repejota/ctest"
)

func TestDummyTest(t *testing.T) {

}

func TestNewCTest(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	_, err := ctest.NewCTest(nil, nil, false)
	if err != nil {
		t.Error(err)
	}
}

func TestNewCTestPathNotExists(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	paths := []string{"/donotexists"}
	_, err := ctest.NewCTest(nil, paths, false)
	if err == nil {
		t.Error(err)
	}
}

func TestNewCTestPathNotExistsRelative(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	paths := []string{"donotexists"}
	_, err := ctest.NewCTest(nil, paths, false)
	if err == nil {
		t.Error(err)
	}
}

func TestRunTests(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	ct, err := ctest.NewCTest(nil, nil, false)
	if err != nil {
		t.Error(err)
	}
	res := ct.RunTests("true")
	if !res {
		t.Error("failed")
	}
}

func TestRunTestsFailCommand(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	ct, err := ctest.NewCTest(nil, nil, false)
	if err != nil {
		t.Error(err)
	}
	res := ct.RunTests("false")
	if res {
		t.Error("not failed")
	}
}
