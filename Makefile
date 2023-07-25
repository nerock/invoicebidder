migrate-invoices:
	migrate \
				-source file://internal/invoice/storage/migrations \
				-database postgres://postgres:postgres@localhost:5432/invoice?sslmode=disable \
				up
migrate-investors:
	migrate \
    		-source file://internal/investor/storage/migrations \
    		-database postgres://postgres:postgres@localhost:5432/investor?sslmode=disable \
    		up
migrate-issuers:
	migrate \
				-source file://internal/issuer/storage/migrations \
				-database postgres://postgres:postgres@localhost:5432/issuer?sslmode=disable \
				up

migrate-all: migrate-investors migrate-issuers migrate-invoices

build-docker:
	docker build -t ibidder

run-docker:
	docker run -dp 8080:8080 ibidder

build:
	go build -o bin/ cmd/ibidder.go

run:
	go run cmd/ibidder.go