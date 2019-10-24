package client

import (
	"github.com/ddliu/go-httpclient"
	"strings"
	"fmt"
)

const (
	TIMEOUT    = 10
	VERIFY  = false
	FILE  = false
	DEBUG_MODE = false
	LOG_LEVEL = "INFO"
	ALLOW_REDIRECTS = false
)

var defaultOptions = map[string] interface{}{
	"allowRedirects": ALLOW_REDIRECTS,
	"timeout": TIMEOUT,
	"verify": VERIFY,
	"file": FILE,
	"debugMode": DEBUG_MODE,
	"logLevel": LOG_LEVEL,
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
	var req *httpclient.HttpClient
	if url == "" {
		return nil, fmt.Errorf("URL cannot be blank")
	}
	req = httpclient.NewHttpClient()
	c.logData()
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
