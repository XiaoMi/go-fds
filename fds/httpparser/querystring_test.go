package httpparser_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/v2tool/galaxy-fds-sdk-go/fds/httpparser"
	"testing"
	"time"
)

func TestQueryString(t *testing.T) {
	type user struct {
		Name string `param:"name"`
		Age  int    `param:"age"`
	}
	type testOptions struct {
		user
		LastModified time.Time `param:"last-modified"`
		Other        string    `param:"-"`
		Can          bool      `param:"can"`
		OtherOption  string
		unexport     string `param:"unexport"`
	}

	option := testOptions{
		user: user{
			"John",
			12,
		},
		LastModified: time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
		Other:        "other",
		OtherOption:  "helloworld",
	}
	values, e := httpparser.QueryString(option)
	assert.Nil(t, e)
	assert.Equal(t, values.Encode(), "OtherOption=helloworld&age=12&can=false&last-modified=01+Jan+18+01%3A01+UTC&name=John")
	option.Can = true
	values, e = httpparser.QueryString(option)
	assert.Nil(t, e)
	assert.Equal(t, values.Encode(), "OtherOption=helloworld&age=12&can=true&last-modified=01+Jan+18+01%3A01+UTC&name=John")
}
