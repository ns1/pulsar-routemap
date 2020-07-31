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
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode"

	"github.com/ns1/pulsar-routemap/pkg/lg"
	"github.com/ns1/pulsar-routemap/pkg/model"
	"go.uber.org/multierr"
)

var errUnparsableNetworkAddr = errors.New("unparsable network address")

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

// ValidateProperCIDR verifies that the IP address corresponds to the network.
func ValidateProperCIDR(ip net.IP, ipnet *net.IPNet) error {
	if !ip.Equal(ipnet.IP) {
		return fmt.Errorf("network address not properly masked")
	}

	return nil
}

// ValidateNetmaskLen verifies that the mask length does not exceeds the maximum
// allowed value.
func ValidateNetmaskLen(ipnet *net.IPNet) error {
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

// ValidateNetwork takes an string representing a network from the JSON source
// file and validates it.
// Returns the parsed IP and IPNet along with an error instance that can possibly
// be a multierr instance.
//
// If the network can not be parsed, only the error instance will contain values.
// The correctness of the network depends on the returned error having a nil value.
func ValidateNetwork(network string) (net.IP, *net.IPNet, error) {
	var errs []error

	ip, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: ", errUnparsableNetworkAddr)
	}

	if err = ValidateProperCIDR(ip, ipnet); err != nil {
		errs = append(errs, err)
	}
	if err = ValidateNetmaskLen(ipnet); err != nil {
		errs = append(errs, err)
	}

	return ip, ipnet, multierr.Combine(errs...)
}

func ValidateNetworks(nets []string, mapIdx int, summary *model.RoutemapSummary) error {
	var (
		allErrs error
		err     error
		ipnet   *net.IPNet
	)

	summary.NumNetworks += len(nets)

	for idx, n := range nets {
		lg.Tracef("visiting networks at index=%d, map segment index=%d...", idx, mapIdx)

		_, ipnet, err = ValidateNetwork(n)
		if err != nil {
			for _, e := range multierr.Errors(err) {
				// Rehydrate the packed errors so we can set the proper individual
				// error messages.
				multierr.AppendInto(&allErrs,
					fmt.Errorf("%v (for CIDR \"%s\" at index=%d, map segment index=%d)",
						e, n, idx, mapIdx))
			}
		}

		// ValidateNetwork will return a nil ipnet if the string was unparsable.
		if ipnet != nil {
			if _, bits := ipnet.Mask.Size(); bits == 32 {
				summary.NumIPv4 += 1
			} else {
				summary.NumIPv6 += 1
			}
		}
	}

	return allErrs
}

// ValidateLabels validates the set of labels of the route map.
func ValidateLabels(labels []string, mapIdx int, summary *model.RoutemapSummary) error {
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

func ValidateVersion(version int) error {
	if version < 1 {
		return fmt.Errorf("invalid or missing meta/version")
	} else if version != 1 {
		return fmt.Errorf("unsupported meta/version [value=%d]", version)
	}

	return nil
}

func startValidate(root *model.RoutemapRoot, summary *model.RoutemapSummary) error {
	if err := ValidateVersion(root.MetaVersion()); err != nil {
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

		multierr.AppendInto(&allErrs, ValidateNetworks(m.Networks, idx, summary))
		multierr.AppendInto(&allErrs, ValidateLabels(m.Labels, idx, summary))

		if lg.EnabledFor(lg.LevelDebug) && (summary.NumNetworks-lastProgressReport) > 500000 {
			numErrs := len(multierr.Errors(allErrs))
			lg.Debugf("validation progress: at map segment index %d/%d; networks visited = %d, errors = %d",
				idx, numSegments, summary.NumNetworks, numErrs)
		}
	}

	return allErrs
}
