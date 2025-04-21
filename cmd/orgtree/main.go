package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/arsants/orgtree"
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

	// Устанавливаем тип обхода

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
}
