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

	node := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "Главный офис",
		SysName: "main_office",
		Type:    nodeType,
	}

	edge := &orgtree.Edge{
		ID:       uuid.New(),
		Name:     "Подчинение",
		SysName:  "subordination",
		FromNode: node.ID,
		ToNode:   uuid.New(), // Для демонстрации
	}

	position := &orgtree.Position{
		ID:      uuid.New(),
		Name:    "Директор",
		SysName: "director",
	}

	relation := &orgtree.PositionNodeRelation{
		ID:       uuid.New(),
		NodeID:   node.ID,
		Position: position,
	}

	// Добавляем данные в построитель
	builder.AddNode(node)

	// Создаем и добавляем дополнительные узлы
	itDept := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "IT отдел",
		SysName: "it_department",
		Type:    nodeType,
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
		node := &orgtree.OrgNode{
			ID:      uuid.New(),
			Name:    team.name,
			SysName: team.sysName,
			Type:    nodeType,
		}
		builder.AddNode(node)
		itTeamNodes[team.sysName] = node

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
					ID:       uuid.New(),
					Name:     "Подчинение",
					SysName:  "subordination",
					FromNode: node.ID,
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
		node := &orgtree.OrgNode{
			ID:      uuid.New(),
			Name:    team.name,
			SysName: team.sysName,
			Type:    nodeType,
		}
		builder.AddNode(node)
		hrTeamNodes[team.sysName] = node

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
					ID:       uuid.New(),
					Name:     "Подчинение",
					SysName:  "subordination",
					FromNode: node.ID,
					ToNode:   subnode.ID,
				})
			}
		}
	}

	// Добавляем связи между узлами
	builder.AddEdge(edge)
	builder.AddEdge(&orgtree.Edge{
		ID:       uuid.New(),
		Name:     "Подчинение",
		SysName:  "subordination",
		FromNode: node.ID,
		ToNode:   itDept.ID,
	})
	builder.AddEdge(&orgtree.Edge{
		ID:       uuid.New(),
		Name:     "Подчинение",
		SysName:  "subordination",
		FromNode: node.ID,
		ToNode:   hrDept.ID,
	})

	// Добавляем связи для IT отдела
	for _, teamNode := range itTeamNodes {
		builder.AddEdge(&orgtree.Edge{
			ID:       uuid.New(),
			Name:     "Подчинение",
			SysName:  "subordination",
			FromNode: itDept.ID,
			ToNode:   teamNode.ID,
		})
	}

	// Добавляем связи для HR отдела
	for _, teamNode := range hrTeamNodes {
		builder.AddEdge(&orgtree.Edge{
			ID:       uuid.New(),
			Name:     "Подчинение",
			SysName:  "subordination",
			FromNode: hrDept.ID,
			ToNode:   teamNode.ID,
		})
	}

	// Добавляем должности
	builder.AddPosition(position)
	itManager := &orgtree.Position{
		ID:      uuid.New(),
		Name:    "IT менеджер",
		SysName: "it_manager",
	}
	builder.AddPosition(itManager)
	hrManager := &orgtree.Position{
		ID:      uuid.New(),
		Name:    "HR менеджер",
		SysName: "hr_manager",
	}
	builder.AddPosition(hrManager)

	// Создаем должности для команд
	teamPositions := []struct {
		name    string
		sysName string
	}{
		{"Руководитель разработки", "development_lead"},
		{"Руководитель тестирования", "qa_lead"},
		{"Руководитель DevOps", "devops_lead"},
		{"Руководитель рекрутинга", "recruitment_lead"},
		{"Руководитель обучения", "training_lead"},
	}

	// Добавляем должности для подкоманд
	subteamPositions := map[string][]struct {
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

	// Добавляем должности и их связи
	positionNodes := make(map[string]*orgtree.Position)
	for _, pos := range teamPositions {
		position := &orgtree.Position{
			ID:      uuid.New(),
			Name:    pos.name,
			SysName: pos.sysName,
		}
		builder.AddPosition(position)
		positionNodes[pos.sysName] = position
	}

	// Добавляем должности для подкоманд
	for _, positions := range subteamPositions {
		for _, pos := range positions {
			position := &orgtree.Position{
				ID:      uuid.New(),
				Name:    pos.name,
				SysName: pos.sysName,
			}
			builder.AddPosition(position)
			positionNodes[pos.sysName] = position
		}
	}

	// Добавляем связи должностей с узлами
	builder.AddPositionNodeRelation(relation)
	builder.AddPositionNodeRelation(&orgtree.PositionNodeRelation{
		ID:       uuid.New(),
		NodeID:   itDept.ID,
		Position: itManager,
	})
	builder.AddPositionNodeRelation(&orgtree.PositionNodeRelation{
		ID:       uuid.New(),
		NodeID:   hrDept.ID,
		Position: hrManager,
	})

	// Добавляем связи должностей для команд
	for sysName, teamNode := range itTeamNodes {
		if pos, ok := positionNodes[sysName+"_lead"]; ok {
			builder.AddPositionNodeRelation(&orgtree.PositionNodeRelation{
				ID:       uuid.New(),
				NodeID:   teamNode.ID,
				Position: pos,
			})
		}
	}

	for sysName, teamNode := range hrTeamNodes {
		if pos, ok := positionNodes[sysName+"_lead"]; ok {
			builder.AddPositionNodeRelation(&orgtree.PositionNodeRelation{
				ID:       uuid.New(),
				NodeID:   teamNode.ID,
				Position: pos,
			})
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
	devSubtree, err = orgTree.FilterSubtreeByRegex(devPattern.String())
	if err != nil {
		log.Printf("Ошибка при поиске команд разработки: %v", err)
	} else if devSubtree != nil {
		devJSON, _ := devSubtree.ToJSON()
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
		teamJSON, _ := teamSubtree.ToJSON()
		fmt.Println("Найденные команды:")
		printJSON(teamJSON)
	}
}
