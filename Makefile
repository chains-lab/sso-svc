DB_URL=postgresql://postgres:postgres@localhost:7000/postgres?sslmode=disable
OPENAPI_GENERATOR := java -jar ~/openapi-generator-cli.jar
CONFIG_FILE := ./config_local.yaml

migrate-up:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./main ./main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./main migrate up

migrate-down:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./main ./main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./main migrate down

run-server:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./main ./main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./main run service

docker-uo:
	docker-compose up -d

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose up -d --build --force-recreate

swagger-docs:
	cd docs && npm run start