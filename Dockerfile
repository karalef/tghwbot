##
## Build
##

FROM golang:1.19-alpine AS build

WORKDIR /src

COPY . ./

RUN go build -o ./bin/hwbot ./


##
## Deploy
##

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin/hwbot ./

ENTRYPOINT [ "./hwbot" ]
