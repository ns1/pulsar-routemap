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

package validate

import (
	"github.com/ns1/pulsar-routemap/internal/config"
	"github.com/spf13/cobra"
)

type Options struct {
	Globals *config.CommandLineGlobals

	InputFilename string
}

func AddCommands(parentCmd *cobra.Command, globals *config.CommandLineGlobals) {
	opts := &Options{Globals: globals}
	sub := &cobra.Command{
		Use:   "validate",
		Short: "Load and validate a route map file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunValidateCommand(opts)
		},
	}

	flags := sub.Flags()

	flags.StringVar(&opts.InputFilename, "file", "",
		"Route map file to validate. Default is STDIN.")

	parentCmd.AddCommand(sub)
}
