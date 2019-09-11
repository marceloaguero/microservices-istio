################################
# STEP 1 build executable binary
################################
FROM golang:1.13-alpine AS builder

# Git required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /go/src/app
COPY . .
RUN go get -d -v

# Build the binary - remove debug info and compile only for linux target
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /go/bin/app .

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy static executable
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]