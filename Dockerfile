FROM alpine:latest AS deploy
RUN apk --no-cache add ca-certificates
COPY store-exporter /
COPY config.sample.toml  /etc/store-exporter/config.toml
VOLUME ["/etc/store-exporter"]
CMD ["./store-exporter", "--config", "/etc/store-exporter/config.toml"]  