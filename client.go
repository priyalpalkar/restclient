package client

import (
	"github.com/ddliu/go-httpclient"
	"strings"
	"fmt"
	"net/http"
)

const (
	TIMEOUT    = 10
	VERIFY  = false
	ALLOW_REDIRECTS = false
	MULTIPART = false
	MAX_REDIRECTS = 10
	CHUNKED = false
)

var defaultOptions = map[string] interface{}{
	"AllowRedirects": ALLOW_REDIRECTS,
	"Timeout": TIMEOUT,
	"Verify": VERIFY,
	"Multipart": MULTIPART,
	"MaxRedirects": MAX_REDIRECTS,
	"Chunked": CHUNKED,
}

type RequestOptions struct {
	AllowRedirects bool
	Timeout int
	Verify bool
	Multipart bool
	MaxRedirects int
	Chunked bool
}

type RestClient struct {
	httpclient.HttpClient
	RequestOptions
	data interface{}
	url string
	params interface{}
}

type RestResponse struct {
	*httpclient.Response
}

func NewRestClient() *RestClient {
	c := new(RestClient)
	setDefaultOptions(defaultOptions, &c.RequestOptions)	
	return c
}

func postRequest(c *RestClient, req *httpclient.HttpClient) (*RestResponse, error) {
	var res *httpclient.Response
	var err error
	contentType := c.Headers["Content-Type"]
	if contentType == "application/json" {
		res, err = req.PostJson(c.url, c.data)
	} else if contentType == "application/x-www-form-urlencoded" {
		res, err = req.Post(c.url, c.data)
	} else if c.RequestOptions.Multipart {
		res, err = req.PostMultipart(c.url, c.data)
	} else {
		//http client blows up in the parsing if the data cannot be converted to URL/form values
		//Handling all these cases here
		if _, ok := c.Headers["Content-Type"]; !ok {
			req.WithHeader("Content-Type", "text/plain")
		}
		res_data, ok := c.data.(string)
		if !ok {
			if c.data == nil {
				c.data = ""
			} else {
				return nil, fmt.Errorf("Payload conversion to string failed")
			}
		}
		body := strings.NewReader(res_data)
		res, err = req.Do("POST", c.url, c.Headers, body)
	}
	response := &RestResponse{res}
	return response, err
}

func (c *RestClient) Post(url string, headers map[string]string, data interface{}) (*RestResponse, error) {
	//TODO Add logging, support constructing requests from YAML/JSON files
	var req *httpclient.HttpClient
	if url == "" {
		return nil, fmt.Errorf("URL cannot be blank")
	}
	req = httpclient.NewHttpClient()
	preProcess(c, headers)
	c.data = data
	c.url = url
	if c.RequestOptions.Chunked {
		c.Headers["Content-Length"]	= "-1"
	}
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
