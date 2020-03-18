// Copyright 2020 NSONE, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lg

import (
	"fmt"
	"log"
	"os"
)

const (
	LoggingOff = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var (
	logger       *log.Logger
	rootPriority int
)

func SetLevel(priority int) error {
	if priority < LoggingOff || priority > LevelTrace {
		return fmt.Errorf("invalid priority level %d", priority)
	}

	rootPriority = priority
	return nil
}

func EnabledFor(priority int) bool {
	return priority <= rootPriority
}

func Tracef(format string, v ...interface{}) {
	if EnabledFor(LevelTrace) {
		f := "[TRACE] " + format
		logger.Printf(f, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if EnabledFor(LevelDebug) {
		f := "[DEBUG] " + format
		logger.Printf(f, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if EnabledFor(LevelInfo) {
		f := "[INFO] " + format
		logger.Printf(f, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if EnabledFor(LevelWarn) {
		f := "[WARN] " + format
		logger.Printf(f, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if EnabledFor(LevelError) {
		f := "[ERROR] " + format
		logger.Printf(f, v...)
	}
}

// Printf writes to configured output regardless of the configured log priority.
// Note that if the log priority is set to LoggingOff this output WILL BE suppressed.
func Printf(format string, v ...interface{}) {
	if rootPriority != LoggingOff {
		logger.Printf(format, v...)
	}
}

func init() {
	rootPriority = LevelError

	// Logs are always written to STDERR.
	logger = log.New(os.Stderr, "", 0)
}
