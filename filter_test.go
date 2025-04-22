package orgtree

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
)

// Вспомогательные функции для тестов
func createTestPositions() (ceoPos, devPos, qaPos *Position) {
	ceoPos = &Position{
		ID:      uuid.New(),
		Name:    "CEO",
		SysName: "ceo",
	}
	devPos = &Position{
		ID:      uuid.New(),
		Name:    "Developer",
		SysName: "developer",
	}
	qaPos = &Position{
		ID:      uuid.New(),
		Name:    "QA Engineer",
		SysName: "qa_engineer",
	}
	return
}

func createTestNodeTypes() (departmentType, teamType, employeeType *NodeType) {
	departmentType = &NodeType{
		ID:      uuid.New(),
		Name:    "Department",
		SysName: "department",
	}
	teamType = &NodeType{
		ID:      uuid.New(),
		Name:    "Team",
		SysName: "team",
	}
	employeeType = &NodeType{
		ID:      uuid.New(),
		Name:    "Employee",
		SysName: "employee",
	}
	return
}

func createTestTree() *Node {
	ceoPos, devPos, qaPos := createTestPositions()
	departmentType, teamType, employeeType := createTestNodeTypes()

	// Создаем тестовые узлы
	root := NewNode(&OrgNode{
		ID:       uuid.New(),
		Name:     "Engineering",
		SysName:  "engineering",
		Position: ceoPos,
		Type:     departmentType,
	})
	team1 := NewNode(&OrgNode{
		ID:       uuid.New(),
		Name:     "Backend Team",
		SysName:  "backend_team",
		Position: devPos,
		Type:     teamType,
	})
	team2 := NewNode(&OrgNode{
		ID:       uuid.New(),
		Name:     "Frontend Team",
		SysName:  "frontend_team",
		Position: devPos,
		Type:     teamType,
	})
	employee1 := NewNode(&OrgNode{
		ID:       uuid.New(),
		Name:     "John Doe",
		SysName:  "john_doe",
		Position: qaPos,
		Type:     employeeType,
	})

	// Строим дерево
	root.AddChild(team1)
	root.AddChild(team2)
	team1.AddChild(employee1)

	return root
}

// Вспомогательная функция для проверки отфильтрованного дерева
func assertFilteredTree(t *testing.T, filteredTree *Node, expectedRootName string, expectedChildrenCount int) {
	if filteredTree == nil {
		t.Fatal("Filtered tree is nil")
	}

	// Проверяем корневой узел
	if orgNode, ok := filteredTree.Value.(*OrgNode); !ok || orgNode.Name != expectedRootName {
		t.Errorf("Expected root node '%s', got '%v'", expectedRootName, filteredTree.Value)
	}

	// Проверяем количество детей
	if len(filteredTree.Children) != expectedChildrenCount {
		t.Fatalf("Expected %d children, got %d", expectedChildrenCount, len(filteredTree.Children))
	}
}

