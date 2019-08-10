// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package proto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	timeout = 10
)

// Response --
type Response struct {
	time time.Duration
	resp *http.Response
}

// Cost --
func (r *Response) Cost() string {
	return r.time.String()
}

// StatusCode --
func (r *Response) StatusCode() int {
	return r.resp.StatusCode
}

// Body -- used to return body bytes.
func (r *Response) Body() string {
	resp := r.resp
	if resp == nil || resp.Body == nil {
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)
}

// Json -- used to return unmarshal value from body.
func (r *Response) Json(v interface{}) error {
	if r.resp.StatusCode != 200 {
		return fmt.Errorf("%v", r.Body())
	}

	resp := r.resp
	if resp == nil || resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

// Request --
type Request struct {
	timeout int
	headers map[string]string
}

// NewRequest -- creates new request.
func NewRequest() *Request {
	return &Request{
		timeout: timeout,
		headers: make(map[string]string),
	}
}

// SetTimeout -- used to set the timeout.
func (r *Request) SetTimeout(t int) *Request {
	r.timeout = t
	return r
}

// SetHeaders -- used to set the headers pair.
func (r *Request) SetHeaders(k string, v string) *Request {
	if k == "Authorization" {
		v = fmt.Sprintf(" Bearer %s", v)
	}
	r.headers[k] = v
	return r
}

func (r *Request) doRequest(method string, url string, body string) (*Response, error) {
	response := &Response{}
	start := time.Now()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: time.Duration(r.timeout) * time.Second,
	}
	rsp, err := client.Do(req)
	response.resp = rsp
	response.time = time.Since(start)
	return response, err
}

// Post -- post request with body.
func (r *Request) Post(url string, body interface{}) (*Response, error) {
	var data string

	if body != nil {
		switch v := body.(type) {
		case string:
			data = v
		default:
			d, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			data = string(d)
		}
	}
	return r.doRequest(http.MethodPost, url, data)
}

// Get -- get request.
func (r *Request) Get(url string) (*Response, error) {
	return r.doRequest(http.MethodGet, url, "")
}
