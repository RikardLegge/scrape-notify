FROM golang:latest
WORKDIR /build/
COPY . /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /build/app .
CMD ["./app"]
