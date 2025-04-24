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
		ID:        uuid.New(),
		Name:      "IT отдел",
		SysName:   "it_department",
		Type:      departmentType,
		Positions: []*Position{managerPosition},
	}
	devTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Команда разработки",
		SysName:   "dev_team",
		Type:      departmentType,
		Positions: []*Position{engineerPosition},
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
	if itDeptNodeValue.Positions == nil {
		t.Errorf("В IT отделе ожидалась должность, получено nil")
	} else if itDeptNodeValue.Positions[0].ID != managerPosition.ID {
		t.Errorf("Неверная должность в IT отделе, ожидался ID %s, получен %s", managerPosition.ID, itDeptNodeValue.Positions[0].ID)
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
	if devTeamNodeValue.Positions == nil {
		t.Errorf("В команде разработки ожидалась должность, получено nil")
	} else if devTeamNodeValue.Positions[0].ID != engineerPosition.ID {
		t.Errorf("Неверная должность в команде разработки, ожидался ID %s, получен %s", engineerPosition.ID, devTeamNodeValue.Positions[0].ID)
	}
}

func TestTreeBuilderWithEmployees(t *testing.T) {
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
		Name:    "Руководитель команды",
		SysName: "team_lead",
	}

	developerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Разработчик",
		SysName: "developer",
	}

	// Создаем отдел
	itDepartment := &OrgNode{
		ID:        uuid.New(),
		Name:      "IT отдел",
		SysName:   "it_department",
		Type:      departmentType,
		Positions: []*Position{managerPosition},
	}

	// Создаем команду
	devTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Команда разработки",
		SysName:   "dev_team",
		Type:      teamType,
		Positions: []*Position{managerPosition},
	}

	// Создаем сотрудников
	employee1 := &OrgNode{
		ID:        uuid.New(),
		Name:      "Иван Иванов",
		SysName:   "ivan_ivanov",
		Type:      employeeType,
		Positions: []*Position{developerPosition},
	}

	employee2 := &OrgNode{
		ID:        uuid.New(),
		Name:      "Петр Петров",
		SysName:   "petr_petrov",
		Type:      employeeType,
		Positions: []*Position{developerPosition},
	}

	// Добавляем узлы
	builder.AddNode(itDepartment)
	builder.AddNode(devTeam)
	builder.AddNode(employee1)
	builder.AddNode(employee2)

	// Добавляем связи
	builder.AddEdge(&Edge{
		FromNode: itDepartment.ID,
		ToNode:   devTeam.ID,
	})

	builder.AddEdge(&Edge{
		FromNode: devTeam.ID,
		ToNode:   employee1.ID,
	})

	builder.AddEdge(&Edge{
		FromNode: devTeam.ID,
		ToNode:   employee2.ID,
	})

	// Строим дерево
	tree := builder.BuildTree()

	// Проверяем структуру дерева
	if len(tree.Children) != 1 {
		t.Fatalf("Ожидался один корневой узел (IT отдел), получено %d", len(tree.Children))
	}

	// Проверяем IT отдел
	itDeptNode := tree.Children[0]
	itDeptValue, ok := itDeptNode.Value.(*OrgNode)
	if !ok || itDeptValue.ID != itDepartment.ID {
		t.Fatalf("Неверный корневой узел, ожидался IT отдел")
	}

	// Проверяем команду разработки
	if len(itDeptNode.Children) != 1 {
		t.Fatalf("Ожидалась одна команда в IT отделе, получено %d", len(itDeptNode.Children))
	}

	devTeamNode := itDeptNode.Children[0]
	devTeamValue, ok := devTeamNode.Value.(*OrgNode)
	if !ok || devTeamValue.ID != devTeam.ID {
		t.Fatalf("Неверный узел команды разработки")
	}

	// Проверяем сотрудников
	if len(devTeamNode.Children) != 2 {
		t.Fatalf("Ожидалось два сотрудника в команде разработки, получено %d", len(devTeamNode.Children))
	}

	// Проверяем первого сотрудника
	employee1Node := devTeamNode.Children[0]
	employee1Value, ok := employee1Node.Value.(*OrgNode)
	if !ok || employee1Value.ID != employee1.ID {
		t.Errorf("Неверный узел первого сотрудника")
	}
	if employee1Value.Type.SysName != "employee" {
		t.Errorf("Неверный тип узла первого сотрудника, ожидался 'employee', получен '%s'", employee1Value.Type.SysName)
	}
	if employee1Value.Positions[0].ID != developerPosition.ID {
		t.Errorf("Неверная должность первого сотрудника")
	}

	// Проверяем второго сотрудника
	employee2Node := devTeamNode.Children[1]
	employee2Value, ok := employee2Node.Value.(*OrgNode)
	if !ok || employee2Value.ID != employee2.ID {
		t.Errorf("Неверный узел второго сотрудника")
	}
	if employee2Value.Type.SysName != "employee" {
		t.Errorf("Неверный тип узла второго сотрудника, ожидался 'employee', получен '%s'", employee2Value.Type.SysName)
	}
	if employee2Value.Positions[0].ID != developerPosition.ID {
		t.Errorf("Неверная должность второго сотрудника")
	}
}

