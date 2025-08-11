# Параметры
COMPOSE_FILE=docker-compose.yml

# Запуск сервиса из контейнера
docker-up:
	docker-compose -f $(COMPOSE_FILE) up -d

# Остановка и удаление контейнеров
docker-down:
	docker-compose -f $(COMPOSE_FILE) down

# Остановка контейнеров
docker-stop:
	docker-compose -f $(COMPOSE_FILE) stop

# Запуск контейнеров
docker-start:
	docker-compose -f $(COMPOSE_FILE) start

# Запуск сервиса локально
run:
	go run ./cmd/main.go

# Запуск unit-тестов локально
test:
	go test -count=1 ./... -v

# Обновление swagger документации
swagger:
	swag init -g cmd/main.go
