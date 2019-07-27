FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY mogutou /mogutou
ENTRYPOINT /mogutou
