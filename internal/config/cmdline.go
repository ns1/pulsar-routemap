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

package config

import (
	"fmt"
	"os"
	"path"

	"go.uber.org/multierr"
)

// CommandLineGlobals are options that apply to all/most commands. These
// are set on the root command level.
type CommandLineGlobals struct {
	// CacheDir is user local directory used to store temporary data.
	CacheDir string

	// Verbosity is the user's preference for logging output.
	Verbosity int

	NS1APIBaseURL string
	NS1APIKey     string
}

// NewCommandLineGlobals creates a new globals with some defaults.
func NewCommandLineGlobals() CommandLineGlobals {
	g := CommandLineGlobals{}

	g.CacheDir = os.TempDir()
	if cacheDir, err := os.UserCacheDir(); err == nil {
		g.CacheDir = cacheDir
	}

	g.CacheDir = path.Join(g.CacheDir, "pulsar-routemap")

	g.NS1APIBaseURL = "https://api.nsone.net/v1"

	return g
}

// RequireAPIAccess validates that global parameters are set appropriately
// for REST API access.
func (g *CommandLineGlobals) RequireAPIAccess() error {
	var err error

	if len(g.NS1APIKey) == 0 {
		multierr.AppendInto(&err, fmt.Errorf("NS1 API key is required"))
	}

	if len(g.NS1APIBaseURL) == 0 {
		multierr.AppendInto(&err, fmt.Errorf("NS1 API base URL is required"))
	}

	return err
}
