
package client

import (
	"testing"
	"github.com/priyalpalkar/restclient"
    "gotest.tools/assert"
    "regexp"
    "strings"
    "fmt"
)

type ResponseInfo struct {
	Gzipped   bool              `json:"gzipped"`
	Method    string            `json:"method"`
	Origin    string            `json:"origin"`
	Useragent string            `json:"user-agent"`
	Form      map[string]string `json:"form"`
	Files     map[string]string `json:"files"`
	Headers   map[string]string `json:"headers"`
	Cookies   map[string]string `json:"cookies"`
}

func TestGetRequest(t *testing.T) {
	var headers = make(map[string]string)
	var params = make(map[string]string)
	var responseInfo ResponseInfo
	c := client.NewRestClient()
	res, err := c.Get("http://httpbin.org/get", headers, params)
	if err != nil {
		t.Error("Error is not nil", err)
	}
	res.Json(&responseInfo)
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	assert.Equal(t, responseInfo.Headers["Accept-Encoding"], "gzip")
}

func TestPostRequest(t *testing.T) {
	var headers = make(map[string]string)
	var responseInfo ResponseInfo
	var testString = "gpThzrynEC1MdenWgAILwvL2CYuNGO9RwtbH1NZJ1GE31ywFOCY%2BLCctUl86jBi8TccpdPI5ppZ%2Bgss%2BNjqGHg=="
	c := client.NewRestClient()
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": testString,
		"d": "d",
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	res, err := c.Post("http://httpbin.org/post", headers, data)
	if err != nil {
		t.Error("Error is not nil", err)
		t.FailNow()
	}
	res.Json(&responseInfo)
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	assert.Equal(t, responseInfo.Form["c"], testString)
}

func TestTextPostRequest(t *testing.T) {
	var headers = make(map[string]string)
	c := client.NewRestClient()
	data := "This is some dummy data"
	res, err := c.Post("http://httpbin.org/post", headers, data)
	if err != nil {
		t.Error("Error is not nil", err)
		t.FailNow()
	}
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	assert.Equal(t, res.Request.Header["Content-Type"][0], "text/plain")
}

func TestPostRequestNilData(t *testing.T) {
	var headers map[string]string
	var data interface{}
	c := client.NewRestClient()
	res, err := c.Post("http://httpbin.org/post", headers, data)
	if err != nil {
		t.Error("Error is not nil", err)
		t.FailNow()
	}
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	assert.Equal(t, res.Request.Header["Content-Type"][0], "text/plain")
}

func TestJSONPostRequest(t *testing.T) {
	var headers = make(map[string]string)
	var responseInfo ResponseInfo
	c := client.NewRestClient()
	data := map[string]string{
		"a": "a",
		"b": "b",
	}
	headers["Content-Type"] = "application/json"
	res, err := c.Post("http://httpbin.org/post", headers, data)
	if err != nil {
		t.Error("Error is not nil", err)
		t.FailNow()
	}
	res.Json(&responseInfo)
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
}

func TestMultipartPostRequest(t *testing.T) {
	var headers = make(map[string]string)
	var responseInfo ResponseInfo
	c := client.NewRestClient()
	data := map[string]string{
		"Hello": "World",
		"@file_name": "sample.txt",
	}
	c.RequestOptions.Multipart = true
	res, err := c.Post("http://httpbin.org/post", headers, data)
	if err != nil {
		t.Error("Error is not nil", err)
		t.FailNow()
	}
	res.Json(&responseInfo)
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	assert.Equal(t, responseInfo.Files["file_name"], "This is a test file")
}

func TestMultipartPostRequestNoFile(t *testing.T) {
	var headers = make(map[string]string)
	var responseInfo ResponseInfo
	c := client.NewRestClient()
	data := map[string]string{
		"Hello": "World",
		"file_name": "sample.txt",
		"test": "SDKMFKSMFKSMFKMSKFMSKFMKSMFKSMFKSMKFMSKMF",
	}
	c.RequestOptions.Multipart = true
	res, err := c.Post("http://httpbin.org/post", headers, data)
	if err != nil {
		t.Error("Error is not nil", err)
		t.FailNow()
	}
	res.Json(&responseInfo)
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	matched, _ := regexp.MatchString(`multipart/form-data.*?boundary=`, res.Request.Header["Content-Type"][0])
	assert.Assert(t, matched)
}

func TestGetHttpsRequest(t *testing.T) {
	var headers = make(map[string]string)
	var params = make(map[string]string)
	c := client.NewRestClient()
	c.RequestOptions.Verify = false
	res, err := c.Get("https://expired.badssl.com/", headers, params)
	if err != nil {
		t.Error("Error is not nil", err)
	}
	assert.Equal(t, res.StatusCode, 200)
}

func TestGetHttpsRequestVerify(t *testing.T) {
	var headers = make(map[string]string)
	var params = make(map[string]string)
	c := client.NewRestClient()
	c.RequestOptions.Verify = true
	_, err := c.Get("https://expired.badssl.com/", headers, params)
	assert.ErrorContains(t, err, "certificate has expired")
}

func TestRedirect(t *testing.T) {
	var headers = make(map[string]string)
	var params = make(map[string]string)
	c := client.NewRestClient()
	res, err := c.Get("http://httpbin.org/redirect/3", headers, params)
	if err == nil {
		t.Error("Must not follow location")
	}
	if !strings.Contains(err.Error(), "redirect not allowed") {
		t.Error(err)
	}
	if res.StatusCode != 302 || res.Header.Get("Location") != "/relative-redirect/2" {
		t.Error("Redirect failed: ", res.StatusCode, res.Header.Get("Location"))
	}
}

func TestRedirectFollow(t *testing.T) {
	var headers = make(map[string]string)
	var params = make(map[string]string)
	var responseInfo ResponseInfo
	c := client.NewRestClient()
	c.RequestOptions.AllowRedirects = true
	res, err := c.Get("http://httpbin.org/redirect/3", headers, params)
	if err != nil {
		t.Error("Error is not nil", err)
	}
	res.Json(&responseInfo)
	if res.StatusCode != 200 {
		t.Error("Status code is not 200", res.StatusCode)
	}
	assert.Equal(t, responseInfo.Headers["Accept-Encoding"], "gzip")
}
