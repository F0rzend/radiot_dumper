FROM golang:1.18.1 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -buildvcs=false -ldflags "-s -w" -o ./app .

FROM gcr.io/distroless/base:latest

COPY --from=builder /build/app /app
COPY --from=builder /etc/mime.types /etc/mime.types

CMD ["/app"]
