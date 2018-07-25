FROM golang:alpine

RUN apk add --no-cache git make
RUN adduser -u 1001 -G root -D gouser
ADD src /go/src/app/

WORKDIR /go/src/app
RUN go get -v -t ./... && go build -v

USER 1001
EXPOSE 8000
CMD ./app

