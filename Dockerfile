FROM golang:1.8
MAINTAINER cmc - cmc@unallocated.net version: 0.1

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
