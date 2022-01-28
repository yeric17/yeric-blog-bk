# FROM golang:1.17.2-alpine3.14 as builder


# LABEL maintainer="Carlos Yeric Fonseca Rios <yeric17@gmail.com>"

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . /app/

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o yeric-blog .

# FROM alpine:3.14

# RUN apk add --no-cache ca-certificates && update-ca-certificates

# COPY --from=builder /app/yeric-blog .
# # COPY --from=builder /app/.env .
# COPY --from=builder /app/images ./images
# COPY --from=builder /app/email ./email

# EXPOSE 7070 7070

# CMD [ "./yeric-blog" ]

FROM heroku/heroku:20-build as build

COPY . /app
WORKDIR /app

# Setup buildpack
RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://buildpack-registry.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go

#Execute Buildpack
RUN STACK=heroku-20 /tmp/buildpack/heroku/go/bin/compile /app /tmp/build_cache /tmp/env

# Prepare final, minimal image
FROM heroku/heroku:20

COPY --from=build /app /app
ENV HOME /app
WORKDIR /app
RUN useradd -m heroku
USER heroku
CMD /app/bin/go-getting-started
