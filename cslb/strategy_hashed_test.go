package cslb

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashed(t *testing.T) {
	s := hashedStrategy{
		hashFunc: func(key interface{}) (uint64, error) {
			return uint64(key.(int)), nil
		},
	}
	nodes := []Node{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	}
	s.SetNodes(nodes)

	next, err := s.Next()
	assert.Nil(t, next)
	assert.NotNil(t, err)

	for i := 0; i < 10; i++ {
		next, err := s.NextFor(i)
		log.Println(next)
		assert.Equal(t, nodes[i%4], next)
		assert.Nil(t, err)
	}
}
