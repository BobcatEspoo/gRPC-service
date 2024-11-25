# Имя бинарного файла
BINARY_NAME=rates

# Имя Docker-образа
DOCKER_IMAGE=grpcbank-app

# Путь к Dockerfile
DOCKERFILE_PATH=./RatesMicroservice/Dockerfile

# Путь к исходному коду
SRC_PATH=./RatesMicroservice

# Сборка приложения
build:
	go build -o $(BINARY_NAME) ./RatesMicroservice/main.go

# Запуск unit-тестов
test:
	go test $(SRC_PATH) -v -cover

# Сборка Docker-образа
docker-build:
	docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE_PATH) .

# Запуск приложения
run: build
	./$(BINARY_NAME)

# Запуск линтера (golangci-lint)
lint:
	golangci-lint run

# Удаление скомпилированного бинарного файла
clean:
	rm -f $(BINARY_NAME)

# Хелп для вывода всех доступных команд
help:
	@echo "Доступные команды:"
	@echo "  make build          - Сборка приложения"
	@echo "  make test           - Запуск unit-тестов"
	@echo "  make docker-build   - Сборка Docker-образа"
	@echo "  make run            - Запуск приложения"
	@echo "  make lint           - Запуск линтера"
	@echo "  make clean          - Очистка проекта"
