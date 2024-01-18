#!/bin/bash

if test -n "$1"; then
    # need -R not -r to copy hidden files
    cp -R "$1/.odin" /root
fi

mkdir -p /root/log
odind start --rpc.laddr tcp://0.0.0.0:26657 --trace
