FROM golang:1.17.2-alpine3.14 as builder


LABEL maintainer="Carlos Yeric Fonseca Rios <yeric17@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o yeric-blog .

FROM alpine:3.14

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/yeric-blog .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD [ "./yeric-blog" ]