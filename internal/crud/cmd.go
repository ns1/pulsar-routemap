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

package crud

import (
	"fmt"

	"github.com/ns1/pulsar-routemap/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/multierr"
)

// Options contains options (command line flags) that are common to all CRUD
// commands.
type Options struct {
	Globals *config.CommandLineGlobals

	InputFilename string
	SkipValidate  bool
	MapID         int
	Name          string
	RawOutput     bool // for list command only.
}

func (o *Options) validateName() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("name parameter is required")
	}

	return nil
}

func (o *Options) validateMapID() error {
	if o.MapID < 1 {
		return fmt.Errorf("mapid parameter is required")
	}

	return nil
}

func (o *Options) addFileFlag(flags *pflag.FlagSet) {
	flags.StringVar(&o.InputFilename, "file", "",
		"Route map file to validate. Default is STDIN.")
}

func (o *Options) addNoValidateFlag(flags *pflag.FlagSet) {
	flags.BoolVar(&o.SkipValidate, "no-validate", false,
		"Validate the route map before uploading.")
}

func (o *Options) addMapIDFlag(flags *pflag.FlagSet, desc string) {
	flags.IntVar(&o.MapID, "mapid", -1, desc)
}

func addCreateCommand(parentCmd *cobra.Command, globals *config.CommandLineGlobals) {
	opts := &Options{Globals: globals}
	sub := &cobra.Command{
		Use:   "create",
		Short: "Create a new route map with optional validation",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return multierr.Combine(
				opts.Globals.RequireAPIAccess(),
				opts.validateName(),
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCreateOrReplaceCommand(opts)
		},
	}

	flags := sub.Flags()

	opts.addFileFlag(flags)
	opts.addNoValidateFlag(flags)

	flags.StringVar(&opts.Name, "name", "",
		"Name of the route map. Required when uploading a new map.")

	parentCmd.AddCommand(sub)
}

func addListCommand(parentCmd *cobra.Command, globals *config.CommandLineGlobals) {
	opts := &Options{Globals: globals}
	sub := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available route maps",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Globals.RequireAPIAccess()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(opts)
		},
	}

	flags := sub.Flags()

	flags.BoolVar(&opts.RawOutput, "raw", false, "Output in JSON. Default is pretty-printed output.")

	parentCmd.AddCommand(sub)
}

func addReplaceCommand(parentCmd *cobra.Command, globals *config.CommandLineGlobals) {
	opts := &Options{Globals: globals}
	sub := &cobra.Command{
		Use:   "replace",
		Short: "Replace an existing route map with optional validation",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return multierr.Combine(
				opts.Globals.RequireAPIAccess(),
				opts.validateMapID(),
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCreateOrReplaceCommand(opts)
		},
	}

	flags := sub.Flags()

	opts.addFileFlag(flags)
	opts.addNoValidateFlag(flags)
	opts.addMapIDFlag(flags, "Replace an existing map identified by this ID.")

	parentCmd.AddCommand(sub)
}

func addDeleteCommand(parentCmd *cobra.Command, globals *config.CommandLineGlobals) {
	opts := &Options{Globals: globals}
	sub := &cobra.Command{
		Use:   "delete",
		Short: "Delete a route map by ID",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return multierr.Combine(
				opts.Globals.RequireAPIAccess(),
				opts.validateMapID(),
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunDeleteCommand(opts)
		},
	}

	flags := sub.Flags()

	opts.addMapIDFlag(flags, "Delete an existing map identified by this ID.")

	parentCmd.AddCommand(sub)
}

func AddCommands(parentCmd *cobra.Command, globals *config.CommandLineGlobals) {
	addCreateCommand(parentCmd, globals)
	addReplaceCommand(parentCmd, globals)
	addListCommand(parentCmd, globals)
	addDeleteCommand(parentCmd, globals)
}
