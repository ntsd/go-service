lint:
	go fmt ./...

test:
	go test ./...

migrate-up:
	go run ./cmd/migrate/migrate.go -mode up

migrate-down:
	go run ./cmd/migrate/migrate.go -mode down

run:
	go run ./cmd/service/main.go

mockgen:
	mockgen -source ./internal/storage/storage.go -destination ./internal/storage/mock_storage/mock_storage.go
	mockgen -source ./internal/crypto/jwt.go -destination ./internal/crypto/mock_jwt/mock_jwt.go

gen-swagger:
	swag init -g ./cmd/service/main.go -o ./docs

docker-postgres:
	docker-compose up service-postgres

docker-build:
	docker-compose build --no-cache

docker-migrate:
	docker-compose up service-migrate

docker-up:
	docker-compose up service-app

docker-down:
	docker-compose down

docker-remove:
	docker-compose down -v
