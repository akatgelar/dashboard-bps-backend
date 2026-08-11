# syntax=docker/dockerfile:1

############################
# Stage 1: build binary
############################
FROM golang:1.24 AS builder

WORKDIR /app

# deps first (cache friendly)
COPY go.mod go.sum ./
RUN go mod download

# source
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server main.go

############################
# Stage 2: runtime
############################
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Jakarta
WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
