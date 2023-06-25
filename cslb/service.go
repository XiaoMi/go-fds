package cslb

// Service represents a group of nodes provide identical functionality, cluster as a logical service
// This type should be thread safe
type Service interface {
	// Nodes returns a new slice of available node
	Nodes() []Node
	// Refresh updates nodes
	Refresh()
	// NodeFailedCallbackFunc returns a callback function which will be triggered in another go routine when certain
	// node exiled by LoadBalancer
	NodeFailedCallbackFunc() func(node Node)
}
