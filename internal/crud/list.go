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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ns1/pulsar-routemap/internal/api"
)

func RunListCommand(opts *Options) error {
	client := api.NewClient(opts.Globals.NS1APIBaseURL, opts.Globals.NS1APIKey)

	if body, err := client.ListRoutemaps(); err != nil {
		return err
	} else if opts.RawOutput {
		fmt.Printf("%s\n", body)
		return nil
	} else {
		return prettyPrint(body)
	}
}

func prettyPrint(body []byte) error {
	var rmaps []api.RoutemapPayload
	if err := json.Unmarshal(body, &rmaps); err != nil {
		return fmt.Errorf("parsing routemap payload: %v", err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)
	defer tw.Flush()

	pp := func(values ...string) {
		line := strings.Join(values, "\t")
		fmt.Fprintf(tw, "%s\t\n", line)
	}

	pp("id", "name", "created", "modified", "status")
	pp("--", "----", "-------", "--------", "------")

	for _, m := range rmaps {
		pp(strconv.Itoa(m.MapID), m.Name, m.CreatedString(), m.ModifiedString(), m.Status)
	}

	return nil
}
