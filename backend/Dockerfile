ARG GOPROXY=https://goproxy.cn,direct

FROM golang:1.22-alpine AS builder

ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV GO111MODULE=on

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o noams .

FROM alpine:3.19

RUN apk add --no-cache tzdata ca-certificates
RUN mkdir -p /app/logs

WORKDIR /app

COPY --from=builder /build/noams .

EXPOSE 8080

CMD ["./noams"]
