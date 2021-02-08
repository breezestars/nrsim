FROM golang:1.15.8-buster

WORKDIR /root/nrsim

COPY . .

RUN go install ./...

ENTRYPOINT ["/go/bin/worker"]
