.PHONY: build test clean help run

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

clean: ## Очистить скомпилированные файлы
	go clean
	rm -rf $(BIN_DIR)

# Дополнительные команды
fmt: ## Отформатировать код
	go fmt ./...

vet: ## Проверить код на потенциальные ошибки
	go vet ./...

lint: ## Запустить линтер
	golangci-lint run

# Полный цикл
all: fmt vet test build ## Выполнить полный цикл: форматирование, проверку, тесты и сборку 