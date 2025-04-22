package orgtree

import (
	"github.com/google/uuid"
)

// TreeBuilder представляет построитель дерева организационной структуры
type TreeBuilder struct {
	nodes         map[uuid.UUID]*OrgNode
	edges         []*Edge
	employeeNodes map[uuid.UUID]*EmployeeNode
}

// NewTreeBuilder создает новый экземпляр TreeBuilder
func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		nodes:         make(map[uuid.UUID]*OrgNode),
		edges:         make([]*Edge, 0),
		employeeNodes: make(map[uuid.UUID]*EmployeeNode),
	}
}

// AddNode добавляет узел в построитель
func (tb *TreeBuilder) AddNode(node interface{}) {
	switch node := node.(type) {
	case *OrgNode:
		tb.nodes[node.ID] = node
	case *EmployeeNode:
		tb.employeeNodes[node.ID] = node
	}
}

// AddEdge добавляет связь в построитель
func (tb *TreeBuilder) AddEdge(edge *Edge) {
	tb.edges = append(tb.edges, edge)
}

// Nodes возвращает карту добавленных узлов
func (tb *TreeBuilder) Nodes() map[uuid.UUID]*OrgNode {
	return tb.nodes
}

// Edges возвращает срез добавленных ребер
func (tb *TreeBuilder) Edges() []*Edge {
	return tb.edges
}

// Node возвращает узел по ID и флаг его существования
func (tb *TreeBuilder) Node(id uuid.UUID) (*OrgNode, bool) {
	node, ok := tb.nodes[id]
	return node, ok
}

// BuildTree строит дерево из добавленных данных
func (tb *TreeBuilder) BuildTree() *Node {
	// Создаем все узлы дерева
	treeNodes := make(map[uuid.UUID]*Node)
	for _, orgNode := range tb.nodes {
		treeNodes[orgNode.ID] = NewNode(orgNode)
	}

	for _, employeeNode := range tb.employeeNodes {
		treeNodes[employeeNode.ID] = NewNode(employeeNode)
	}

	// Ищем входящие связи
	hasIncoming := make(map[uuid.UUID]bool)
	for _, edge := range tb.edges {
		hasIncoming[edge.ToNode] = true
	}

	// Определяем корневые узлы
	rootNodes := []*Node{}
	for id, node := range treeNodes {
		if !hasIncoming[id] {
			rootNodes = append(rootNodes, node)
		}
	}

	// Устанавливаем дочерние связи
	for _, edge := range tb.edges {
		if fromNode, ok := treeNodes[edge.FromNode]; ok {
			if toNode, ok2 := treeNodes[edge.ToNode]; ok2 {
				fromNode.AddChild(toNode)
			}
		}
	}

	// Оборачиваем корневые узлы в общий узел-заглушку
	wrapper := NewNode(nil)
	for _, rn := range rootNodes {
		wrapper.AddChild(rn)
	}
	return wrapper
}
