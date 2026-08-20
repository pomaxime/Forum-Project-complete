# ── Stage 1 : compilation ──────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

# Dépendances système nécessaires à go-sqlite3 (cgo)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -o forum ./main.go

FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/forum .

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static    ./static

EXPOSE 8080

RUN addgroup -S forum && adduser -S forum -G forum
USER forum

CMD ["./forum"]
