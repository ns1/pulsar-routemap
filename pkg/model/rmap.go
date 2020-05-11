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
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	MaxNetworkBitsV4 = 26
	MaxNetworkBitsV6 = 64
)

type Routemap struct {
	Networks []string `json:"networks"`
	Labels   []string `json:"labels"`
}

type RoutemapRoot struct {
	Meta     map[string]interface{} `json:"meta"`
	Routemap []Routemap             `json:"map"`

	SHA1        []byte `json:"-"`
	SizeInBytes int    `json:"-"`
	Raw         []byte `json:"-"`
}

// LoadRoutemapFileOrStdin loads a route map from the named file (if name is not empty)
// or falls back to STDIN.
func LoadRoutemapFileOrStdin(optionalFilename string) (*RoutemapRoot, error) {
	if len(optionalFilename) == 0 {
		return LoadRoutemapFile(os.Stdin)
	} else {
		return LoadRoutemapFilename(optionalFilename)
	}
}

// LoadRoutemapFilename loads a route map from a filename.
func LoadRoutemapFilename(filename string) (*RoutemapRoot, error) {
	if source, err := os.Open(filename); err == nil {
		defer source.Close()
		return LoadRoutemapFile(source)
	} else {
		return nil, err
	}
}

// LoadRoutemapFile loads a route map from an already-opened file.
func LoadRoutemapFile(source *os.File) (*RoutemapRoot, error) {
	rmap := &RoutemapRoot{}

	r := bufio.NewReader(source)

	// Stream updates to hash function
	fileHash := sha1.New()
	teeHash := io.TeeReader(r, fileHash)

	// Save raw bytes.
	bytesBuf := bytes.Buffer{}
	teeBuf := io.TeeReader(teeHash, &bytesBuf)

	dec := json.NewDecoder(teeBuf)

	if decErr := dec.Decode(rmap); decErr != nil {
		switch decErr.(type) {
		case *json.SyntaxError:
			return nil, fmt.Errorf("parsing route map: %s at byte offset %d", decErr, decErr.(*json.SyntaxError).Offset)
		default:
			return nil, decErr
		}
	}

	rmap.SHA1 = fileHash.Sum(nil)
	rmap.Raw = bytesBuf.Bytes()
	rmap.SizeInBytes = len(rmap.Raw)

	return rmap, nil
}

func (r *RoutemapRoot) MetaVersion() int {
	if v, ok := r.Meta["version"]; ok {
		switch t := v.(type) {
		case int:
			return t
		case float64:
			return int(t)
		case float32:
			return int(t)
		case string:
			if s, err := strconv.Atoi(t); err == nil {
				return s
			}
		}
	}

	return -1
}

func (r *RoutemapRoot) SetMetaVersion(v int) {
	r.Meta["version"] = v
}

func (r *RoutemapRoot) ClearRaw() {
	r.Raw = nil
}
