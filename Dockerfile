FROM golang:alpine AS build-stage
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /runescape-exporter
FROM gcr.io/distroless/base-debian11:latest
COPY --from=build-stage /runescape-exporter /runescape-exporter
USER nonroot:nonroot
ENTRYPOINT [ "/runescape-exporter" ]