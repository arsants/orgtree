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
		TypeID:  nodeType.ID,
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
		ID:         uuid.New(),
		NodeID:     node.ID,
		PositionID: position.ID,
		Position:   position,
	}

	// Добавляем данные в построитель
	builder.AddNode(node)

	// Создаем и добавляем дополнительные узлы
	itDept := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "IT отдел",
		SysName: "it_department",
		TypeID:  nodeType.ID,
		Type:    nodeType,
	}
	builder.AddNode(itDept)

	hrDept := &orgtree.OrgNode{
		ID:      uuid.New(),
		Name:    "Отдел кадров",
		SysName: "hr_department",
		TypeID:  nodeType.ID,
		Type:    nodeType,
	}
	builder.AddNode(hrDept)

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

	// Добавляем связи должностей с узлами
	builder.AddPositionNodeRelation(relation)
	builder.AddPositionNodeRelation(&orgtree.PositionNodeRelation{
		ID:         uuid.New(),
		NodeID:     itDept.ID,
		PositionID: itManager.ID,
		Position:   itManager,
	})
	builder.AddPositionNodeRelation(&orgtree.PositionNodeRelation{
		ID:         uuid.New(),
		NodeID:     hrDept.ID,
		PositionID: hrManager.ID,
		Position:   hrManager,
	})

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

}
