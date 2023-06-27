FROM golang:1.20-alpine as builder
WORKDIR /workspace

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-w -s" -o /out/custom_metrics_app_linux .
FROM alpine

WORKDIR /build

COPY --from=builder /out/custom_metrics_app_linux .

EXPOSE 3000
ENTRYPOINT ["./custom_metrics_app_linux"]