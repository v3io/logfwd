FROM golang:1.11 as builder

# copy sources
ADD . /go/src/github.com/v3io/logfwd

# build the logfwd
RUN make -C /go/src/github.com/v3io/logfwd bin

FROM debian:stretch-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

COPY --from=builder /go/bin/logfwd /usr/local/bin/logfwd

CMD [ "logfwd" ]
