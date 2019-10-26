FROM alpine:3.9

ADD jukebox_server /
EXPOSE 8080
ENTRYPOINT /jukebox_server
