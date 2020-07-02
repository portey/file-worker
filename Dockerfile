FROM golang:1.14.4-alpine as builder
WORKDIR /app
COPY . ./
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ./bin/svc

FROM scratch
COPY --from=builder /app/bin/svc /svc
EXPOSE 8080 8888 9100
CMD ["./svc"]
