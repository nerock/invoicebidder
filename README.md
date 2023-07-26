# Invoice Bidder

Publish invoices and bid on them

## Try it
Hosted in my [droplet](ib.manueladalid.dev/swagger/index.html) route (careful please cheapest DO tier)

## Design choices
Although unnecessary for the scope the application is designed as a modular monolith,
with the idea of showing a design pattern oriented to individual services
It makes it easy to further separate them in said services and have them communicate through rpc calls or other communication methods

Although each service has its own responsibilities the idea is that the orchestrator overviews the collaboration between so individual services/domains
do not need each other

There's also included a "fake" event broker for asynchronous operations that could be a real broker like Kafka or RabbitMQ although
a broker could probably be better utilized from the inner services by them publishing their actions and other services reacting to them
instead of through the orchestrator

There are other things about the design and its potential pitfalls to be discussed

## Migrations
- Install [tool](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- Update connection details in Makefile
- Run the migrations either individually or together with `migrate-all`

## How to run locally
- Go run with `make run`
- Build with `make build`

## How to run in docker
- Build image with `make build-docker`
- Run image with `make run-docker`

## Api docs
- See at `localhost:8080/swagger/index.html`
- Regenerate with `make gen-docs`

## Testing
Due to time lmit and not wanting to sacrifice showing other stuff only issuer service and router are tested as examples
