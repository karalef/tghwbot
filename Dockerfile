##
## Build
##

FROM golang:1.23.2 AS build

WORKDIR /go/src/tghwbot

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/hwbot 


##
## Deploy
##

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/hwbot /

ENTRYPOINT [ "/hwbot" ]
