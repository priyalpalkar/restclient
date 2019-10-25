
package client

import (
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"reflect"
)

func setDefaultOptions(options map[string]interface{}, o *RequestOptions) {
	r := reflect.ValueOf(o).Elem()
	for k, v := range options {
		switch val := v.(type) {
		case int:
			reflect.Indirect(r).FieldByName(k).SetInt(int64(val))
		default:
			reflect.Indirect(r).FieldByName(k).Set(reflect.ValueOf(v))
		}
	}
}

func preProcess(c *RestClient, headers map[string]string) {
	if headers != nil {
		c.Headers = headers
	} else if c.Headers != nil {
		c.Headers = make(map[string]string)
	}
}

func setOptions(req *httpclient.HttpClient, c *RestClient) {
	req.WithOption(httpclient.OPT_FOLLOWLOCATION, c.RequestOptions.AllowRedirects)
	req.WithOption(httpclient.OPT_CONNECTTIMEOUT, c.RequestOptions.Timeout)
	req.WithOption(httpclient.OPT_TIMEOUT, c.RequestOptions.Timeout)
	if !c.RequestOptions.Verify {
		req.WithOption(httpclient.OPT_UNSAFE_TLS, true)
	}
	req.WithOption(httpclient.OPT_MAXREDIRS, c.RequestOptions.MaxRedirects)
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
