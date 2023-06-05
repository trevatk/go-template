
FROM golang:1.20-alpine as builder

RUN mkdir -p /usr/src/app

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian11

COPY --from=builder /usr/src/app/server /

VOLUME [ "/usr/src/migrations", "/usr/src/sqlite" ]

EXPOSE 8080

CMD [ "./server" ]