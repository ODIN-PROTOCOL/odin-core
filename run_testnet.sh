#!/bin/bash
docker run --rm -it -p 26657:26657 -p 26656:26656 -p 1317:1317 -v ./odin_data:/root geodbodinprotocol/core:v0.7.7-2 /core/run_odind.sh