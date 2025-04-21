# OrgTree

Проект для работы с организационной структурой.

## Установка

```bash
go get github.com/yourusername/orgtree
```

## Использование

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/orgtree"
)

func main() {
    // Создаем новое дерево
    tree := orgtree.NewTree()

    // Добавляем узлы
    err := tree.AddNode("1", "CEO", "root")
    if err != nil {
        log.Fatal(err)
    }

    // Получаем узел по ID
    node := tree.GetNode("1")
    if node != nil {
        fmt.Printf("Найден узел: %s\n", node.Name)
    }
}
```

Полный пример использования можно найти в директории `cmd/main.go`.

## Запуск примера

```bash
go run cmd/main.go
```

## Лицензия

MIT 