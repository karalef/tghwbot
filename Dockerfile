##
## Build
##

FROM golang:1.18-alpine AS build

WORKDIR /src

COPY ./* ./

RUN go build -o ./bin/hwbot ./


##
## Deploy
##

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin .

ENTRYPOINT [ "./bin/hwbot" ]
