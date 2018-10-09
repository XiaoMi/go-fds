package httpparser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v2tool/galaxy-fds-sdk-go/fds/httpparser"
)

func TestOriginRange(t *testing.T) {
	r := "bytes=1-3,-5"

	i := strings.Index(r, "=")
	assert.Equal(t, i, 5)
	assert.Equal(t, r[0:i], "bytes")

	rangeBody := r[i+1:]
	assert.Equal(t, rangeBody, "1-3,-5")

	ranges := strings.Split(rangeBody, ",")
	assert.Equal(t, ranges[0], "1-3")
	assert.Equal(t, ranges[1], "-5")

}

func TestRange(t *testing.T) {
	r := "bytes=-1,2-4,-5"

	ranges, err := httpparser.Range(r)
	assert.Nil(t, err)

	assert.Equal(t, len(ranges), 3)

	assert.Equal(t, ranges[0].Start, int64(0))
	assert.Equal(t, ranges[0].End, int64(1))

	assert.Equal(t, ranges[1].Start, int64(2))
	assert.Equal(t, ranges[1].End, int64(4))

	assert.Equal(t, ranges[2].Start, int64(0))
	assert.Equal(t, ranges[2].End, int64(5))
}
