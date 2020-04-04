FROM golang:1.14-alpine as builder
WORKDIR /app
COPY ./go.* ./
RUN go mod download
COPY ./go ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/app /bin/app
EXPOSE 8080
ENTRYPOINT ["/bin/app"]
