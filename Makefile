DB_HOST = localhost
DB_PORT = 5432
DB_USERNAME = root
DB_PASSWORD = secret
DB_DATABASE = hehe_business
DB_SCHEMA = hehe_business

DB_URL = postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable&search_path=$(DB_SCHEMA)
# ----------------------------- Setup database ---------------------------------
databaseup:
	docker compose -f deployments/docker-compose.yaml up -d

databasedown:
	docker compose -f deployments/docker-compose.yaml down

# ------------------- Read schema sql -> crete or update database --------------
# Migarte database all
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

# Migarte database lastest
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

# ------------------- Read schema and query sqlc -> generate code golang -------
sqlc:
	sqlc generate -f ./etc/sqlc.yaml

# Start server
server:
	go run main.go

proto:
	rm -f rd_portfolio_rpc/*.go
	protoc --proto_path=proto --go_out=rd_portfolio_rpc --go_opt=paths=source_relative \
	--go-grpc_out=rd_portfolio_rpc --go-grpc_opt=paths=source_relative \
	proto/*.proto

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: databaseup databasedown migrateup migratedown migrateup1 migratedown1 sqlc server proto redis
