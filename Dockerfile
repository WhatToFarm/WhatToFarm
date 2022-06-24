FROM golang:1.17.3-alpine3.15 as Builder
COPY . /go/src/tg-bot
WORKDIR /go/src/tg-bot
RUN go mod tidy -compat=1.17
RUN go build -o /go/src/tg-bot/bin/bot

FROM alpine:3.15
COPY --from=Builder /go/src/tg-bot/bin/* /go/src/tg-bot/