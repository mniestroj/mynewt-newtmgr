/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cli

import (
	"fmt"
	"os"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"mynewt.apache.org/newt/util"
)

var onExit func()
var exiting int32
var silenceErrors bool

func SetOnExit(cb func()) {
	onExit = cb
}

func SilenceErrors() {
	silenceErrors = true
}

// Performs some cleanup and terminates the application.
func NmExit(status int) {
	// If we are already exiting, just block forever.  We don't want to perform
	// a second round of cleanup or quit before the current one completes.
	if !atomic.CompareAndSwapInt32(&exiting, 0, 1) {
		select {}
	}

	if onExit != nil {
		onExit()
	}
	os.Exit(status)
}

func nmUsage(cmd *cobra.Command, err error) {
	if !silenceErrors {
		if err != nil {
			sErr, ok := err.(*util.NewtError)
			if !ok {
				sErr = util.ChildNewtError(err)
			}

			log.Debugf("%s", sErr.StackTrace)
			fmt.Fprintf(os.Stderr, "Error: %s\n", sErr.Text)
		}

		if cmd != nil {
			fmt.Printf("\n")
			fmt.Printf("%s - ", cmd.Name())
			cmd.Help()
		}
	}

	NmExit(1)
}
