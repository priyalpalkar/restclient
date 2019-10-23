package client

import (
	"github.com/ddliu/go-httpclient"
	"encoding/json"
	"strings"
	"fmt"
)

var defaultOptions = map[string] interface{}{
	"allowRedirects": false,
	"timeout": 10,
	"verify": true,
	"file": false,
}

type RestClient struct {
	httpclient.HttpClient
	RequestOptions map[string] interface{}
	data interface{}
	url string
	params interface{}
}

type RestResponse struct {
	*httpclient.Response
}

func setOptions(req *httpclient.HttpClient, c *RestClient) {
	// req.WithOption(httpclient.OPT_FOLLOWLOCATION, c.RequestOptions["allowRedirects"])
	// req.WithOption(httpclient.OPT_CONNECTTIMEOUT, c.RequestOptions["timeout"])
	// req.WithOption(httpclient.OPT_TIMEOUT, c.RequestOptions["timeout"])
	// req.WithOption(httpclient.OPT_UNSAFE_TLS, c.RequestOptions["verify"])
}

func mergeOptions(options ...map[string]interface{}) map[string]interface{} {
	rst := make(map[string]interface{})
	for _, m := range options {
		for k, v := range m {
			rst[k] = v
		}
	}
	return rst
}

func preProcess(c *RestClient, headers map[string]string) {
	if c.RequestOptions == nil {
		c.RequestOptions = defaultOptions
	} else {
		c.RequestOptions = mergeOptions(defaultOptions, c.RequestOptions)
	}
	if headers != nil {
		c.Headers = headers
	} else if c.Headers != nil {
		c.Headers = make(map[string]string)
	}
}

func postRequest(c *RestClient, req *httpclient.HttpClient) (*RestResponse, error) {
	var res *httpclient.Response
	var err error
	contentType := c.Headers["Content-Type"]
	if contentType == "application/json" {
		res, err = req.PostJson(c.url, c.data)
	} else if contentType == "application/x-www-form-urlencoded" {
		res, err = req.Post(c.url, c.data)
	} else if strings.Contains(c.Headers["Content-Type"], "multipart/form-data") {
		if _, ok := c.RequestOptions["file"]; ok {
			//Similar to curl, parameter needs to start with @ to be considered as a filepath
			res, err = req.Post(c.url, c.data)
		} else {
			res, err = req.PostMultipart(c.url, c.data)
		}
	} else {
		if c.Headers["Content-Type"] == "" {
			req.WithHeader("Content-Type", "text/plain")
		}
		body := strings.NewReader(c.data.(string))
		res, err = req.Do("POST", c.url, c.Headers, body)
	}
	response := &RestResponse{res}
	return response, err
}

func (c *RestClient) Post(url string, headers map[string]string, data interface{}) (*RestResponse, error) {
	//Handle definitions in YML files
	//Logging failures etc.
	//Handle nil values
	var req *httpclient.HttpClient
	if url == "" {
		return nil, fmt.Errorf("URL cannot be blank")
	}
	req = httpclient.NewHttpClient()
	preProcess(c, headers)
	c.data = data
	c.url = url
	req.WithHeaders(c.Headers)
	setOptions(req, c)
	response, err := postRequest(c, req)
	return response, err
}

func (c *RestClient) Get(url string, headers map[string]string, params map[string]string) (*RestResponse, error) {
	var req *httpclient.HttpClient
	if url == "" {
		return nil, fmt.Errorf("URL cannot be blank")
	}
	req = httpclient.NewHttpClient()
	preProcess(c, headers)
	c.params = params
	c.url = url
	req.WithHeaders(c.Headers)
	setOptions(req, c)
	response, err := req.Get(c.url, c.params)
	res := &RestResponse{response}
	return res, err
}

func (res *RestResponse) Text() (string, error) {
	result, err := res.ToString()
	return result, err
}

func (res *RestResponse) Content() ([]byte, error) {
	result, err := res.ReadAll()
	return result, err
}

func (res *RestResponse) Json(target interface{}) (error) {
    defer res.Body.Close()
    return json.NewDecoder(res.Body).Decode(target)
}

func (res *RestResponse) Headers() (interface{}) {
	return res.Header
}
