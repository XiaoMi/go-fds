package httpparser_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/XiaoMi/go-fds/fds/httpparser"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	type User struct {
		Name string `header:"name"`
		Age  int    `header:"age"`
	}
	type testOptions struct {
		User
		LastModified time.Time   `header:"last-modified"`
		Other        string      `header:"-"`
		Can          bool        `header:"can"`
		SomeHeader   http.Header `header:",omitempty"`
		OtherOption  string
	}

	someHeader := http.Header{}
	someHeader.Add("content-length", "10")
	option := testOptions{
		User: User{
			"John",
			12,
		},
		LastModified: time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
		Other:        "other",
	}
	headers, e := httpparser.Header(option)
	assert.Nil(t, e)

	assert.Equal(t, headers.Get("name"), "John")
	assert.Equal(t, headers.Get("age"), "12")
	assert.Equal(t, headers.Get("last-modified"), "01 Jan 18 01:01 UTC")
	assert.Empty(t, headers.Get("other"))
	assert.Equal(t, headers.Get("can"), "false")
	assert.Empty(t, headers.Get("content-length"))
	assert.Empty(t, headers.Get("other"))
	assert.Empty(t, headers.Get("OtherOption"))

	option.Can = true
	option.SomeHeader = someHeader
	option.OtherOption = "helloworld"
	headers, e = httpparser.Header(option)
	assert.Nil(t, e)
	assert.Equal(t, headers.Get("can"), "true")
	assert.Equal(t, headers.Get("content-length"), "10")
	assert.Equal(t, headers.Get("OtherOption"), "helloworld")

}
