FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o scheduler .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/scheduler .
COPY web ./web
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
CMD ["./scheduler"]