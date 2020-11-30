FROM golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go build -o /go/bin/h24 ./cmd/

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/h24 /
ENTRYPOINT ["/h24"]