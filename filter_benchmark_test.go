package orgtree

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkFilterLargeTree(b *testing.B) {
	// Создаем построитель дерева
	builder := NewTreeBuilder()

	// Создаем типы узлов
	departmentType := &NodeType{
		ID:      uuid.New(),
		Name:    "Отдел",
		SysName: "department",
	}

	teamType := &NodeType{
		ID:      uuid.New(),
		Name:    "Команда",
		SysName: "team",
	}

	employeeType := &NodeType{
		ID:      uuid.New(),
		Name:    "Сотрудник",
		SysName: "employee",
	}

	// Создаем должности
	managerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Руководитель",
		SysName: "manager",
	}

	developerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Разработчик",
		SysName: "developer",
	}

	// Создаем корневой узел
	root := &OrgNode{
		ID:       uuid.New(),
		Name:     "Главный офис",
		SysName:  "main_office",
		Type:     departmentType,
		Position: managerPosition,
	}
	builder.AddNode(root)

	// Создаем большое дерево с 1_000_000 узлов
	// Структура: корень -> отделы -> команды -> сотрудники
	// Примерно 100 отделов, по 100 команд в каждом отделе, по 100 сотрудников в каждой команде
	const (
		numDepartments = 10
		numTeams       = 10
		numEmployees   = 10
	)

	departmentNodes := make([]*OrgNode, numDepartments)

	// Создаем отделы
	for i := 0; i < numDepartments; i++ {
		dept := &OrgNode{
			ID:       uuid.New(),
			Name:     "Отдел " + string(rune('A'+i)),
			SysName:  "department_" + string(rune('a'+i)),
			Type:     departmentType,
			Position: managerPosition,
		}
		builder.AddNode(dept)
		builder.AddEdge(&Edge{
			FromNode: root.ID,
			ToNode:   dept.ID,
		})
		departmentNodes[i] = dept
	}

	// Создаем команды для каждого отдела
	for _, dept := range departmentNodes {
		for i := 0; i < numTeams; i++ {
			team := &OrgNode{
				ID:       uuid.New(),
				Name:     "Команда " + string(rune('A'+i)),
				SysName:  "team_" + string(rune('a'+i)),
				Type:     teamType,
				Position: managerPosition,
			}
			builder.AddNode(team)
			builder.AddEdge(&Edge{
				FromNode: dept.ID,
				ToNode:   team.ID,
			})

			// Создаем сотрудников для каждой команды
			for j := 0; j < numEmployees; j++ {
				employee := &OrgNode{
					ID:       uuid.New(),
					Name:     "Сотрудник " + string(rune('A'+j)),
					SysName:  "employee_" + string(rune('a'+j)),
					Type:     employeeType,
					Position: developerPosition,
				}
				builder.AddNode(employee)
				builder.AddEdge(&Edge{
					FromNode: team.ID,
					ToNode:   employee.ID,
				})
			}
		}
	}

	// Строим дерево
	tree := builder.BuildTree()

	// Создаем предикат для фильтрации сотрудников
	employeePredicate := func(value interface{}) bool {
		if orgNode, ok := value.(*OrgNode); ok {
			return orgNode.Type != nil && orgNode.Type.SysName == "employee"
		}
		return false
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filteredTree := tree.FilterSubtree(employeePredicate)
		if filteredTree == nil {
			b.Fatal("Filtered tree is nil")
		}
	}
}

func BenchmarkFilterDifferentSizes(b *testing.B) {
	sizes := []struct {
		name        string
		departments int
		teams       int
		employees   int
		totalNodes  int
	}{
		{"Small (100)", 5, 5, 5, 100},
		{"Medium (1K)", 5, 5, 50, 1000},
		{"Large (10K)", 5, 50, 50, 10000},
		{"Huge (100K)", 50, 50, 50, 100000},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			// Создаем построитель дерева
			builder := NewTreeBuilder()

			// Создаем типы узлов
			departmentType := &NodeType{
				ID:      uuid.New(),
				Name:    "Отдел",
				SysName: "department",
			}

			teamType := &NodeType{
				ID:      uuid.New(),
				Name:    "Команда",
				SysName: "team",
			}

			employeeType := &NodeType{
				ID:      uuid.New(),
				Name:    "Сотрудник",
				SysName: "employee",
			}

			// Создаем должности
			managerPosition := &Position{
				ID:      uuid.New(),
				Name:    "Руководитель",
				SysName: "manager",
			}

			developerPosition := &Position{
				ID:      uuid.New(),
				Name:    "Разработчик",
				SysName: "developer",
			}

			// Создаем корневой узел
			root := &OrgNode{
				ID:       uuid.New(),
				Name:     "Главный офис",
				SysName:  "main_office",
				Type:     departmentType,
				Position: managerPosition,
			}
			builder.AddNode(root)

			departmentNodes := make([]*OrgNode, size.departments)

			// Создаем отделы
			for i := 0; i < size.departments; i++ {
				dept := &OrgNode{
					ID:       uuid.New(),
					Name:     fmt.Sprintf("Отдел %d", i),
					SysName:  fmt.Sprintf("department_%d", i),
					Type:     departmentType,
					Position: managerPosition,
				}
				builder.AddNode(dept)
				builder.AddEdge(&Edge{
					FromNode: root.ID,
					ToNode:   dept.ID,
				})
				departmentNodes[i] = dept
			}

			// Создаем команды для каждого отдела
			for _, dept := range departmentNodes {
				for i := 0; i < size.teams; i++ {
					team := &OrgNode{
						ID:       uuid.New(),
						Name:     fmt.Sprintf("Команда %d", i),
						SysName:  fmt.Sprintf("team_%d", i),
						Type:     teamType,
						Position: managerPosition,
					}
					builder.AddNode(team)
					builder.AddEdge(&Edge{
						FromNode: dept.ID,
						ToNode:   team.ID,
					})

					// Создаем сотрудников для каждой команды
					for j := 0; j < size.employees; j++ {
						employee := &OrgNode{
							ID:       uuid.New(),
							Name:     fmt.Sprintf("Сотрудник %d", j),
							SysName:  fmt.Sprintf("employee_%d", j),
							Type:     employeeType,
							Position: developerPosition,
						}
						builder.AddNode(employee)
						builder.AddEdge(&Edge{
							FromNode: team.ID,
							ToNode:   employee.ID,
						})
					}
				}
			}

			// Строим дерево
			tree := builder.BuildTree()

			// Создаем предикат для фильтрации сотрудников
			employeePredicate := func(value interface{}) bool {
				if orgNode, ok := value.(*OrgNode); ok {
					return orgNode.Type != nil && orgNode.Type.SysName == "employee"
				}
				return false
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				filteredTree := tree.FilterSubtree(employeePredicate)
				if filteredTree == nil {
					b.Fatal("Filtered tree is nil")
				}
			}
		})
	}
}

