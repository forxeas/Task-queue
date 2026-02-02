include .env
export

migrate-up:
	migrate -path ./migrations -database ${MIGRATION_URL} up

migrate-down:
	migrate -path ./migrations -database ${MIGRATION_URL} down