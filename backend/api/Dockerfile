FROM golang:1.26.1 AS builder
WORKDIR /usr/src/app
COPY backend/api/go.mod ./
COPY ackend/api/go.sum ./
RUN go mod download
COPY ackend/api/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/.env .
EXPOSE 5005
CMD ["./main"]