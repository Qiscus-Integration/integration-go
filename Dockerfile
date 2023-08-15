FROM golang:1.19 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o binary cmd/main.go

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/binary .