func TestTreeBuilderFilterEmployees(t *testing.T) {
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

	subteamType := &NodeType{
		ID:      uuid.New(),
		Name:    "Подкоманда",
		SysName: "subteam",
	}

	employeeType := &NodeType{
		ID:      uuid.New(),
		Name:    "Сотрудник",
		SysName: "employee",
	}

	// Создаем должности
	managerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Руководитель команды",
		SysName: "team_lead",
	}

	developerPosition := &Position{
		ID:      uuid.New(),
		Name:    "Разработчик",
		SysName: "developer",
	}

	qaPosition := &Position{
		ID:      uuid.New(),
		Name:    "Тестировщик",
		SysName: "qa_engineer",
	}

	// Создаем отдел
	itDepartment := &OrgNode{
		ID:        uuid.New(),
		Name:      "IT отдел",
		SysName:   "it_department",
		Type:      departmentType,
		Positions: []*Position{managerPosition},
	}

	// Создаем команды
	devTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Команда разработки",
		SysName:   "dev_team",
		Type:      teamType,
		Positions: []*Position{managerPosition},
	}

	qaTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Команда тестирования",
		SysName:   "qa_team",
		Type:      teamType,
		Positions: []*Position{managerPosition},
	}

	// Создаем подкоманды
	frontendTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Frontend команда",
		SysName:   "frontend_team",
		Type:      subteamType,
		Positions: []*Position{managerPosition},
	}

	backendTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Backend команда",
		SysName:   "backend_team",
		Type:      subteamType,
		Positions: []*Position{managerPosition},
	}

	automationTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Команда автоматизации",
		SysName:   "automation_team",
		Type:      subteamType,
		Positions: []*Position{managerPosition},
	}

	manualTeam := &OrgNode{
		ID:        uuid.New(),
		Name:      "Команда ручного тестирования",
		SysName:   "manual_team",
		Type:      subteamType,
		Positions: []*Position{managerPosition},
	}

	// Создаем сотрудников
	employees := []struct {
		name     string
		sysName  string
		position *Position
	}{
		// Разработчики Frontend
		{"Иван Иванов", "ivan_ivanov", developerPosition},
		{"Петр Петров", "petr_petrov", developerPosition},
		{"Анна Сидорова", "anna_sidorova", developerPosition},

		// Разработчики Backend
		{"Сергей Сергеев", "sergey_sergeev", developerPosition},
		{"Мария Петрова", "maria_petrova", developerPosition},
		{"Алексей Алексеев", "alexey_alexeev", developerPosition},

		// Тестировщики автоматизации
		{"Елена Еленова", "elena_elenova", qaPosition},
		{"Дмитрий Дмитриев", "dmitry_dmitriev", qaPosition},

		// Тестировщики ручного тестирования
		{"Ольга Ольгова", "olga_olgova", qaPosition},
		{"Николай Николаев", "nikolay_nikolaev", qaPosition},
	}

	// Добавляем узлы
	builder.AddNode(itDepartment)
	builder.AddNode(devTeam)
	builder.AddNode(qaTeam)
	builder.AddNode(frontendTeam)
	builder.AddNode(backendTeam)
	builder.AddNode(automationTeam)
	builder.AddNode(manualTeam)

	// Создаем и добавляем сотрудников
	employeeNodes := make([]*EmployeeNode, len(employees))
	for i, emp := range employees {
		employee := &EmployeeNode{
			ID:   uuid.New(),
			Name: emp.name,
			Type: employeeType,
		}
		employeeNodes[i] = employee
		builder.AddNode(employee)
	}

	// Добавляем связи между отделами и командами
	builder.AddEdge(&Edge{
		FromNode: itDepartment.ID,
		ToNode:   devTeam.ID,
	})
	builder.AddEdge(&Edge{
		FromNode: itDepartment.ID,
		ToNode:   qaTeam.ID,
	})

	// Добавляем связи между командами и подкомандами
	builder.AddEdge(&Edge{
		FromNode: devTeam.ID,
		ToNode:   frontendTeam.ID,
	})
	builder.AddEdge(&Edge{
		FromNode: devTeam.ID,
		ToNode:   backendTeam.ID,
	})
	builder.AddEdge(&Edge{
		FromNode: qaTeam.ID,
		ToNode:   automationTeam.ID,
	})
	builder.AddEdge(&Edge{
		FromNode: qaTeam.ID,
		ToNode:   manualTeam.ID,
	})

	// Распределяем сотрудников по подкомандам
	// Frontend команда
	builder.AddEdge(&Edge{
		FromNode: frontendTeam.ID,
		ToNode:   employeeNodes[0].ID,
	})
	builder.AddEdge(&Edge{
		FromNode: frontendTeam.ID,
		ToNode:   employeeNodes[1].ID,
	})
	builder.AddEdge(&Edge{
		FromNode: frontendTeam.ID,
		ToNode:   employeeNodes[2].ID,
	})

	// Backend команда
	builder.AddEdge(&Edge{
		FromNode: backendTeam.ID,
		ToNode:   employeeNodes[3].ID,
	})
	builder.AddEdge(&Edge{
		FromNode: backendTeam.ID,
		ToNode:   employeeNodes[4].ID,
	})
	builder.AddEdge(&Edge{
		FromNode: backendTeam.ID,
		ToNode:   employeeNodes[5].ID,
	})

	// Команда автоматизации
	builder.AddEdge(&Edge{
		FromNode: automationTeam.ID,
		ToNode:   employeeNodes[6].ID,
	})
	builder.AddEdge(&Edge{
		FromNode: automationTeam.ID,
		ToNode:   employeeNodes[7].ID,
	})

	// Команда ручного тестирования
	builder.AddEdge(&Edge{
		FromNode: manualTeam.ID,
		ToNode:   employeeNodes[8].ID,
	})
	builder.AddEdge(&Edge{
		FromNode: manualTeam.ID,
		ToNode:   employeeNodes[9].ID,
	})

	// Строим дерево
	tree := builder.BuildTree()

	// Создаем предикат для фильтрации сотрудников
	employeePredicate := func(value interface{}) bool {
		if employeeNode, ok := value.(*EmployeeNode); ok {
			return employeeNode.Type != nil && employeeNode.Type.SysName == "employee"
		}
		return false
	}

	// Фильтруем дерево
	filteredTree := tree.FilterSubtree(employeePredicate)

	// Проверяем результат фильтрации
	if filteredTree == nil {
		t.Fatal("Ожидалось непустое отфильтрованное дерево")
	}

	// Проверяем, что в отфильтрованном дереве только сотрудники
	iterator := NewPreOrderIterator(filteredTree)
	employeeCount := 0
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		if node.Value == nil {
			continue // Пропускаем корневой узел-обертку
		}

		switch node.Value.(type) {
		case *OrgNode:
			continue
		case *EmployeeNode:
			employeeCount++
		default:
			t.Errorf("Неизвестный тип узла: %T", node.Value)
		}

	}

	// Проверяем количество найденных сотрудников
	if employeeCount != len(employees) {
		t.Errorf("Ожидалось %d сотрудников, найдено %d", len(employees), employeeCount)
	}
}
