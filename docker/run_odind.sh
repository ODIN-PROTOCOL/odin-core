#!/bin/bash

if test -n "$1"; then
    # need -R not -r to copy hidden files
    cp -R "$1/.odin" /root
fi

perl -i -pe 's/^minimum-gas-prices = .+?$/minimum-gas-prices = "0.0125loki"/' /root/.odin/config/app.toml

mkdir -p /root/log
odind start --rpc.laddr tcp://0.0.0.0:26657 --trace
