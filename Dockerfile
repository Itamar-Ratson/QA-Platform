FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o qa-test-app cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates terraform
WORKDIR /root/
COPY --from=builder /app/qa-test-app .
COPY --from=builder /app/test-cases ./test-cases
COPY --from=builder /app/terraform ./terraform
CMD ["./qa-test-app"]
