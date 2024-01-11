FROM golang:1.20.6-bullseye as builder

WORKDIR /core
COPY ./ /core

RUN apt-get update && \
    apt-get install -y ca-certificates wget libc6-dev && \
    update-ca-certificates && \
    wget https://github.com/WebAssembly/wabt/releases/download/1.0.17/wabt-1.0.17-ubuntu.tar.gz && \
    tar -zxf wabt-1.0.17-ubuntu.tar.gz && \
    cp wabt-1.0.17/bin/wat2wasm /usr/local/bin && \
    make install && \
    make faucet && rm -rf /core/* \
    && rm -rf /var/lib/apt/lists/*

CMD ["odind", "--help"]