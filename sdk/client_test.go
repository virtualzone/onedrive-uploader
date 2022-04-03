package sdk

import (
	"sort"
	"strings"
	"testing"
)

func TestIsHTTPStatusOK(t *testing.T) {
	checkTestBool(t, false, IsHTTPStatusOK(199))
	checkTestBool(t, true, IsHTTPStatusOK(200))
	checkTestBool(t, true, IsHTTPStatusOK(201))
	checkTestBool(t, true, IsHTTPStatusOK(299))
	checkTestBool(t, false, IsHTTPStatusOK(300))
}

func TestBuildURIParams0(t *testing.T) {
	params := HTTPRequestParams{}
	c := &Client{}
	s := c.buildURIParams(params)
	checkTestString(t, "", s)
}

func TestBuildURIParams1(t *testing.T) {
	params := HTTPRequestParams{
		"p1": "test1",
	}
	c := &Client{}
	s := c.buildURIParams(params)
	checkTestString(t, "p1=test1", s)
}

func TestBuildURIParams3(t *testing.T) {
	params := HTTPRequestParams{
		"p1": "test1",
		"p2": "test2",
		"p3": "test3",
	}
	c := &Client{}
	s := c.buildURIParams(params)
	tokens := strings.Split(s, "&")
	sort.Strings(tokens)
	checkTestString(t, "p1=test1", tokens[0])
	checkTestString(t, "p2=test2", tokens[1])
	checkTestString(t, "p3=test3", tokens[2])
}

func TestBuildURIParamsSpecialChars(t *testing.T) {
	params := HTTPRequestParams{
		"p1": "test 1",
		"p2": "test&2",
		"p3": "testü3",
	}
	c := &Client{}
	s := c.buildURIParams(params)
	tokens := strings.Split(s, "&")
	sort.Strings(tokens)
	checkTestString(t, "p1=test+1", tokens[0])
	checkTestString(t, "p2=test%262", tokens[1])
	checkTestString(t, "p3=test%C3%BC3", tokens[2])
}

func TestBuildURINoParams(t *testing.T) {
	params := HTTPRequestParams{}
	c := &Client{}
	s := c.buildURI("http://test", params)
	checkTestString(t, "http://test", s)
}

func TestBuildURIWithParams(t *testing.T) {
	params := HTTPRequestParams{
		"p1": "test1",
	}
	c := &Client{}
	s := c.buildURI("http://test", params)
	checkTestString(t, "http://test?p1=test1", s)
}
