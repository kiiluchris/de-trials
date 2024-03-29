build:
	go build -v -o bin/de .

run: build
	./bin/de

startdb:
	docker compose up db

stopdb:
	docker compose stop db

rmdb:
	docker compose rm db
