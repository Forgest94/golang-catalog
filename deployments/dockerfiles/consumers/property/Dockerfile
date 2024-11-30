FROM golang:alpine AS builder
WORKDIR /consumers

RUN apk add --no-progress --no-cache gcc musl-dev git bash make

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download

COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -tags musl -ldflags '-extldflags "-static"' -o ./bin/propertyConsumer cmd/consumers/property/main.go

FROM alpine AS runner

COPY --from=builder /consumers/bin/propertyConsumer /
COPY --from=builder /consumers/.env /.env

CMD ["/propertyConsumer"]