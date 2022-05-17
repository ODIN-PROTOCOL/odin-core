FROM golang:1.15.8-buster as builder

WORKDIR /core
COPY ./ /core

RUN wget https://github.com/WebAssembly/wabt/releases/download/1.0.17/wabt-1.0.17-ubuntu.tar.gz
RUN tar -zxf wabt-1.0.17-ubuntu.tar.gz
RUN cp wabt-1.0.17/bin/wat2wasm /usr/local/bin

RUN make build && make faucet


FROM debian:buster-slim

RUN apt update && apt install -y ca-certificates && update-ca-certificates

RUN apt update && apt install -y wget


COPY --from=builder /core/vendor/github.com/bandprotocol/go-owasm/api/ /usr/lib/
COPY --from=builder /core/build/odind /usr/local/bin/odind
COPY --from=builder /core/build/yoda /usr/local/bin/yoda

ENTRYPOINT ["odind"]

#COPY ./docker-config/validator1/ validator1/
#COPY ./docker-config/validator2/ validator2/
#COPY ./docker-config/validator3/ validator3/
#COPY ./docker-config/validator4/ validator4/

# generated genesis
#COPY ./docker-config/genesis.json genesis.json

# CMD ["bandd", "--help"]
