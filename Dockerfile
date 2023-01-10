FROM golang:1.19.4-buster as build

WORKDIR /go/src/app
COPY go.mod .
COPY . .

RUN go mod download

RUN cd /go/src/app/cmd/agent && go build -ldflags="-s" -o /go/bin/app

FROM frolvlad/alpine-glibc

COPY --from=build /go/bin/app /
CMD ["/app"]