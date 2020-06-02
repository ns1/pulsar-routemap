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
	"encoding/hex"

	"github.com/ns1/pulsar-routemap/internal/api"
	"github.com/ns1/pulsar-routemap/internal/validate"
	"github.com/ns1/pulsar-routemap/pkg/lg"
	"github.com/ns1/pulsar-routemap/pkg/model"
	"github.com/ns1/pulsar-routemap/pkg/validator"
)

func RunCreateOrReplaceCommand(opts *Options) error {
	var (
		root *model.RoutemapRoot
		err  error
	)

	if opts.SkipValidate {
		lg.Infof("skipping validation on upload")
		if root, err = model.LoadRoutemapFileOrStdin(opts.InputFilename); err != nil {
			return err
		}
	} else {
		if root, _, err = validator.LoadAndValidate(opts.InputFilename); err != nil {
			errSummary := validate.PrettyPrintErrors(err)
			lg.Errorf("map is invalid; halting upload process")
			return errSummary
		}

		lg.Infof("map is valid; ready for upload")
	}

	lg.Infof("uploading route map: meta version = %d, sha1 = %s, size = %d",
		root.MetaVersion(),
		hex.EncodeToString(root.SHA1),
		root.SizeInBytes)

	client := api.NewClient(opts.Globals.NS1APIBaseURL, opts.Globals.NS1APIKey)

	if opts.MapID > 0 {
		lg.Infof("replacing existing mapid = %d", opts.MapID)
		return client.ReplaceRoutemap(root, opts.MapID)
	} else {
		lg.Infof("creating new map: %s", opts.Name)
		return client.CreateRoutemap(root, opts.Name)
	}
}
