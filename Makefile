build:
	docker-compose build timeslot-app

run: 
	docker-compose up timeslot-app

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up

swag:
	swag init -g cmd/app/main.go