func TestFilterSubtree(t *testing.T) {
	root := createTestTree()

	// Тест 1: Фильтрация по имени (используя регулярное выражение)
	t.Run("Filter by Name", func(t *testing.T) {
		predicate := func(value interface{}) bool {
			if orgNode, ok := value.(*OrgNode); ok {
				return orgNode.Name == "Backend Team" || orgNode.Name == "Frontend Team"
			}
			return false
		}

		filteredTree := root.FilterSubtree(predicate)
		assertFilteredTree(t, filteredTree, "Engineering", 2)

		// Проверяем команды
		team1Filtered := filteredTree.Children[0]
		if orgNode, ok := team1Filtered.Value.(*OrgNode); !ok || orgNode.Name != "Backend Team" {
			t.Errorf("Expected first child 'Backend Team', got '%v'", team1Filtered.Value)
		}

		team2Filtered := filteredTree.Children[1]
		if orgNode, ok := team2Filtered.Value.(*OrgNode); !ok || orgNode.Name != "Frontend Team" {
			t.Errorf("Expected second child 'Frontend Team', got '%v'", team2Filtered.Value)
		}

		// Проверяем, что у команд нет детей
		if len(team1Filtered.Children) != 0 {
			t.Errorf("Expected no children for Backend Team, got %d", len(team1Filtered.Children))
		}
		if len(team2Filtered.Children) != 0 {
			t.Errorf("Expected no children for Frontend Team, got %d", len(team2Filtered.Children))
		}
	})

	// Тест 2: Фильтрация по должности Developer
	t.Run("Filter by Position", func(t *testing.T) {
		predicate := func(value interface{}) bool {
			if orgNode, ok := value.(*OrgNode); ok {
				return orgNode.Position != nil && orgNode.Position.Name == "Developer"
			}
			return false
		}

		filteredTree := root.FilterSubtree(predicate)
		assertFilteredTree(t, filteredTree, "Engineering", 2)

		// Проверяем команды
		team1Filtered := filteredTree.Children[0]
		if orgNode, ok := team1Filtered.Value.(*OrgNode); !ok || orgNode.Name != "Backend Team" {
			t.Errorf("Expected first child 'Backend Team', got '%v'", team1Filtered.Value)
		}

		team2Filtered := filteredTree.Children[1]
		if orgNode, ok := team2Filtered.Value.(*OrgNode); !ok || orgNode.Name != "Frontend Team" {
			t.Errorf("Expected second child 'Frontend Team', got '%v'", team2Filtered.Value)
		}

		// Проверяем, что у команд нет детей
		if len(team1Filtered.Children) != 0 {
			t.Errorf("Expected no children for Backend Team, got %d", len(team1Filtered.Children))
		}
		if len(team2Filtered.Children) != 0 {
			t.Errorf("Expected no children for Frontend Team, got %d", len(team2Filtered.Children))
		}
	})

	// Тест 3: Фильтрация по типу Team
	t.Run("Filter by Type", func(t *testing.T) {
		predicate := func(value interface{}) bool {
			if orgNode, ok := value.(*OrgNode); ok {
				return orgNode.Type != nil && orgNode.Type.Name == "Team"
			}
			return false
		}

		filteredTree := root.FilterSubtree(predicate)
		assertFilteredTree(t, filteredTree, "Engineering", 2)

		// Проверяем команды
		team1Filtered := filteredTree.Children[0]
		if orgNode, ok := team1Filtered.Value.(*OrgNode); !ok || orgNode.Name != "Backend Team" {
			t.Errorf("Expected first child 'Backend Team', got '%v'", team1Filtered.Value)
		}

		team2Filtered := filteredTree.Children[1]
		if orgNode, ok := team2Filtered.Value.(*OrgNode); !ok || orgNode.Name != "Frontend Team" {
			t.Errorf("Expected second child 'Frontend Team', got '%v'", team2Filtered.Value)
		}

		// Проверяем, что у команд нет детей
		if len(team1Filtered.Children) != 0 {
			t.Errorf("Expected no children for Backend Team, got %d", len(team1Filtered.Children))
		}
		if len(team2Filtered.Children) != 0 {
			t.Errorf("Expected no children for Frontend Team, got %d", len(team2Filtered.Children))
		}
	})

	// Тест 4: Фильтрация по регулярному выражению (для строковых значений)
	t.Run("Filter by Regex", func(t *testing.T) {
		// Создаем простое дерево со строковыми значениями для тестирования регулярных выражений
		stringRoot := NewNode("root")
		stringChild1 := NewNode("test1")
		stringChild2 := NewNode("test2")
		stringGrandchild1 := NewNode("test3")
		stringGrandchild2 := NewNode("other")

		stringRoot.AddChild(stringChild1)
		stringRoot.AddChild(stringChild2)
		stringChild1.AddChild(stringGrandchild1)
		stringChild2.AddChild(stringGrandchild2)

		pattern := regexp.MustCompile("test")
		filteredTree, err := stringRoot.FilterSubtreeByRegex(pattern.String())
		if err != nil {
			t.Fatalf("Failed to filter subtree: %v", err)
		}

		// Проверяем результат
		if filteredTree == nil {
			t.Fatal("Filtered tree is nil")
		}

		// Проверяем корневой узел
		if filteredTree.Value != "root" {
			t.Errorf("Expected root value 'root', got '%v'", filteredTree.Value)
		}

		// Проверяем детей корневого узла
		if len(filteredTree.Children) != 2 {
			t.Errorf("Expected 2 children, got %d", len(filteredTree.Children))
		}

		// Проверяем первого ребенка
		child1Filtered := filteredTree.Children[0]
		if child1Filtered.Value != "test1" {
			t.Errorf("Expected child1 value 'test1', got '%v'", child1Filtered.Value)
		}

		// Проверяем второго ребенка
		child2Filtered := filteredTree.Children[1]
		if child2Filtered.Value != "test2" {
			t.Errorf("Expected child2 value 'test2', got '%v'", child2Filtered.Value)
		}

		// Проверяем внука
		if len(child1Filtered.Children) != 1 {
			t.Errorf("Expected 1 grandchild for child1, got %d", len(child1Filtered.Children))
		}
		grandchild1Filtered := child1Filtered.Children[0]
		if grandchild1Filtered.Value != "test3" {
			t.Errorf("Expected grandchild1 value 'test3', got '%v'", grandchild1Filtered.Value)
		}
	})
}
