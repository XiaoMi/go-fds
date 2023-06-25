package cslb

// Strategy controls how the nodes are chosen
// This type should be thread safe
type Strategy interface {
	// SetNodes update saved nodes
	SetNodes(nodes []Node)
	// Next returns a node address
	Next() (Node, error)
	// NextFor returns a node address assigned to request
	NextFor(interface{}) (Node, error)
}
