FROM golang:1.25 AS builder
RUN useradd -u 10001 appuser
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main-app ./cmd/app/main/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate-tool ./cmd/migration/main/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/main-app .
COPY --from=builder /bin/migrate-tool .
USER appuser
EXPOSE 8080
CMD ["./main-app"]