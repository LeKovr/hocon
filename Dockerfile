# Docker image versions
ARG go_ver=v1.18.5-alpine3.16.2

# Docker images
ARG go_img=ghcr.io/dopos/golang-alpine

FROM ${go_img}:${go_ver}

#RUN apk add --no-cache git curl

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=`git describe --tags --always`" -a ./cmd/hocon

FROM scratch

WORKDIR /
COPY --from=0 /app/hocon .
# Need for SSL
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
ENTRYPOINT ["/hocon"]
