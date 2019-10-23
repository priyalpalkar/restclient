package client

import (
	"github.com/ddliu/go-httpclient"
	"fmt"
)

var defaultOptions = map[string] interface{}{
	"allowRedirects": false,
	"timeout": 10,
	"verify": true,
}



type RestClient struct {
	*httpclient.HttpClient
	options map[string] interface{}
	data interface{}
}

type RestResponse struct {
	*httpclient.Response
	// Json()
	// Text()
	// Content()
}

func setOptions(req *httpclient.HttpClient, c *RestClient) {
	req.WithOption(httpclient.OPT_FOLLOWLOCATION, c.options["allowRedirects"])
	req.WithOption(httpclient.OPT_CONNECTTIMEOUT, c.options["timeout"])
	req.WithOption(httpclient.OPT_UNSAFE_TLS, c.options["verify"])
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

func (c *RestClient) Post(url string, headers map[string]string, data string) (*RestResponse, error) {
	//Handle definitions in YML files
	//Add the text, content and Json() methods to the response structure
	//Logging failures etc.
	//Handle nil values
	var req *httpclient.HttpClient
	if c.options == nil {
		c.options = defaultOptions
	} else {
		c.options = mergeOptions(defaultOptions, c.options)
	}
	// if headers != nil {
	// 	c.Headers = headers
	// } else {
	// 	c.Headers = make(map[string]string)
	// }
	// req.WithHeaders(c.Headers)
	setOptions(req, c)
	fmt.Println(c.Headers, c.options)
	response, err := req.Post(url, nil)
	res := &RestResponse{response}
	return res, err
}

func (c *RestClient) Get(url string, headers map[string]string) (*RestResponse, error) {
	var req *httpclient.HttpClient
	// if headers != nil {
	// 	res = httpclient.WithHeaders(headers)
	// }
	// res.WithOption(httpclient.OPT_FOLLOWLOCATION, options.allowRedirects)
	// res.WithOption(httpclient.OPT_CONNECTTIMEOUT, option.timeout)
	// res.WithOption(httpclient.OPT_UNSAFE_TLS, option.verify)
	response, err := req.Get(url, nil)
	res := &RestResponse{response}
	return res, err
}

func (res RestResponse) Json() {

}

func ResponseBody(response *httpclient.Response) string {
	body, _  := response.ToString()
	return body
}
