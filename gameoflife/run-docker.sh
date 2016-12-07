#!/bin/bash

docker run --publish 8080:8080 \
  --link pub:pub \
  --name web --rm \
  web

