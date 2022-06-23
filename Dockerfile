FROM golang:1.17-alpine

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /go/src/tg-bot/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags nethttpomithttp2 -ldflags="-w -s" -o /go/src/tg-bot/bot
