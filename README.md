# OrgTree

OrgTree - это библиотека на Go для работы с организационными структурами в виде деревьев. Она позволяет создавать, манипулировать и сериализовать иерархические структуры данных.

## Особенности

- Создание и манипуляция древовидными структурами
- Сериализация в JSON и десериализация из JSON
- Вычисление хеша структуры для отслеживания изменений
- Фильтрация поддеревьев по условию или регулярному выражению
- Обход дерева в ширину (BFS) и глубину (DFS)
- Простой и понятный API

## Требования

- Go 1.16 или выше
- Поддерживаемые операционные системы:
  - Linux
  - macOS
  - Windows

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

### Получение пути и глубины

```go
// Получение пути к узлу
path := node.GetPath(root)
for _, n := range path {
    fmt.Println(n.Value)
}

// Получение глубины узла
depth, ok := node.GetDepth(root)
if ok {
    fmt.Printf("Глубина узла: %d\n", depth)
}
```

### Работа с должностями

```go
// Создание должностей
position := orgtree.Position{
    ID:      uuid.New(),
    Name:    "Senior Developer",
    Sysname: "senior_developer",
}

// Создание связи между должностью и узлом
relation := orgtree.PositionNodeRelation{
    ID:         uuid.New(),
    NodeID:     cto.ID,
    PositionID: position.ID,
}
```

### Получение связанных узлов

```go
// Получение всех узлов, связанных с текущим
connectedNodes := ceo.GetConnectedNodes()

// Получение всех должностей для узла
positions := cto.GetPositions()
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

### Фильтрация поддерева

```go
// Фильтрация по условию
filteredTree := root.FilterSubtree(func(value interface{}) bool {
    if str, ok := value.(string); ok {
        return strings.Contains(str, "Senior")
    }
    return false
})

// Фильтрация по регулярному выражению
pattern := regexp.MustCompile("Senior")
filteredTree, err := root.FilterSubtreeByRegex(pattern)
if err != nil {
    log.Fatal(err)
}
```

### Обход дерева

```go
// Установка типа обхода (BFS по умолчанию)
root.SetTraversalType(orgtree.BFS)

// Обход дерева
for node := root; node != nil; node = node.Next() {
    fmt.Println(node.Value)
}
```

## Структура проекта

```
orgtree/
├── cmd/                  # Директория с исполняемыми файлами
│   └── orgtree/         # Основной исполняемый файл
│       └── main.go      # Точка входа приложения
├── binaries/            # Директория для скомпилированных бинарных файлов
├── node.go              # Основные структуры и интерфейсы
├── node_utils.go        # Вспомогательные функции для работы с узлами
├── iterator.go          # Реализация итераторов для обхода дерева
├── filter.go            # Функции фильтрации дерева
├── node_test.go         # Тесты базовых операций с узлами
├── filter_test.go       # Тесты фильтрации
├── hash_test.go         # Тесты хеширования
├── json_test.go         # Тесты сериализации
├── orgtree.json         # Пример JSON-структуры
├── go.mod               # Go модуль
├── Makefile             # Команды для сборки и тестирования
├── .gitignore          # Игнорируемые Git файлы
└── README.md            # Документация
```

## Сборка и тестирование

Проект использует Makefile для автоматизации задач:

```bash
make help      # Показать список доступных команд
make build     # Собрать проект
make run       # Собрать и запустить приложение
make test      # Запустить тесты
make test-all  # Запустить все тесты с генерацией отчета о покрытии
make clean     # Очистить скомпилированные файлы
make fmt       # Отформатировать код
make vet       # Проверить код на потенциальные ошибки
make lint      # Запустить линтер
make all       # Выполнить полный цикл: форматирование, проверку, тесты и сборку
```

Все скомпилированные файлы сохраняются в директории `binaries/`.

## Ограничения

- Максимальная глубина дерева ограничена только доступной памятью
- Сериализация поддерживает только JSON формат
- Значения узлов должны быть сериализуемыми в JSON

## Планы развития

- [ ] Добавление поддержки XML сериализации
- [ ] Реализация параллельного обхода дерева
- [ ] Добавление поддержки версионирования структуры
- [ ] Интеграция с базами данных
- [ ] Добавление веб-интерфейса для визуализации

## Лицензия

MIT 