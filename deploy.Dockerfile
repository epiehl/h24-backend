FROM 942928664695.dkr.ecr.eu-central-1.amazonaws.com/golang:1.15-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go build -o /go/bin/h24 ./cmd/

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10:latest
COPY --from=build /go/bin/h24 /
COPY --from=build /go/src/app/assets/swagger /assets/swagger

EXPOSE 3000

ENTRYPOINT ["/h24"]