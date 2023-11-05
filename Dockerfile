FROM golang:1.20.6-bullseye as builder

WORKDIR /core
COPY ./ /core

RUN wget https://github.com/WebAssembly/wabt/releases/download/1.0.17/wabt-1.0.17-ubuntu.tar.gz
RUN tar -zxf wabt-1.0.17-ubuntu.tar.gz
RUN cp wabt-1.0.17/bin/wat2wasm /usr/local/bin

RUN make install && make faucet

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

RUN apt-get update && apt-get install -y wget

COPY ./docker-config/validator1/ validator1/
COPY ./docker-config/validator2/ validator2/
COPY ./docker-config/validator3/ validator3/
COPY ./docker-config/validator4/ validator4/

# generated genesis
COPY ./docker-config/genesis.json genesis.json

CMD ["bandd", "--help"]
