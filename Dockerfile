FROM golang:1.17.2-alpine3.14 as builder


LABEL maintainer="Carlos Yeric Fonseca Rios <yeric17@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o yeric-blog .

FROM alpine:3.14

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /app/yeric-blog .
# COPY --from=builder /app/.env .
COPY --from=builder /app/images ./images
COPY --from=builder /app/email ./email

EXPOSE 7070 7070

CMD [ "./yeric-blog" ]