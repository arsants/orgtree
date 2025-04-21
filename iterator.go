package orgtree

type Iterator interface {
	Next() *Node
}

// --- PreOrder Iterator ---

type PreOrderIterator struct {
	stack []*Node
}

func NewPreOrderIterator(root *Node) *PreOrderIterator {
	return &PreOrderIterator{stack: []*Node{root}}
}

func (it *PreOrderIterator) Next() *Node {
	if len(it.stack) == 0 {
		return nil
	}

	node := it.stack[len(it.stack)-1]
	it.stack = it.stack[:len(it.stack)-1]

	// Добавляем детей в стек в обратном порядке
	for i := len(node.Children) - 1; i >= 0; i-- {
		it.stack = append(it.stack, node.Children[i])
	}

	return node
}

// --- PostOrder Iterator ---

type PostOrderIterator struct {
	stack    []*Node
	visited  map[*Node]bool
	initDone bool
	root     *Node
}

func NewPostOrderIterator(root *Node) *PostOrderIterator {
	return &PostOrderIterator{
		stack:   []*Node{},
		visited: make(map[*Node]bool),
		root:    root,
	}
}

func (it *PostOrderIterator) Next() *Node {
	if !it.initDone {
		it.stack = append(it.stack, it.root)
		it.initDone = true
	}

	for len(it.stack) > 0 {
		node := it.stack[len(it.stack)-1]
		if len(node.Children) == 0 || it.visited[node] {
			it.stack = it.stack[:len(it.stack)-1]
			return node
		}

		// добавляем детей в стек в обратном порядке
		for i := len(node.Children) - 1; i >= 0; i-- {
			it.stack = append(it.stack, node.Children[i])
		}
		it.visited[node] = true
	}

	return nil
}

// --- BFS Iterator ---

type BFSIterator struct {
	queue []*Node
}

func NewBFSIterator(root *Node) *BFSIterator {
	return &BFSIterator{queue: []*Node{root}}
}

func (it *BFSIterator) Next() *Node {
	if len(it.queue) == 0 {
		return nil
	}

	node := it.queue[0]
	it.queue = it.queue[1:]

	for _, child := range node.Children {
		it.queue = append(it.queue, child)
	}

	return node
}
