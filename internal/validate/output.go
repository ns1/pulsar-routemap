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
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ns1/pulsar-routemap/pkg/model"
	"go.uber.org/multierr"
)

// PrettyPrintSuccess outputs to STDOUT info from the model and summary objects.
func PrettyPrintSuccess(root *model.RoutemapRoot, summary model.RoutemapSummary) {
	fmt.Printf("version: %d\n", root.MetaVersion())
	fmt.Printf("sha1: %s\n", hex.EncodeToString(root.SHA1))
	fmt.Printf("size in bytes: %d\n", root.SizeInBytes)
	summary.PrettyPrint(os.Stdout)
}

// PrettyPrintErrors outputs to STDOUT the single or multi-valued error (supported by multierr package).
// Returns a new error that indicates the number of errors contained within the input.
func PrettyPrintErrors(err error) error {
	if err == nil {
		return nil
	}

	allErrs := multierr.Errors(err)

	for _, e := range allErrs {
		fmt.Printf("! %s\n", e)
	}

	return fmt.Errorf("found %d errors", len(allErrs))
}
