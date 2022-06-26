##
## Build
##

FROM golang:1.18-alpine AS build

RUN apk add --no-cache --update make

WORKDIR /src

COPY ./* ./

RUN go mod download
RUN go build -o ./bin/hwbot ./


##
## Deploy
##

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin ./

ENTRYPOINT [ "./hwbot" ]
