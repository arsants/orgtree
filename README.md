# OrgTree

OrgTree - это библиотека на Go для работы с организационными структурами в виде деревьев. Она позволяет создавать, манипулировать и сериализовать иерархические структуры данных.

## Особенности

- Создание и манипуляция древовидными структурами
- Сериализация в JSON и десериализация из JSON
- Вычисление хеша структуры для отслеживания изменений
- Простой и понятный API

## Установка

```bash
go get github.com/arsants/orgtree
```

## Использование

### Создание дерева

```go
import "github.com/arsants/orgtree"

// Создание корневого узла
root := orgtree.NewNode("CEO")

// Добавление дочерних узлов
cto := orgtree.NewNode("CTO")
root.AddChild(cto)

dev1 := orgtree.NewNode("Senior Developer")
dev2 := orgtree.NewNode("Junior Developer")
cto.AddChild(dev1)
cto.AddChild(dev2)
```

### Сериализация в JSON

```go
// Сериализация в JSON
jsonData, err := root.ToJSON()
if err != nil {
    log.Fatal(err)
}

// Десериализация из JSON
tree, err := orgtree.FromJSON(jsonData)
if err != nil {
    log.Fatal(err)
}
```

### Вычисление хеша

```go
// Получение хеша в виде байтов
hash := root.Hash()

// Получение хеша в виде строки
hashString := root.HashString()
```

## Структура проекта

```
orgtree/
├── cmd/
│   └── orgtree/
│       └── main.go       # Пример использования
├── orgtree.go            # Основной код библиотеки
├── orgtree_test.go       # Тесты
├── go.mod                # Go модуль
├── Makefile              # Команды для сборки и тестирования
└── README.md             # Документация
```

## Сборка и тестирование

Проект использует Makefile для автоматизации задач:

```bash
make build    # Собрать проект
make run      # Собрать и запустить приложение
make test     # Запустить тесты
make clean    # Очистить скомпилированные файлы
make fmt      # Отформатировать код
make vet      # Проверить код на потенциальные ошибки
make lint     # Запустить линтер
make all      # Выполнить полный цикл: форматирование, проверку, тесты и сборку
```

## Лицензия

MIT 