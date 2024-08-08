FROM golang:1.22.2


COPY ./scripts/healthcheck.sh /usr/local/bin/healthcheck.sh
RUN chmod +x /usr/local/bin/healthcheck.sh

COPY ./1st /go/bin
EXPOSE 8080

ENTRYPOINT ["/go/bin/1st"]