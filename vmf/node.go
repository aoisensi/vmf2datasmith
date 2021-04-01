package vmf

import "strconv"

type Node struct {
	data  map[string]string
	nodes map[string][]*Node
}

func (n *Node) ID() int {
	return n.Int("id")
}

func (n *Node) CountNodes(key string) int {
	if _, ok := n.nodes[key]; !ok {
		return 0
	}
	return len(n.nodes[key])
}

func (n *Node) Nodes(key string) []*Node {
	return n.nodes[key]
}

func (n *Node) String(key string) string {
	return n.data[key]
}

func (n *Node) Int(key string) int {
	v, err := strconv.Atoi(n.String(key))
	if err != nil {
		panic(err)
	}
	return v
}

func newNode() *Node {
	return &Node{
		data:  make(map[string]string),
		nodes: make(map[string][]*Node),
	}
}

func (n *Node) addNode(key string, node *Node) {
	nodes := n.nodes[key]
	if nodes == nil {
		nodes = make([]*Node, 0, 1024)
	}
	nodes = append(nodes, node)
	n.nodes[key] = nodes
}
