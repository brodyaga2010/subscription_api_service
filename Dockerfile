FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o eff_mobile ./cmd/main.go


FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/eff_mobile .
COPY config.yaml . 

EXPOSE 8080

CMD ["./eff_mobile", "-config", "config.yaml"]