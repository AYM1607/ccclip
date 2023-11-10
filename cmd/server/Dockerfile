FROM golang:1.21-alpine AS builder
# Ensure we have a c compiler.
RUN apk add --no-cache build-base ca-certificates fuse3 sqlite

# Install LiteFS for distribute SQLite
COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

COPY cmd/server/litefs.yml /etc/litefs.yml

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags='-s -w' -tags 'linux' -trimpath -o /dist/app ./cmd/server

EXPOSE 3000

ENTRYPOINT ["litefs", "mount"]
