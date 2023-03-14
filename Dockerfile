FROM golang:1.19-alpine AS builder

ARG GOOS=linux
ARG GOARCH=amd64

COPY . /go/src/leaderElectionSample
WORKDIR /go/src/leaderElectionSample/cmd
RUN go build -o bin/leaderelectionsample
ENTRYPOINT ["bin/leaderelectionsample"]