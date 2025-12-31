FROM golang:1.21-alpine AS builder
WORKDIR /app

# no external deps, so no go.sum, go.mod etc required
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o go-grab-geocode main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/go-grab-geocode /go-grab-geocode
ENTRYPOINT ["/go-grab-geocode"]
