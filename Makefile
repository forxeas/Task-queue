include .env
export

migrate up:
	migrate -source ./migrate -database ${DATABASE_URL} up

migrate down:
	migrate -source ./migrate -database ${DATABASE_URL} down