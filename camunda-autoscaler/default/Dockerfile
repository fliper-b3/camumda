FROM golang:alpine as build

ADD ./src/main.go /src/main.go

WORKDIR /src
RUN apk add git 
RUN go get k8s.io/apimachinery/pkg/apis/meta/v1
RUN go get k8s.io/client-go/kubernetes
RUN go get k8s.io/client-go/rest
RUN go get k8s.io/client-go/tools/clientcmd
RUN go get github.com/sirupsen/logrus

RUN env GOOS=linux GOARCH=amd64 go build -o autoscaler

FROM golang:alpine

COPY --from=build /src/autoscaler /usr/bin/autoscaler

CMD ["/usr/bin/autoscaler"]
