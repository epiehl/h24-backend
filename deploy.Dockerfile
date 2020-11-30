FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go build -o /go/bin/server ./cmd/server
RUN go build -o /go/bin/notificator ./cmd/notificator
RUN go build -o /go/bin/aggregator ./cmd/aggregator

# Now copy it into our base image.
ENV RUN_APP "server”

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/server /
CMD ["/$RUN_APP"]