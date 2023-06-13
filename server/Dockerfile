FROM golang:1.20.5 AS builder
ARG BUILD_VERSION
ARG BUILD_DATE
RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crosshairs -v -trimpath -ldflags="-s -w -X github.com/devusSs/crosshairs/updater.BuildVersion=${BUILD_VERSION} -X github.com/devusSs/crosshairs/updater.BuildDate=$(BUILD_DATE) -X github.com/devusSs/crosshairs/updater.BuildMode=dev" ./...

FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/docker.env .
COPY --from=builder /app/files/sample.png ./files/sample.png
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/crosshairs ./
CMD ["./crosshairs", "-docker"]