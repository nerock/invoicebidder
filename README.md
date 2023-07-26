# invoicebidder

Publish invoices and bid on them

## Migrations
- Install tool `go install github.com/golang-migrate/migrate/v4`
- Update connection details in Makefile
- Run the migrations either individually or together with `migrate-all`

## How to run locally
- Go run with `make run`
- Build with `make build`

## How to run in docker
- Build image with `make build-docker`
- Run image with `make run-docker`

## Api docs
- See at `localhost:8080/swagger/`
- Regenerate with `make gen-docs`