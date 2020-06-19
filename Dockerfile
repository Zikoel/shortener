FROM golang:1.13 AS build-platform

WORKDIR /app
COPY . .

RUN make build

FROM golang:1.13-alpine

LABEL maintainer="Elia Mazzuoli <zikoel@gmail.com>"
WORKDIR /app

COPY --from=build-platform /app/bin/shortener /app/shortener

# Read this https://stackoverflow.com/a/35613430/2381099 for understand the next line
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 5000

CMD ["/app/shortener"]