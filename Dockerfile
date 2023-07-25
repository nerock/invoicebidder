FROM golang:1.20.6-alpine3.18 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/bin/ibidder cmd/ibidder.go

FROM alpine:3.18 AS runtime

WORKDIR /

COPY --from=build /app/bin /bin
COPY --from=build /app/config.json /config.json

ENTRYPOINT ["./bin/ibidder"]