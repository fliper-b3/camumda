FROM golang:alpine as build

ADD ./src/main.go /src/main.go

WORKDIR /src
RUN apk add git
RUN go get github.com/sirupsen/logrus
RUN env GOOS=linux GOARCH=amd64 go build -o metric-collect

FROM golang:alpine

COPY --from=build /src/metric-collect /usr/bin/metric-collect

CMD ["/usr/bin/metric-collect"]
