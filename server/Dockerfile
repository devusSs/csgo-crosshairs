FROM golang:1.20.5 AS builder
ARG BUILD_VERSION
ARG BUILD_MODE
RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crosshairs -v -trimpath -ldflags="-s -w -X 'github.com/devusSs/crosshairs/updater.BuildVersion=${BUILD_VERSION}' -X 'github.com/devusSs/crosshairs/updater.BuildDate=$(date)' -X 'github.com/devusSs/crosshairs/updater.BuildMode=${BUILD_MODE}'" ./...

FROM alpine:latest AS production
COPY --from=builder /app/files/sample.png ./files/sample.png
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/crosshairs ./
CMD ["./crosshairs", "-docker"]