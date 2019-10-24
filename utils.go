package client

import (
	"github.com/netSkope/go-kestrel/pkg/log"
	"encoding/json"
	"github.com/ddliu/go-httpclient"
)

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

func setOptions(req *httpclient.HttpClient, c *RestClient) {
	req.WithOption(httpclient.OPT_FOLLOWLOCATION, c.RequestOptions["allowRedirects"])
	req.WithOption(httpclient.OPT_CONNECTTIMEOUT, c.RequestOptions["timeout"])
	req.WithOption(httpclient.OPT_TIMEOUT, c.RequestOptions["timeout"])
	req.WithOption(httpclient.OPT_UNSAFE_TLS, c.RequestOptions["verify"])
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

func (c *RestClient) logData() () {
	logger := log.NewLogger("main-logger")
	logger.Debug("This is a test log message")
}
