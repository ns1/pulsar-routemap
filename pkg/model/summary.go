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

package model

import (
	"fmt"
	"io"
	"strings"
)

type RoutemapSummary struct {
	NumNetworks int
	NumIPv4     int
	NumIPv6     int

	LabelDistribution map[string]int
}

func NewRoutemapSummary() RoutemapSummary {
	return RoutemapSummary{LabelDistribution: map[string]int{}}
}

// SummarizeLabel adds label to LabelDistribution and returns
// the number of labels of this value seen thus far.
func (s *RoutemapSummary) SummarizeLabel(label string) int {
	v := s.LabelDistribution[label]
	v += 1
	s.LabelDistribution[label] = v
	return v
}

// PrettyPrint prints a formatted summary using the given Writer.
func (s *RoutemapSummary) PrettyPrint(w io.Writer) {
	fmt.Fprintf(w, "total networks: %d\n", s.NumNetworks)
	fmt.Fprintf(w, "v4 addresses: %d\n", s.NumIPv4)
	fmt.Fprintf(w, "v6 addresses: %d\n", s.NumIPv6)
	fmt.Fprintf(w, "total unique labels: %d\n", len(s.LabelDistribution))

	var labels []string
	for k, v := range s.LabelDistribution {
		labels = append(labels, fmt.Sprintf("%s: %d", k, v))
	}

	fmt.Fprintf(w, "label histogram: %s\n", strings.Join(labels, ", "))
}
