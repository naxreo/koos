## Build
FROM golang:alpine AS builder
## idock.daumkakao.io/steve_kim/rhainkoos:v0.1
## naxreo/rhainkoos:v0.1

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY * ./
RUN go mod download
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

## Deploy
FROM scratch
COPY --from=builder /dist/main .
ENTRYPOINT ["/main"]