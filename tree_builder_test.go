package orgtree

import (
	"testing"

	"github.com/google/uuid"
)

func TestTreeBuilder(t *testing.T) {
	// Создаем построитель дерева
	builder := NewTreeBuilder()

	// Создаем тестовые данные
	nodeType := &NodeType{
		ID:      uuid.New(),
		Name:    "Отдел",
		SysName: "department",
	}

	node1 := &OrgNode{
		ID:      uuid.New(),
		Name:    "Главный офис",
		SysName: "main_office",
		TypeID:  nodeType.ID,
		Type:    nodeType,
	}

	node2 := &OrgNode{
		ID:      uuid.New(),
		Name:    "IT отдел",
		SysName: "it_department",
		TypeID:  nodeType.ID,
		Type:    nodeType,
	}

	edge := &Edge{
		ID:       uuid.New(),
		Name:     "Подчинение",
		SysName:  "subordination",
		FromNode: node1.ID,
		ToNode:   node2.ID,
	}

	// Добавляем данные в построитель
	builder.AddNode(node1)
	builder.AddNode(node2)
	builder.AddEdge(edge)

	// Строим дерево
	tree := builder.BuildTree()

	// Проверяем структуру дерева с помощью итератора
	iterator := NewPreOrderIterator(tree)

	// Пропускаем корневой узел (он пустой)
	iterator.Next()

	// Проверяем первый узел (Главный офис)
	node := iterator.Next()
	if node == nil {
		t.Fatal("Ожидался узел 'Главный офис', получен nil")
	}

	if orgNode, ok := node.Value.(*OrgNode); ok {
		if orgNode.Name != "Главный офис" {
			t.Errorf("Ожидалось имя 'Главный офис', получено '%s'", orgNode.Name)
		}
	} else {
		t.Error("Значение узла не является *OrgNode")
	}

	// Проверяем второй узел (IT отдел)
	node = iterator.Next()
	if node == nil {
		t.Fatal("Ожидался узел 'IT отдел', получен nil")
	}

	if orgNode, ok := node.Value.(*OrgNode); ok {
		if orgNode.Name != "IT отдел" {
			t.Errorf("Ожидалось имя 'IT отдел', получено '%s'", orgNode.Name)
		}
	} else {
		t.Error("Значение узла не является *OrgNode")
	}

	// Проверяем, что больше узлов нет
	if iterator.Next() != nil {
		t.Error("Ожидалось окончание итерации")
	}
}
