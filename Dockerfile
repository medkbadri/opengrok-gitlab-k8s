FROM golang:latest as build
WORKDIR /go/src/app
COPY . .
RUN go build -o opengrok-gitlab

FROM debian:stable-slim
COPY --from=build /go/src/app/opengrok-gitlab /usr/local/bin/opengrok-gitlab
COPY --from=build /go/src/app/scripts/start-sync.sh /scripts/start-sync.sh
RUN apt-get update && apt-get install -y git && apt-get install -y inotify-tools &&  \
    mkdir -p /opengrok/src && \
    chmod +x /scripts/start-sync.sh

CMD ["/scripts/start-sync.sh"]