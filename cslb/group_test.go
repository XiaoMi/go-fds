package cslb

import (
	"net"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	var actual []Node
	example := []Node{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	}
	g := NewGroup(NodeCountUnlimited)

	g.Set(example)
	actual = g.Get()
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].String() < actual[j].String()
	})
	assert.Equal(t, example, actual)

	g.Exile(&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)})
	actual = g.Get()
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].String() < actual[j].String()
	})
	assert.Equal(t, []Node{
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	}, actual)
}
