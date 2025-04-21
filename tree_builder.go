package orgtree

import (
	"github.com/google/uuid"
)

// TreeBuilder представляет построитель дерева организационной структуры
type TreeBuilder struct {
	nodes     map[uuid.UUID]*OrgNode
	edges     map[uuid.UUID]*Edge
	positions map[uuid.UUID]*Position
	relations map[uuid.UUID]*PositionNodeRelation
}

// NewTreeBuilder создает новый экземпляр TreeBuilder
func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		nodes:     make(map[uuid.UUID]*OrgNode),
		edges:     make(map[uuid.UUID]*Edge),
		positions: make(map[uuid.UUID]*Position),
		relations: make(map[uuid.UUID]*PositionNodeRelation),
	}
}

// AddNode добавляет узел в построитель
func (tb *TreeBuilder) AddNode(node *OrgNode) {
	tb.nodes[node.ID] = node
}

// AddEdge добавляет связь в построитель
func (tb *TreeBuilder) AddEdge(edge *Edge) {
	tb.edges[edge.ID] = edge
}

// AddPosition добавляет должность в построитель
func (tb *TreeBuilder) AddPosition(position *Position) {
	tb.positions[position.ID] = position
}

// AddPositionNodeRelation добавляет связь должности с узлом
func (tb *TreeBuilder) AddPositionNodeRelation(relation *PositionNodeRelation) {
	tb.relations[relation.ID] = relation
}

// BuildTree строит дерево из добавленных данных
func (tb *TreeBuilder) BuildTree() *Node {
	// Создаем корневой узел
	root := NewNode(nil)

	// Создаем карту для хранения узлов дерева
	treeNodes := make(map[uuid.UUID]*Node)

	// Сначала создаем все узлы дерева
	for _, orgNode := range tb.nodes {
		treeNode := NewNode(orgNode)
		treeNodes[orgNode.ID] = treeNode
	}

	// Затем добавляем связи между узлами
	for _, edge := range tb.edges {
		fromNode, exists := treeNodes[edge.FromNode]
		if !exists {
			continue
		}

		toNode, exists := treeNodes[edge.ToNode]
		if !exists {
			continue
		}

		fromNode.AddChild(toNode)
	}

	// Находим корневые узлы (те, у которых нет входящих связей)
	hasIncoming := make(map[uuid.UUID]bool)
	for _, edge := range tb.edges {
		hasIncoming[edge.ToNode] = true
	}

	// Добавляем корневые узлы в корневой узел дерева
	for id, treeNode := range treeNodes {
		if !hasIncoming[id] {
			root.AddChild(treeNode)
		}
	}

	return root
}
