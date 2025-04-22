.PHONY: build test clean help run test-all

# Переменные
BINARY_NAME=orgtree
GO=go
BIN_DIR=binaries
CMD_DIR=./cmd/orgtree

# Команда по умолчанию
.DEFAULT_GOAL := help

help: ## Показать это сообщение помощи
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Основные команды
build: ## Собрать проект
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

run: build ## Собрать и запустить приложение
	./$(BIN_DIR)/$(BINARY_NAME)

test: ## Запустить тесты
	go test -v ./...

test-all: ## Запустить все тесты с покрытием
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет о покрытии сохранен в coverage.html"

clean: ## Очистить скомпилированные файлы
	go clean
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

# Дополнительные команды
fmt: ## Отформатировать код
	go fmt ./...

vet: ## Проверить код на потенциальные ошибки
	go vet ./...

lint: ## Запустить линтер
	golangci-lint run

bench: ## Запустить тесты производительности
	go test -bench=. -benchmem ./...

# Полный цикл
all: fmt vet test-all build ## Выполнить полный цикл: форматирование, проверку, тесты и сборку 