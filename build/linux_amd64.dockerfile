FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY ./bin/ethlevel_linux/ethlevel_amd64 ethlevel

CMD [ "/app/ethlevel" ]