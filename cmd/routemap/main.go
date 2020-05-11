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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ns1/pulsar-routemap/internal/config"
	"github.com/ns1/pulsar-routemap/internal/crud"
	"github.com/ns1/pulsar-routemap/internal/validate"
	"github.com/ns1/pulsar-routemap/pkg/lg"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func setupVerbosity(g *config.CommandLineGlobals) error {
	switch g.Verbosity {
	case 2:
		lg.SetLevel(lg.LevelDebug)
	case 1:
		lg.SetLevel(lg.LevelInfo)
	case 0:
		lg.SetLevel(lg.LevelWarn)
	default:
		if g.Verbosity > 2 {
			lg.SetLevel(lg.LevelTrace)
		}
	}

	return nil
}

func setupCacheDir(g *config.CommandLineGlobals) error {
	if err := os.MkdirAll(g.CacheDir, os.ModePerm); err != nil {
		return fmt.Errorf("creating cachedir: %v", err)
	}

	return nil
}

func setupAPIKey(g *config.CommandLineGlobals) error {
	if len(g.NS1APIKey) == 0 {
		g.NS1APIKey = os.Getenv("NS1_APIKEY")
	}

	return nil
}

func main() {
	globals := config.NewCommandLineGlobals()

	var rootCmd = cobra.Command{}

	rootCmd.Use = filepath.Base(os.Args[0])
	rootCmd.Version = fmt.Sprintf("%s from %s (%s)", version, date, commit)
	rootCmd.Short = fmt.Sprintf("Manage Pulsar Route Maps [%s]", rootCmd.Version)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return multierr.Combine(
			setupVerbosity(&globals),
			setupAPIKey(&globals),
			setupCacheDir(&globals))
	}

	{
		// Setup usage for flags to wrap.
		usageTmpl := strings.ReplaceAll(
			rootCmd.UsageTemplate(),
			".FlagUsages",
			".FlagUsagesWrapped 78")

		rootCmd.SetUsageTemplate(usageTmpl)
	}

	pf := rootCmd.PersistentFlags()

	pf.Bool("version", false, "Display version and build information.")

	pf.CountVarP(&globals.Verbosity, "verbose", "v",
		"Increase the verbosity of output messages. Repeatable up to 3 times.")

	pf.StringVar(&globals.CacheDir, "cachedir", globals.CacheDir,
		"Where to store cached data.")
	pf.MarkHidden("cachedir") // Not being used yet.

	pf.StringVar(&globals.NS1APIBaseURL, "api-baseurl", globals.NS1APIBaseURL,
		"Base URL for NS1 REST API. Normally the default will suffice.")
	pf.MarkHidden("api-baseurl")

	pf.StringVar(&globals.NS1APIKey, "api-key", "",
		"NS1 API key for commands that require it. You may specify either this option or "+
			"use the NS1_APIKEY environment variable. The value of this command line option "+
			"takes precedence over the environment setting.")

	validate.AddCommands(&rootCmd, &globals)
	crud.AddCommands(&rootCmd, &globals)

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		lg.Errorf("%v", err)
		os.Exit(1)
	}
}
