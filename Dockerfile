FROM golang:1.14-alpine AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0

WORKDIR /go/src/github.com/bradhe/hobo
COPY . .
RUN mkdir -p ./bin
RUN go build -a -installsuffix cgo -mod vendor -o ./bin/hobo ./cmd/hobo

FROM alpine:latest
RUN apk update
RUN apk --no-cache add git gcc bind-dev musl-dev ca-certificates tzdata
RUN update-ca-certificates
COPY --from=builder /go/src/github.com/bradhe/hobo/bin/hobo /usr/bin/hobo
CMD ["hobo"]
