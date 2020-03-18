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
	"net"
	"testing"

	"github.com/ns1/pulsar-routemap/pkg/model"
	"github.com/stretchr/testify/assert"
)

func Test_validateProperCIDR(t *testing.T) {
	fixtures := []struct {
		addr  string
		valid bool
	}{
		{addr: "192.168.22.22/24", valid: false},
		{addr: "192.168.22.22/32", valid: true},
		{addr: "192.168.22.0/24", valid: true},
		{addr: "2001:db8:1234::/48", valid: true},
		{addr: "2001:db8:85a3:0:0:8a2e:370:7334/56", valid: false},
		{addr: "2001:db8:85a3:0:0:8a2e:370:7334/128", valid: true},
	}

	for _, fx := range fixtures {
		ip, ipnet, err := net.ParseCIDR(fx.addr)
		assert.NoError(t, err, fx.addr)

		err = validateProperCIDR(ip, ipnet)
		if fx.valid {
			assert.NoError(t, err, fx.addr)
		} else {
			assert.Error(t, err, fx.addr)
		}
	}
}

func Test_validateNetmaskLen(t *testing.T) {
	fixtures := []struct {
		addr  string
		valid bool
	}{
		{addr: "192.168.22.0/27", valid: false},
		{addr: "192.168.22.0/26", valid: true},
		{addr: "2001:db8:1234::/56", valid: true},
		{addr: "2001:db8:1234::/64", valid: true},
		{addr: "2001:db8:1234::/96", valid: false},
	}

	for _, fx := range fixtures {
		_, ipnet, err := net.ParseCIDR(fx.addr)
		assert.NoError(t, err, fx.addr)

		err = validateNetmaskLen(ipnet)
		if fx.valid {
			assert.NoError(t, err, fx.addr)
		} else {
			assert.Error(t, err, fx.addr)
		}
	}
}

func Test_validateLabels(t *testing.T) {
	fixtures := []struct {
		labels []string
		valid  bool
	}{
		{labels: []string{""}, valid: false},
		{labels: []string{"    "}, valid: false},
		{labels: []string{"\n"}, valid: false},
		{labels: []string{"\x98"}, valid: false},
		{labels: []string{"bags", "time"}, valid: true},
	}

	summary := model.NewRoutemapSummary()

	for _, fx := range fixtures {
		err := validateLabels(fx.labels, 1, &summary)
		if fx.valid {
			assert.NoError(t, err, fx.labels)
		} else {
			assert.Error(t, err, fx.labels)
		}
	}
}
