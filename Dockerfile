FROM golang:1.10 as build

WORKDIR /go/src/github.com/DanShu93/martialarts-tracker

COPY . .

RUN go get ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app cmd/server/server.go

FROM alpine:3.8

EXPOSE 80

COPY --from=build /app /go/bin/app

ENTRYPOINT /go/bin/app
