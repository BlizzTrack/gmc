FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN go build -o gmc ./cmd/gmc/main.go

FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
COPY --from=builder /build/gmc /app/
WORKDIR /app

EXPOSE 11211

CMD ["./gmc"]
