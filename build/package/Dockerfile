FROM alpine:latest
MAINTAINER Ieyasu (ieyasu@protonmail.com)
WORKDIR /bingo
COPY build/server .
COPY web ./web
ENTRYPOINT ["/bingo/server"]
EXPOSE 80
