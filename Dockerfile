ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -o ./app .

FROM gcr.io/distroless/base:latest

COPY --from=builder /build/app /app

ENV OUTPUT_DIRECTORY=/tmp/records

CMD ["/app"]
