FROM golang:1.13-alpine AS build-platform

WORKDIR /app
COPY . .

RUN apk add --no-cache make
RUN make build

FROM alpine:3.12

LABEL maintainer="Elia Mazzuoli <zikoel@gmail.com>"
WORKDIR /app

COPY --from=build-platform /app/bin/shortener /app/shortener

EXPOSE 5000

CMD ["/app/shortener"]