FROM --platform=$BUILDPLATFORM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH
RUN if [ -n "$TARGETARCH" ]; then \
      CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/backend .; \
    else \
      CGO_ENABLED=0 GOOS=linux go build -o /out/backend .; \
    fi

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/backend ./backend

ENV PORT=3000

EXPOSE 3000

CMD ["./backend"]
