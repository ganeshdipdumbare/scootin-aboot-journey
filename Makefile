all : start stop test
.PHONY : all

start: 
	docker compose -f docker-compose.yml up --build -d
stop:
	docker compose down
test:
	go test -v ./... 