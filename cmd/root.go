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

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/repejota/ctest"
	"github.com/spf13/cobra"
)

var (
	extensionFlag []string
	verboseFlag   bool
	versionFlag   bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "goctest",
	Short: "Watch files and execute tests",
	Long:  `goctest continuouslly watch for file changes and execute tests`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup default logger
		log.SetFlags(0)
		t := time.Now()
		tf := t.Format(time.RFC3339)
		prefix := fmt.Sprintf("ctest[%s]: ", tf)
		log.SetPrefix(prefix)
		log.SetOutput(ioutil.Discard)

		// --verbose
		if verboseFlag {
			log.SetOutput(os.Stderr)
		}

		// --version
		if versionFlag {
			// TODO:
			// * Do not hardcode version and build
			output := ctest.ShowVersionInfo("1.2.3", "abcdefg")
			fmt.Println(output)
			os.Exit(0)
		}

		ctest := ctest.NewCTest(extensionFlag, args)

		ctest.Start()
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.SetUsageFunc(UsageFunc)
	RootCmd.Flags().StringArrayVarP(&extensionFlag, "extension", "", []string{"go"}, "set extensions to watch")
	RootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "enable verbose mode")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "V", false, "show version number")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Unimplemented
}

// UsageFunc prints command usage help message
func UsageFunc(cmd *cobra.Command) error {
	fmt.Println("Usage:")
	fmt.Println("  goctest [flags] [paths]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h, --help		help for misto")
	fmt.Println("  -v, --verbose   	enable verbose mode")
	fmt.Println("  -V, --version   	show version number")
	return nil
}