func BenchmarkFilterDifferentTypes(b *testing.B) {
	// Создаем построитель дерева
	builder := NewTreeBuilder()

	// Создаем типы узлов
	departmentType := &NodeType{
		ID:      uuid.New(),
		Name:    "Отдел",
		SysName: "department",
	}

	teamType := &NodeType{
		ID:      uuid.New(),
		Name:    "Команда",
		SysName: "team",
	}

	employeeType := &NodeType{
		ID:      uuid.New(),
		Name:    "Сотрудник",
		SysName: "employee",
	}

	// Создаем должности
	managerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Руководитель",
		SysName: "manager",
	}

	// Создаем корневой узел
	root := &OrgNode{
		ID:       uuid.New(),
		Name:     "Главный офис",
		SysName:  "main_office",
		Type:     departmentType,
		Position: managerPosition,
	}
	builder.AddNode(root)

	// Создаем дерево среднего размера (10K узлов)
	const (
		numDepartments = 100
		numTeams       = 100
		numEmployees   = 10
	)

	departmentNodes := make([]*OrgNode, numDepartments)

	// Создаем отделы
	for i := 0; i < numDepartments; i++ {
		dept := &OrgNode{
			ID:       uuid.New(),
			Name:     fmt.Sprintf("Отдел %d", i),
			SysName:  fmt.Sprintf("department_%d", i),
			Type:     departmentType,
			Position: managerPosition,
		}
		builder.AddNode(dept)
		builder.AddEdge(&Edge{
			FromNode: root.ID,
			ToNode:   dept.ID,
		})
		departmentNodes[i] = dept
	}

	// Создаем команды для каждого отдела
	for _, dept := range departmentNodes {
		for i := 0; i < numTeams; i++ {
			team := &OrgNode{
				ID:       uuid.New(),
				Name:     fmt.Sprintf("Команда %d", i),
				SysName:  fmt.Sprintf("team_%d", i),
				Type:     teamType,
				Position: managerPosition,
			}
			builder.AddNode(team)
			builder.AddEdge(&Edge{
				FromNode: dept.ID,
				ToNode:   team.ID,
			})

			// Создаем сотрудников для каждой команды
			for j := 0; j < numEmployees; j++ {
				employee := &EmployeeNode{
					ID:   uuid.New(),
					Name: fmt.Sprintf("Сотрудник %d", j),
					Type: employeeType,
				}
				builder.AddNode(employee)
				builder.AddEdge(&Edge{
					FromNode: team.ID,
					ToNode:   employee.ID,
				})
			}
		}
	}

	// Строим дерево
	tree := builder.BuildTree()

	// Определяем различные типы фильтрации
	filters := []struct {
		name      string
		predicate func(interface{}) bool
	}{
		{
			"По типу (employee)",
			func(value interface{}) bool {
				if orgNode, ok := value.(*EmployeeNode); ok {
					return orgNode.Type != nil && orgNode.Type.SysName == "employee"
				}
				return false
			},
		},
		{
			"По имени (содержит 'Команда')",
			func(value interface{}) bool {
				if orgNode, ok := value.(*OrgNode); ok {
					return orgNode.Name != "" && strings.Contains(orgNode.Name, "Команда")
				}
				return false
			},
		},
		{
			"По регулярному выражению (team_)",
			func(value interface{}) bool {
				if orgNode, ok := value.(*OrgNode); ok {
					matched, _ := regexp.MatchString("team_", orgNode.SysName)
					return matched
				}
				return false
			},
		},
	}

	// Запускаем бенчмарки для каждого типа фильтрации
	for _, filter := range filters {
		b.Run(filter.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				filteredTree := tree.FilterSubtree(filter.predicate)
				if filteredTree == nil {
					b.Fatal("Filtered tree is nil")
				}
			}
		})
	}
}
