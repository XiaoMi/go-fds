package cslb

import (
	"log"
	"testing"
)

func TestRoundRobinDNSService(t *testing.T) {
	s := NewRRDNSService(
		[]string{
			"example.com",
		}, true, true)
	s.Refresh()
	log.Println(s.Nodes())
}
