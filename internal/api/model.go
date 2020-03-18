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

package api

import "time"

// RoutemapPayload is a catch-all struct for responses from the routemap API.
type RoutemapPayload struct {
	Customer  int    `json:"customer"`
	MapID     int    `json:"mapid"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Created   int64  `json:"created"`
	Modified  int64  `json:"modified"`
	ErrorCode string `json:"errorCode"`
}

// CreatedString returns a formatted date-time string when the routemap was created.
func (r *RoutemapPayload) CreatedString() string {
	if r.Created < 1 {
		return ""
	}

	return time.Unix(r.Created, 0).String()
}

// ModifiedString returns a formatted date-time string when the routemap was last
// modified.
func (r *RoutemapPayload) ModifiedString() string {
	if r.Modified < 1 {
		return ""
	}

	return time.Unix(r.Modified, 0).String()
}
