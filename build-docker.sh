#!/bin/bash

# Change to the directory of this script.
# http://stackoverflow.com/a/246128
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

docker build -t web -f Dockerfile.web .

