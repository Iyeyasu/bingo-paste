FROM alpine:latest
MAINTAINER Ieyasu (ieyasu@protonmail.com)
WORKDIR /bingo
COPY server .
COPY web ./web
ENTRYPOINT ["/bingo/server"]
EXPOSE 80
