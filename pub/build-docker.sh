#!/bin/bash

# Change to the directory of above this script.
# http://stackoverflow.com/a/246128
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR
cd ..

docker build -t pub -f Dockerfile.pub .

