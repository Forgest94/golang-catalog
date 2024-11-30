FROM golang:alpine AS builder
WORKDIR /data

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -tags musl -o ./bin/app cmd/app/main.go

FROM alpine AS runner

COPY --from=builder /data/bin/app /
COPY --from=builder /data/.env /.env

CMD ["/app"]