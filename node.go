package orgtree

// Node представляет узел в дереве
type Node struct {
	Value    interface{} `json:"value"`
	Children []*Node     `json:"children"`
}

// NewNode создает новый узел
func NewNode(value interface{}) *Node {
	return &Node{
		Value:    value,
		Children: []*Node{},
	}
}

// AddChild добавляет дочерний узел
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
