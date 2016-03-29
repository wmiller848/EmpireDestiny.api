FROM golang:1.5.3

ADD ./src /go/src/ed

WORKDIR /go/src/ed
RUN ls
RUN go build -v
