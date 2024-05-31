FROM golang:1.22.2-bullseye as builder

WORKDIR /core
COPY ./ /core

RUN apt-get update && \
    apt-get install -y ca-certificates wget libc6-dev && \
    update-ca-certificates && \
    make install && \
    make faucet && rm -rf /core/* \
    && rm -rf /var/lib/apt/lists/*

COPY docker/* /core/
RUN chmod +x /core/*.sh

CMD ["odind", "--help"]
