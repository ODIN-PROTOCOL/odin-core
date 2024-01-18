#!/bin/bash
docker run --rm -it -e PASSWORD=xxxxxxxxx -v ./odin_data:/root geodbodinprotocol/core:v0.7.7-2 /core/setup_odind.sh