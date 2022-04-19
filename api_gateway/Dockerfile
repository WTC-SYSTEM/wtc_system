FROM golang:1.18-alpine3.14 AS builder

WORKDIR /usr/local/go/src/
COPY app .

RUN go build --mod=mod -o app cmd/main/app.go

FROM alpine:3.14
# copy binary app
COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/config.yml /

CMD ["/app"]
