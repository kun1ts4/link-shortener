.PHONY: run-postgres run-memory stop clean test

run-postgres:
	STORAGE_TYPE=postgres docker compose up --build

run-memory:
	STORAGE_TYPE=memory docker compose up --build link-shortener

stop:
	docker compose down

clean:
	docker compose down -v

test:
	go test -v ./...

