FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY api api
COPY biz biz
COPY internal internal
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -o main main.go

FROM scratch AS runner

COPY ./certs ./certs
COPY ./public ./public
COPY ./templates ./templates
COPY ./configs.json ./configs.json
COPY --from=builder /build/main main
CMD [ "./main" ]