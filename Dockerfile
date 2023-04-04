# Docker image versions
ARG go_ver=v1.19.7-alpine3.17.2
# Docker images
ARG go_img=ghcr.io/dopos/golang-alpine

FROM ${go_img}:${go_ver} as builder

RUN apk add --no-cache git curl

WORKDIR /build

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-X main.version=`git describe --tags --always`" -a ./cmd/hocon

FROM scratch

WORKDIR /
COPY --from=builder /build/hocon .

# SSL support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
ENTRYPOINT ["/hocon"]
