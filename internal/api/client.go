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

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ns1/pulsar-routemap/pkg/model"
)

// Client defines the interactions with NS1's REST API.
type Client interface {
	// ListRoutemaps returns all routemaps for the customer. Result is raw JSON.
	ListRoutemaps() ([]byte, error)

	// CreateRoutemap creates a new routemap of the given name.
	CreateRoutemap(root *model.RoutemapRoot, name string) error

	// ReplaceRoutemap replaces the existing routemap given by mapid.
	ReplaceRoutemap(root *model.RoutemapRoot, mapid int) error

	// DeleteRoutemap deletes an existing routemap given by mapid. It's an error
	// to delete a routemap that does not exist.
	DeleteRoutemap(mapid int) error
}

type httpClient struct {
	apiKey  string
	baseURL string

	inst *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL string, apiKey string) Client {
	return &httpClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		inst:    &http.Client{},
	}
}

func (c *httpClient) ListRoutemaps() ([]byte, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	if req, err = c.newRequest("GET", "/pulsar/routemaps"); err != nil {
		return nil, fmt.Errorf("creating API request: %v", err)
	}

	if resp, err = c.inst.Do(req); err != nil {
		return nil, fmt.Errorf("issuing API request: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, notOKToError(resp)
	}

	reader := resp.Body
	defer reader.Close()

	var body []byte
	if body, err = ioutil.ReadAll(reader); err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	} else {
		return body, nil
	}
}

func (c *httpClient) CreateRoutemap(root *model.RoutemapRoot, name string) error {
	var (
		req *http.Request
		err error
	)

	if req, err = c.newRequest("GET", "/pulsar/routemaps/create"); err != nil {
		return fmt.Errorf("creating API request: %v", err)
	}

	q := url.Values{}
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	return c.uploadMap(root, req)
}

func (c *httpClient) ReplaceRoutemap(root *model.RoutemapRoot, mapid int) error {
	var (
		req *http.Request
		err error
	)

	if req, err = c.newRequest("GET", fmt.Sprintf("/pulsar/routemaps/%d/replace", mapid)); err != nil {
		return fmt.Errorf("creating API request: %v", err)
	}

	return c.uploadMap(root, req)
}

func (c *httpClient) DeleteRoutemap(mapid int) error {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	if req, err = c.newRequest("DELETE", fmt.Sprintf("/pulsar/routemaps/%d", mapid)); err != nil {
		return fmt.Errorf("creating API request: %v", err)
	}

	if resp, err = c.inst.Do(req); err != nil {
		return fmt.Errorf("deleting routemap: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deleting routemap: %v", notOKToError(resp))
	}

	return nil
}

func (c *httpClient) fetchUploadURL(req *http.Request) (string, error) {
	var (
		resp *http.Response
		err  error
	)

	if resp, err = c.inst.Do(req); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", notOKToError(resp)
	}

	body := resp.Body
	defer body.Close()

	var bytes []byte
	if bytes, err = ioutil.ReadAll(body); err != nil {
		return "", fmt.Errorf("reading response body: %v", err)
	}

	return string(bytes), nil
}

func (c *httpClient) uploadMap(root *model.RoutemapRoot, startUploadReq *http.Request) error {
	var (
		uploadURL string
		err       error
	)

	if uploadURL, err = c.fetchUploadURL(startUploadReq); err != nil {
		return fmt.Errorf("starting map upload: %v", err)
	}

	body := bytes.NewReader(root.Raw)

	// Note: we aren't issuing an API request here; it's a fully-qualified URL.
	var req *http.Request
	if req, err = http.NewRequest("PUT", uploadURL, body); err != nil {
		return fmt.Errorf("creating API request: %v", err)
	}

	var resp *http.Response
	if resp, err = c.inst.Do(req); err != nil {
		return fmt.Errorf("uploading routemap: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("transferring routemap: %v", notOKToError(resp))
	}

	return nil
}

func (c *httpClient) newRequest(method string, requestURI string) (*http.Request, error) {
	return c.newRequestWithBody(method, requestURI, nil)
}

func (c *httpClient) newRequestWithBody(method string, requestURI string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(requestURI, "/"))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-NSONE-Key", c.apiKey)

	return req, nil
}

func notOKToError(resp *http.Response) error {
	statusCodeErr := fmt.Errorf("unexpected status: %s", resp.Status)

	body := resp.Body
	if body == nil {
		return statusCodeErr
	}

	defer body.Close()

	if bytes, err := ioutil.ReadAll(body); err != nil {
		return statusCodeErr
	} else {
		return fmt.Errorf("unexpected status: (%d) %s", resp.StatusCode, bytes)
	}
}
