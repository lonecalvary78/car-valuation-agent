FROM golang:alpine AS builder

WORKDIR /builder
ENV GOTOOLCHAIN=auto
COPY go.mod go.sum /builder/
RUN go mod download

COPY cmd /builder/cmd
COPY internal /builder/internal

ENV CGO_ENABLED=0
RUN go build -o /builder/car-valuer-agent ./cmd/api

FROM debian:trixie-slim
ARG SERVER_PORT="8080"
LABEL image-name=car-valuation-agent
LABEL BUILD_VERSION=1.0.0
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /builder/car-valuer-agent /app/
COPY skills /app/skills
ENV SERVER_PORT=${SERVER_PORT}
EXPOSE ${SERVER_PORT}
WORKDIR /app
ENTRYPOINT [ "/app/car-valuer-agent" ]