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
		Type:    nodeType,
	}

	node2 := &OrgNode{
		ID:      uuid.New(),
		Name:    "IT отдел",
		SysName: "it_department",
		Type:    nodeType,
	}

	edge := &Edge{
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

func TestTreeBuilderWithDifferentDepths(t *testing.T) {
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

	// Создаем корневой узел
	root := &OrgNode{
		ID:      uuid.New(),
		Name:    "Главный офис",
		SysName: "main_office",
		Type:    departmentType,
	}
	builder.AddNode(root)

	// Создаем отделы первого уровня
	departments := []struct {
		name    string
		sysName string
	}{
		{"IT отдел", "it_department"},
		{"HR отдел", "hr_department"},
		{"Финансовый отдел", "finance_department"},
	}

	// Создаем команды для каждого отдела
	teams := map[string][]struct {
		name    string
		sysName string
	}{
		"it_department": {
			{"Команда разработки", "dev_team"},
			{"Команда тестирования", "qa_team"},
		},
		"hr_department": {
			{"Команда рекрутинга", "recruitment_team"},
			{"Команда обучения", "training_team"},
		},
		"finance_department": {
			{"Команда бухгалтерии", "accounting_team"},
			{"Команда аудита", "audit_team"},
		},
	}

	// Создаем подкоманды для каждой команды
	subteams := map[string][]struct {
		name    string
		sysName string
	}{
		"dev_team": {
			{"Frontend команда", "frontend_team"},
			{"Backend команда", "backend_team"},
		},
		"qa_team": {
			{"Команда автоматизации", "automation_team"},
			{"Команда ручного тестирования", "manual_team"},
		},
		"recruitment_team": {
			{"Команда технического рекрутинга", "tech_recruitment_team"},
			{"Команда HR рекрутинга", "hr_recruitment_team"},
		},
		"training_team": {
			{"Команда корпоративного обучения", "corporate_training_team"},
			{"Команда адаптации", "onboarding_team"},
		},
		"accounting_team": {
			{"Команда первичной документации", "primary_docs_team"},
			{"Команда отчетности", "reporting_team"},
		},
		"audit_team": {
			{"Команда внутреннего аудита", "internal_audit_team"},
			{"Команда внешнего аудита", "external_audit_team"},
		},
	}

	// Создаем и добавляем узлы для отделов
	departmentNodes := make(map[string]*OrgNode)
	for _, dept := range departments {
		node := &OrgNode{
			ID:      uuid.New(),
			Name:    dept.name,
			SysName: dept.sysName,
			Type:    departmentType,
		}
		builder.AddNode(node)
		builder.AddEdge(&Edge{
			FromNode: root.ID,
			ToNode:   node.ID,
		})
		departmentNodes[dept.sysName] = node

		// Создаем и добавляем узлы для команд
		if deptTeams, ok := teams[dept.sysName]; ok {
			for _, team := range deptTeams {
				teamNode := &OrgNode{
					ID:      uuid.New(),
					Name:    team.name,
					SysName: team.sysName,
					Type:    teamType,
				}
				builder.AddNode(teamNode)
				builder.AddEdge(&Edge{
					FromNode: node.ID,
					ToNode:   teamNode.ID,
				})

				// Создаем и добавляем узлы для подкоманд
				if teamSubteams, ok := subteams[team.sysName]; ok {
					for _, subteam := range teamSubteams {
						subteamNode := &OrgNode{
							ID:      uuid.New(),
							Name:    subteam.name,
							SysName: subteam.sysName,
							Type:    teamType,
						}
						builder.AddNode(subteamNode)
						builder.AddEdge(&Edge{
							FromNode: teamNode.ID,
							ToNode:   subteamNode.ID,
						})
					}
				}
			}
		}
	}

	// Строим дерево
	tree := builder.BuildTree()

	// Проверяем структуру дерева
	// Ожидаем, что корнем будет wrapper, у которого один child "Главный офис"
	if tree.Value != nil || len(tree.Children) != 1 {
		t.Fatalf("Ожидался wrapper с одним дочерним узлом 'Главный офис'")
	}
	mainOfficeNode := tree.Children[0]
	if mainOfficeNode == nil {
		t.Fatal("Не найден узел 'Главный офис'")
	}
	if mainOfficeOrgNode, ok := mainOfficeNode.Value.(*OrgNode); !ok || mainOfficeOrgNode.Name != "Главный офис" {
		t.Fatalf("Неверный корневой узел, ожидался 'Главный офис', получен %v", mainOfficeNode.Value)
	}

	// Проверяем отделы под "Главным офисом"
	if len(mainOfficeNode.Children) != len(departments) {
		t.Fatalf("Ожидалось %d отделов под 'Главным офисом', получено %d", len(departments), len(mainOfficeNode.Children))
	}

	// Проверяем каждый отдел
	for _, deptNode := range mainOfficeNode.Children {
		deptOrgNode, ok := deptNode.Value.(*OrgNode)
		if !ok || deptOrgNode.Type.SysName != "department" {
			t.Errorf("Ожидался узел типа 'department', получен %v", deptNode.Value)
			continue
		}

		// Проверяем количество команд (прямых детей отдела)
		expectedTeamCount := 0
		if deptTeams, exists := teams[deptOrgNode.SysName]; exists {
			expectedTeamCount = len(deptTeams)
		}
		if len(deptNode.Children) != expectedTeamCount {
			t.Errorf("В отделе '%s' ожидалось %d команд(ы), получено %d", deptOrgNode.Name, expectedTeamCount, len(deptNode.Children))
			continue
		}

		// Проверяем каждую команду
		for _, teamNode := range deptNode.Children {
			teamOrgNode, ok := teamNode.Value.(*OrgNode)
			if !ok || teamOrgNode.Type.SysName != "team" {
				t.Errorf("Ожидался узел типа 'team' под отделом '%s', получен %v", deptOrgNode.Name, teamNode.Value)
				continue
			}

			// Проверяем количество подкоманд (прямых детей команды)
			expectedSubteamCount := 0
			if teamSubteams, exists := subteams[teamOrgNode.SysName]; exists {
				expectedSubteamCount = len(teamSubteams)
			}
			if len(teamNode.Children) != expectedSubteamCount {
				t.Errorf("В команде '%s' (отдел '%s') ожидалось %d подкоманд(ы), получено %d", teamOrgNode.Name, deptOrgNode.Name, expectedSubteamCount, len(teamNode.Children))
			}
		}
	}

}

func TestTreeBuilderWithPositions(t *testing.T) {
	builder := NewTreeBuilder()

	// Типы узлов
	departmentType := &NodeType{ID: uuid.New(), Name: "Отдел", SysName: "department"}

	// Должности (с корректными полями)
	managerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Руководитель IT",
		SysName: "it_manager",
	}
	engineerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Инженер-разработчик",
		SysName: "dev_engineer",
	}

	// Узлы (с присвоением Position)
	itDepartment := &OrgNode{
		ID:       uuid.New(),
		Name:     "IT отдел",
		SysName:  "it_department",
		Type:     departmentType,
		Position: managerPosition,
	}
	devTeam := &OrgNode{
		ID:       uuid.New(),
		Name:     "Команда разработки",
		SysName:  "dev_team",
		Type:     departmentType,
		Position: engineerPosition,
	}

	// Добавляем узлы
	builder.AddNode(itDepartment)
	builder.AddNode(devTeam)

	// Добавляем связь
	builder.AddEdge(&Edge{
		FromNode: itDepartment.ID,
		ToNode:   devTeam.ID,
	})

	// Строим дерево
	tree := builder.BuildTree()

	// Проверка: Ищем узлы и проверяем наличие должностей
	// Корень - это обертка, пропускаем его
	if len(tree.Children) != 1 {
		t.Fatalf("Ожидался один корневой узел (IT отдел), получено %d", len(tree.Children))
	}
	rootOrgNode := tree.Children[0]

	// Проверяем IT отдел
	itDeptNodeValue, ok := rootOrgNode.Value.(*OrgNode)
	if !ok || itDeptNodeValue.ID != itDepartment.ID {
		t.Fatalf("Неверный корневой узел, ожидался IT отдел")
	}
	if itDeptNodeValue.Position == nil {
		t.Errorf("В IT отделе ожидалась должность, получено nil")
	} else if itDeptNodeValue.Position.ID != managerPosition.ID {
		t.Errorf("Неверная должность в IT отделе, ожидался ID %s, получен %s", managerPosition.ID, itDeptNodeValue.Position.ID)
	}

	// Проверяем команду разработки
	if len(rootOrgNode.Children) != 1 {
		t.Fatalf("Ожидался один дочерний узел (Команда разработки) у IT отдела, получено %d", len(rootOrgNode.Children))
	}
	devTeamNode := rootOrgNode.Children[0]
	devTeamNodeValue, ok := devTeamNode.Value.(*OrgNode)
	if !ok || devTeamNodeValue.ID != devTeam.ID {
		t.Fatalf("Неверный узел команды разработки")
	}
	if devTeamNodeValue.Position == nil {
		t.Errorf("В команде разработки ожидалась должность, получено nil")
	} else if devTeamNodeValue.Position.ID != engineerPosition.ID {
		t.Errorf("Неверная должность в команде разработки, ожидался ID %s, получен %s", engineerPosition.ID, devTeamNodeValue.Position.ID)
	}
}
