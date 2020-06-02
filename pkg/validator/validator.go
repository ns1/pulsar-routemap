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

package validator

import (
	"fmt"
	"net"
	"strings"
	"unicode"

	"github.com/ns1/pulsar-routemap/pkg/lg"
	"github.com/ns1/pulsar-routemap/pkg/model"
	"go.uber.org/multierr"
)

func LoadAndValidate(filename string) (*model.RoutemapRoot, model.RoutemapSummary, error) {
	var (
		rmap    *model.RoutemapRoot
		summary = model.NewRoutemapSummary()
		err     error
	)

	if rmap, err = model.LoadRoutemapFileOrStdin(filename); err != nil {
		return nil, summary, err
	}

	err = startValidate(rmap, &summary)
	return rmap, summary, err
}

func isAsciiOnly(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}

	return true
}

func validateProperCIDR(ip net.IP, ipnet *net.IPNet) error {
	if !ip.Equal(ipnet.IP) {
		return fmt.Errorf("network address not properly masked")
	}

	return nil
}

func validateNetmaskLen(ipnet *net.IPNet) error {
	numOnes, numBits := ipnet.Mask.Size()
	switch {
	case numOnes == 0 && numBits == 0:
		return fmt.Errorf("invalid network mask")
	case numBits == 128 && numOnes > model.MaxNetworkBitsV6:
		return fmt.Errorf("network bits %d > %d (max)", numOnes, model.MaxNetworkBitsV6)
	case numBits == 32 && numOnes > model.MaxNetworkBitsV4:
		return fmt.Errorf("network bits %d > %d (max)", numOnes, model.MaxNetworkBitsV4)
	default:
		return nil
	}
}

func validateNetworks(nets []string, mapIdx int, summary *model.RoutemapSummary) error {
	var allErrs error

	summary.NumNetworks += len(nets)

	for idx, n := range nets {
		lg.Tracef("visiting networks at index=%d, map segment index=%d...", idx, mapIdx)

		ip, ipnet, err := net.ParseCIDR(n)
		if err != nil {
			multierr.AppendInto(&allErrs,
				fmt.Errorf("unparsable network address (for CIDR \"%s\" at index=%d, map segment index=%d)", n, idx, mapIdx))
			continue
		}

		if _, bits := ipnet.Mask.Size(); bits == 32 {
			summary.NumIPv4 += 1
		} else {
			summary.NumIPv6 += 1
		}

		var errs []error

		errs = append(errs, validateProperCIDR(ip, ipnet))
		errs = append(errs, validateNetmaskLen(ipnet))

		for _, e := range errs {
			if e != nil {
				multierr.AppendInto(&allErrs,
					fmt.Errorf("%v (for CIDR \"%s\" at index=%d, map segment index=%d)", e, n, idx, mapIdx))
			}
		}
	}

	return allErrs
}

func validateLabels(labels []string, mapIdx int, summary *model.RoutemapSummary) error {
	var allErrs error

	if len(labels) == 0 {
		return fmt.Errorf("empty labels list (at map segment index=%d)", mapIdx)
	}

	// For detecting duplicate labels in this map segment.
	uniqueLabels := map[string]bool{}

	for idx, lbl := range labels {
		lg.Tracef("visiting labels at index=%d, map segment index=%d...", idx, mapIdx)

		if len(lbl) == 0 || len(strings.TrimSpace(lbl)) == 0 {
			multierr.AppendInto(&allErrs,
				fmt.Errorf("empty or whitespace-only label (at index=%d, map segment index=%d)", idx, mapIdx))
		} else if !isAsciiOnly(lbl) {
			multierr.AppendInto(&allErrs,
				fmt.Errorf("label with non-ASCII characters (at index=%d, map segment index=%d)", idx, mapIdx))
		}

		lc := strings.ToLower(lbl)
		if _, ok := uniqueLabels[lc]; ok {
			multierr.AppendInto(&allErrs,
				fmt.Errorf("duplicate label \"%s\" (at index=%d, map segment index=%d)", lbl, idx, mapIdx))
		}

		// Add label to the unique "set".
		uniqueLabels[lc] = true

		// Update summary based on true case of the label. Because if the map passes validation then
		// the case will be uniform.
		summary.SummarizeLabel(lbl)
	}

	return allErrs
}

func validateVersion(version int) error {
	if version < 1 {
		return fmt.Errorf("invalid or missing meta/version")
	} else if version != 1 {
		return fmt.Errorf("unsupported meta/version [value=%d]", version)
	}

	return nil
}

func startValidate(root *model.RoutemapRoot, summary *model.RoutemapSummary) error {
	if err := validateVersion(root.MetaVersion()); err != nil {
		return err
	}

	if len(root.Routemap) == 0 {
		lg.Warnf("route map is empty; skipping all validation")
		summary.NumNetworks = 0
		return nil
	}

	var (
		allErrs            error
		lastProgressReport int
		numSegments        = len(root.Routemap)
	)

	for idx, m := range root.Routemap {
		lg.Tracef("visiting map segment at index %d...", idx)
		if len(m.Networks) == 0 {
			multierr.AppendInto(&allErrs, fmt.Errorf("map segment at index %d has no networks defined", idx))
			continue
		}

		multierr.AppendInto(&allErrs, validateNetworks(m.Networks, idx, summary))
		multierr.AppendInto(&allErrs, validateLabels(m.Labels, idx, summary))

		if lg.EnabledFor(lg.LevelDebug) && (summary.NumNetworks-lastProgressReport) > 500000 {
			numErrs := len(multierr.Errors(allErrs))
			lg.Debugf("validation progress: at map segment index %d/%d; networks visited = %d, errors = %d",
				idx, numSegments, summary.NumNetworks, numErrs)
		}
	}

	return allErrs
}
