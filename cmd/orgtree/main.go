package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arsants/orgtree"
)

func main() {
	// Создаем новое дерево
	root := orgtree.NewNode("CEO")

	// Добавляем несколько узлов для демонстрации
	root.Value = "CEO"

	// Создаем поддерево
	cto := orgtree.NewNode("CTO")
	root.AddChild(cto)

	dev1 := orgtree.NewNode("Senior Developer")
	dev2 := orgtree.NewNode("Junior Developer")
	cto.AddChild(dev1)
	cto.AddChild(dev2)

	// Сериализуем дерево в JSON
	jsonData, err := root.ToJSON()
	if err != nil {
		log.Fatalf("Ошибка при сериализации: %v", err)
	}

	// Выводим JSON в stdout
	fmt.Println("Структура организации:")
	fmt.Println(string(jsonData))

	// Выводим хеш дерева
	fmt.Printf("\nХеш дерева: %s\n", root.HashString())

	// Сохраняем JSON в файл
	err = os.WriteFile("orgtree.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Ошибка при сохранении файла: %v", err)
	}
	fmt.Println("\nСтруктура сохранена в файл orgtree.json")
}
