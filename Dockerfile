FROM alpine:3.8

COPY bin/kruise-state-metrics /

USER nobody

ENTRYPOINT ["/kruise-state-metrics", "--port=8080", "--telemetry-port=8081"]

EXPOSE 8080 8081

