FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build
RUN go install github.com/swaggo/swag/cmd/swag@v1.5.0
RUN swag init
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=direct go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/migration /app/migration/

WORKDIR /app
EXPOSE 8080
CMD ["./main"]