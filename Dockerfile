FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/app

FROM alpine:3.20
RUN adduser -D -H appuser
WORKDIR /app
COPY --from=builder /bin/app /app/app
COPY config.yaml /app/config.yaml
USER appuser
EXPOSE 8080
CMD ["/app/app", "-config", "/app/config.yaml"]
