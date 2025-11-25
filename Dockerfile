FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o farmers_backend ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

COPY --from=builder /app/farmers_backend .

EXPOSE 6060

CMD ["./farmers_backend"]

