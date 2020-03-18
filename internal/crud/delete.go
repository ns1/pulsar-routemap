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
	"github.com/ns1/pulsar-routemap/internal/api"
	"github.com/ns1/pulsar-routemap/pkg/lg"
)

func RunDeleteCommand(opts *Options) error {
	client := api.NewClient(opts.Globals.NS1APIBaseURL, opts.Globals.NS1APIKey)

	if err := client.DeleteRoutemap(opts.MapID); err != nil {
		return err
	} else {
		lg.Printf("deleted routemap %d", opts.MapID)
	}

	return nil
}
