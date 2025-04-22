package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/arsants/orgtree"
	"github.com/google/uuid"
)

func printJSON(data []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "    "); err != nil {
		log.Printf("Ошибка при форматировании JSON: %v", err)
		fmt.Println(string(data))
		return
	}
	fmt.Println(prettyJSON.String())
}

func main() {
	// Создаем новое дерево
	root := orgtree.NewNode("CEO")

	// Добавляем несколько узлов для демонстрации
	root.Value = "CEO"

	// Создаем поддерево для технического отдела
	cto := orgtree.NewNode("CTO")
	root.AddChild(cto)

	dev1 := orgtree.NewNode("Senior Developer")
	dev2 := orgtree.NewNode("Junior Developer")
	dev3 := orgtree.NewNode("Middle Developer")
	cto.AddChild(dev1)
	cto.AddChild(dev2)
	cto.AddChild(dev3)

	// Добавляем поддерево для финансового отдела
	cfo := orgtree.NewNode("CFO")
	root.AddChild(cfo)

	acc1 := orgtree.NewNode("Senior Accountant")
	acc2 := orgtree.NewNode("Junior Accountant")
	acc3 := orgtree.NewNode("Middle Accountant")
	cfo.AddChild(acc1)
	cfo.AddChild(acc2)
	cfo.AddChild(acc3)

	// Добавляем поддерево для HR отдела
	hr := orgtree.NewNode("HR Manager")
	root.AddChild(hr)

	hr1 := orgtree.NewNode("Senior HR Specialist")
	hr2 := orgtree.NewNode("Junior HR Specialist")
	hr.AddChild(hr1)
	hr.AddChild(hr2)

	fmt.Println("=== Структура организации ===")

	// Сериализуем дерево в JSON
	jsonData, err := root.ToJSON()
	if err != nil {
		log.Fatalf("Ошибка при сериализации: %v", err)
	}

	// Выводим отформатированный JSON
	printJSON(jsonData)

	// Выводим хеш дерева
	fmt.Printf("\n=== Хеш дерева ===\n%s\n", root.HashString())

	// Сохраняем JSON в файл
	err = os.WriteFile("orgtree.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Ошибка при сохранении файла: %v", err)
	}
	fmt.Println("\nСтруктура сохранена в файл orgtree.json")

	// Демонстрация фильтрации поддерева
	fmt.Println("\n=== Демонстрация фильтрации поддерева ===")

	// Фильтруем только узлы, содержащие "Senior"
	pattern := regexp.MustCompile("Senior")
	seniorSubtree, err := root.FilterSubtreeByRegex(pattern.String())
	if err != nil {
		log.Fatalf("Ошибка при фильтрации Senior: %v", err)
	}

	if seniorSubtree != nil {
		seniorJSON, _ := seniorSubtree.ToJSON()
		fmt.Println("Senior сотрудники:")
		printJSON(seniorJSON)
	} else {
		fmt.Println("Senior сотрудники не найдены")
	}

	// Фильтруем только узлы, содержащие "Developer"
	pattern = regexp.MustCompile("Developer")
	devSubtree, err := root.FilterSubtreeByRegex(pattern.String())
	if err != nil {
		log.Fatalf("Ошибка при фильтрации Developer: %v", err)
	}

	if devSubtree != nil {
		devJSON, _ := devSubtree.ToJSON()
		fmt.Println("\nDeveloper сотрудники:")
		printJSON(devJSON)
	} else {
		fmt.Println("Developer сотрудники не найдены")
	}

	// Фильтруем только узлы, содержащие "Accountant"
	pattern = regexp.MustCompile("Accountant")
	accSubtree, err := root.FilterSubtreeByRegex(pattern.String())
	if err != nil {
		log.Fatalf("Ошибка при фильтрации Accountant: %v", err)
	}

	if accSubtree != nil {
		accJSON, _ := accSubtree.ToJSON()
		fmt.Println("\nAccountant сотрудники:")
		printJSON(accJSON)
	} else {
		fmt.Println("Accountant сотрудники не найдены")
	}

	// Демонстрация обхода дерева
	fmt.Println("\n=== Демонстрация обхода дерева ===")

	// Демонстрация фильтрации по условию
	fmt.Println("\n=== Демонстрация фильтрации по условию ===")
	predicate := func(value interface{}) bool {
		if str, ok := value.(string); ok {
			return str == "CEO" || str == "CTO" || str == "CFO"
		}
		return false
	}

	managementSubtree := root.FilterSubtree(predicate)
	if managementSubtree != nil {
		managementJSON, _ := managementSubtree.ToJSON()
		fmt.Println("Топ-менеджмент:")
		printJSON(managementJSON)
	}

	// Демонстрация получения пути к узлу
	fmt.Println("\n=== Демонстрация получения пути к узлу ===")
	path := dev1.GetPath(root)
	fmt.Println("Путь к Senior Developer:")
	for _, node := range path {
		fmt.Printf("- %s\n", node.Value)
	}

	// Демонстрация получения глубины узла
	fmt.Println("\n=== Демонстрация получения глубины узла ===")
	depth, ok := dev1.GetDepth(root)
	if !ok {
		fmt.Println("Не удалось получить глубину узла")
	} else {
		fmt.Printf("Глубина узла Senior Developer: %d\n", depth)
	}

	// Демонстрация работы с TreeBuilder
	fmt.Println("\n=== Демонстрация работы с TreeBuilder ===")

	// Создаем построитель дерева
	builder := orgtree.NewTreeBuilder()

	// Создаем тестовые данные
	nodeType := &orgtree.NodeType{
		ID:      uuid.New(),
		Name:    "Отдел",
		SysName: "department",
	}

	mainOffice := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "Главный офис",
		SysName: "main_office",
		Type:    nodeType,
		Position: &orgtree.Position{
			ID:      uuid.New(),
			Name:    "Директор",
			SysName: "director",
		},
	}

	// Добавляем данные в построитель
	builder.AddNode(mainOffice)

	// Создаем и добавляем дополнительные узлы
	itDept := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "IT отдел",
		SysName: "it_department",
		Type:    nodeType,
		Position: &orgtree.Position{
			ID:      uuid.New(),
			Name:    "IT менеджер",
			SysName: "it_manager",
		},
	}
	builder.AddNode(itDept)

	// Создаем подразделения IT отдела
	itTeams := []struct {
		name    string
		sysName string
	}{
		{"Команда разработки", "development_team"},
		{"Команда тестирования", "qa_team"},
		{"DevOps команда", "devops_team"},
	}

	itSubteams := map[string][]struct {
		name    string
		sysName string
	}{
		"development_team": {
			{"Frontend команда", "frontend_team"},
			{"Backend команда", "backend_team"},
			{"Mobile команда", "mobile_team"},
		},
		"qa_team": {
			{"Команда автоматизации", "automation_team"},
			{"Команда ручного тестирования", "manual_team"},
		},
	}

	// Создаем узлы для IT команд
	itTeamNodes := make(map[string]*orgtree.OrgNode)
	for _, team := range itTeams {
		teamNode := &orgtree.OrgNode{
			ID:      uuid.New(),
			Name:    team.name,
			SysName: team.sysName,
			Type:    nodeType,
		}
		builder.AddNode(teamNode)
		itTeamNodes[team.sysName] = teamNode

		// Добавляем подкоманды
		if subteams, ok := itSubteams[team.sysName]; ok {
			for _, subteam := range subteams {
				subnode := &orgtree.OrgNode{
					ID:      uuid.New(),
					Name:    subteam.name,
					SysName: subteam.sysName,
					Type:    nodeType,
				}
				builder.AddNode(subnode)
				builder.AddEdge(&orgtree.Edge{
					FromNode: teamNode.ID,
					ToNode:   subnode.ID,
				})
			}
		}
	}

	// Создаем HR отдел
	hrDept := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "Отдел кадров",
		SysName: "hr_department",
		Type:    nodeType,
		Position: &orgtree.Position{
			ID:      uuid.New(),
			Name:    "HR менеджер",
			SysName: "hr_manager",
		},
	}
	builder.AddNode(hrDept)

	// Создаем подразделения HR отдела
	hrTeams := []struct {
		name    string
		sysName string
	}{
		{"Команда рекрутинга", "recruitment_team"},
		{"Команда обучения", "training_team"},
	}

	hrSubteams := map[string][]struct {
		name    string
		sysName string
	}{
		"recruitment_team": {
			{"Команда технического рекрутинга", "tech_recruitment_team"},
			{"Команда HR рекрутинга", "hr_recruitment_team"},
		},
		"training_team": {
			{"Команда корпоративного обучения", "corporate_training_team"},
			{"Команда адаптации", "onboarding_team"},
		},
	}

	// Создаем узлы для HR команд
	hrTeamNodes := make(map[string]*orgtree.OrgNode)
	for _, team := range hrTeams {
		teamNode := &orgtree.OrgNode{
			ID:      uuid.New(),
			Name:    team.name,
			SysName: team.sysName,
			Type:    nodeType,
		}
		builder.AddNode(teamNode)
		hrTeamNodes[team.sysName] = teamNode

		// Добавляем подкоманды
		if subteams, ok := hrSubteams[team.sysName]; ok {
			for _, subteam := range subteams {
				subnode := &orgtree.OrgNode{
					ID:      uuid.New(),
					Name:    subteam.name,
					SysName: subteam.sysName,
					Type:    nodeType,
				}
				builder.AddNode(subnode)
				builder.AddEdge(&orgtree.Edge{
					FromNode: teamNode.ID,
					ToNode:   subnode.ID,
				})
			}
		}
	}

	// Создаем должность по умолчанию для сотрудников
	defaultPosition := &orgtree.Position{
		ID:      uuid.New(),
		Name:    "Сотрудник",
		SysName: "employee",
	}

	// Добавляем связи между узлами
	builder.AddEdge(&orgtree.Edge{
		FromNode: mainOffice.ID,
		ToNode:   itDept.ID,
	})
	builder.AddEdge(&orgtree.Edge{
		FromNode: mainOffice.ID,
		ToNode:   hrDept.ID,
	})

	// Добавляем связи для IT отдела
	for _, teamNode := range itTeamNodes {
		builder.AddEdge(&orgtree.Edge{
			FromNode: itDept.ID,
			ToNode:   teamNode.ID,
		})
	}

	// Добавляем связи для HR отдела
	for _, teamNode := range hrTeamNodes {
		builder.AddEdge(&orgtree.Edge{
			FromNode: hrDept.ID,
			ToNode:   teamNode.ID,
		})
	}

	// Создаем должности для команд
	teamPositionsData := []struct {
		name    string
		sysName string
	}{
		{"Руководитель разработки", "development_lead"},
		{"Руководитель тестирования", "qa_lead"},
		{"Руководитель DevOps", "devops_lead"},
		{"Руководитель рекрутинга", "recruitment_lead"},
		{"Руководитель обучения", "training_lead"},
	}

	// Создаем должности для подкоманд
	subteamPositionsData := map[string][]struct {
		name    string
		sysName string
	}{
		"development_team": {
			{"Руководитель Frontend", "frontend_lead"},
			{"Руководитель Backend", "backend_lead"},
			{"Руководитель Mobile", "mobile_lead"},
		},
		"qa_team": {
			{"Руководитель автоматизации", "automation_lead"},
			{"Руководитель ручного тестирования", "manual_qa_lead"},
		},
		"recruitment_team": {
			{"Руководитель технического рекрутинга", "tech_recruitment_lead"},
			{"Руководитель HR рекрутинга", "hr_recruitment_lead"},
		},
		"training_team": {
			{"Руководитель корпоративного обучения", "corporate_training_lead"},
			{"Руководитель адаптации", "onboarding_lead"},
		},
	}

	// Создаем и сохраняем должности для дальнейшего присвоения
	positionsMap := make(map[string]*orgtree.Position)
	for _, posData := range teamPositionsData {
		position := &orgtree.Position{
			ID:      uuid.New(),
			Name:    posData.name,
			SysName: posData.sysName,
		}
		positionsMap[posData.sysName] = position
	}

	// Добавляем должности для подкоманд
	for _, subPosDataList := range subteamPositionsData {
		for _, posData := range subPosDataList {
			position := &orgtree.Position{
				ID:      uuid.New(),
				Name:    posData.name,
				SysName: posData.sysName,
			}
			positionsMap[posData.sysName] = position
		}
	}

	// Добавляем связи должностей для команд (Присваиваем лидов)
	leadPositionMap := map[string]string{
		"development_team": "development_lead",
		"qa_team":          "qa_lead",
		"devops_team":      "devops_lead",
		"recruitment_team": "recruitment_lead",
		"training_team":    "training_lead",
	}
	for sysName, teamNode := range itTeamNodes {
		if leadSysName, ok := leadPositionMap[sysName]; ok {
			if pos, posOK := positionsMap[leadSysName]; posOK {
				teamNode.Position = pos
			}
		}
	}
	for sysName, teamNode := range hrTeamNodes {
		if leadSysName, ok := leadPositionMap[sysName]; ok {
			if pos, posOK := positionsMap[leadSysName]; posOK {
				teamNode.Position = pos
			}
		}
	}

	// Добавляем связи должностей для подкоманд (Присваиваем лидов)
	subteamLeadPositionMap := map[string]string{
		"frontend_team":           "frontend_lead",
		"backend_team":            "backend_lead",
		"mobile_team":             "mobile_lead",
		"automation_team":         "automation_lead",
		"manual_team":             "manual_qa_lead",
		"tech_recruitment_team":   "tech_recruitment_lead",
		"hr_recruitment_team":     "hr_recruitment_lead",
		"corporate_training_team": "corporate_training_lead",
		"onboarding_team":         "onboarding_lead",
	}

	// Проходим по всем узлам и присваиваем должности подкомандам или дефолтную
	for _, orgNode := range builder.Nodes() {
		// Если должность уже установлена (например, для корневых отделов или команд), пропускаем
		if orgNode.Position != nil {
			continue
		}
		// Проверяем, является ли узел подкомандой и есть ли для него лид-должность
		if leadSysName, leadExists := subteamLeadPositionMap[orgNode.SysName]; leadExists {
			if pos, posOK := positionsMap[leadSysName]; posOK {
				orgNode.Position = pos // Assign specific subteam lead position
			} else {
				orgNode.Position = defaultPosition // Assign default if lead pos not found in map
			}
		} else {
			orgNode.Position = defaultPosition // Assign default if not a subteam lead
		}
	}

	// Построение дерева
	orgTree := builder.BuildTree()

	// Выводим результат
	orgTreeJSON, _ := orgTree.ToJSON()
	fmt.Println("Дерево, построенное с помощью TreeBuilder:")
	printJSON(orgTreeJSON)

	fmt.Println("\n=== Структура дерева ===")
	orgTree.PrintTree()

	// Пример использования WalkTree
	fmt.Println("\n=== Обход дерева с отступами ===")
	orgTree.WalkTree(func(node *orgtree.Node, depth int) {
		indent := strings.Repeat("  ", depth)
		if orgNode, ok := node.Value.(*orgtree.OrgNode); ok {
			fmt.Printf("%s- %s\n", indent, orgNode.Name)
		}
	})

	// Демонстрация поиска в построенном дереве
	fmt.Println("\n=== Демонстрация поиска в построенном дереве ===")

	// Поиск отдела по имени
	fmt.Println("\nПоиск отдела по имени:")
	searchPredicate := func(value interface{}) bool {
		if orgNode, ok := value.(*orgtree.OrgNode); ok {
			return orgNode.Name == "IT отдел"
		}
		return false
	}
	if found := orgTree.Filter(searchPredicate); found != nil {
		if orgNode, ok := found.Value.(*orgtree.OrgNode); ok {
			fmt.Printf("Найден отдел: %s (ID: %s)\n", orgNode.Name, orgNode.ID)
		}
	}

	// Поиск всех команд разработки
	fmt.Println("\nПоиск всех команд разработки:")
	devPattern := regexp.MustCompile("разработки")
	devSubtree, searchErr := orgTree.FilterSubtreeByRegex(devPattern.String())
	if searchErr != nil {
		log.Fatalf("Ошибка при поиске команд разработки: %v", searchErr)
	} else if devSubtree != nil {
		devJSON, err := devSubtree.ToJSON()
		if err != nil {
			log.Fatalf("Ошибка при сериализации найденных команд разработки: %v", err)
		}
		fmt.Println("Найденные команды разработки:")
		printJSON(devJSON)
	}

	// Поиск всех команд, содержащих определенное слово
	fmt.Println("\nПоиск всех команд, содержащих 'команда':")
	teamPredicate := func(value interface{}) bool {
		if orgNode, ok := value.(*orgtree.OrgNode); ok {
			return strings.Contains(orgNode.Name, "команда")
		}
		return false
	}

	teamSubtree := orgTree.FilterSubtree(teamPredicate)
	if teamSubtree != nil {
		teamJSON, err := teamSubtree.ToJSON()
		if err != nil {
			log.Fatalf("Ошибка при сериализации найденных команд: %v", err)
		}
		fmt.Println("Найденные команды:")
		printJSON(teamJSON)
	}
}
