# Game of Life

This is a [Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)
simulation. The simulation runs on a single machine and publishes state updates
over ZeroMQ to the web frontend servers, which in turn send updates to browsers
via [WebSockets](http://www.html5rocks.com/en/tutorials/websockets/basics/).

The simulation is different in that it uses a 32-bit integer to represent the
state instead of a boolean. There are actually many games of life running in
parallel in this implementation. By playing with the mask used to select which
bits are used to represent the "alive" state, you can make the parallel versions
easier to see or even get them to interact with each other.

My goal is to be able to run other kinds of simulations with this framework. By
having small random programs flow data into each other, we can make some pretty
interesting pixel art. I used the Game of Life as a first example since I
realized it only requires a few very small programs to talk to each other and
can be run in parallel quite easily.

## Demo

https://youtu.be/WBenbd_HZgw

## Building

Use [Docker](http://www.docker.com) to build the server. I upload to
[Google container registry](https://cloud.google.com/container-registry/)
multiple images to make for faster builds when I don't have to update
dependencies. Replace golang-game-of-life with your Google Cloud project name.

Layer 1, Go with C++ build tools:

    % docker build -t golang-cpp -f Dockerfile.golang-cpp .
    % docker tag -f golang-cpp gcr.io/golang-game-of-life/golang-cpp:v1
    % docker tag -f golang-cpp gcr.io/golang-game-of-life/golang-cpp:latest
    % gcloud docker push gcr.io/golang-game-of-life/golang-cpp:latest

Layer 2, ZeroMQ and goczmq:

    % docker build -t goczmq -f Dockerfile.goczmq .
    % docker tag -f goczmq gcr.io/golang-game-of-life/goczmq:v1
    % docker tag -f goczmq gcr.io/golang-game-of-life/goczmq:latest
    % gcloud docker push gcr.io/golang-game-of-life/goczmq:latest

Publisher (simulation server):

    % docker build -t pub -f Dockerfile.pub .
    % docker tag -f pub gcr.io/golang-game-of-life/pub:v1-6

Web server:

    % docker build -t web -f Dockerfile.web .
    % docker tag -f web gcr.io/golang-game-of-life/web:v1-8

## Running locally

After building the Docker images, you can run them locally before pushing to a
cloud service.

    % ./pub/run-docker.sh &
    % ./run-docker.sh &

Go to http://localhost:8080 to see it in action.

To kill when you are done:

    % ./kill-local.sh
    % ./pub/kill-local.sh

## Deploy the server

The server is deployed using Kubernetes.
[Google Container Engine](https://cloud.google.com/container-engine/) is a good
option for running this without much set up.

Make sure all the Docker images are available. Push them to
[Google container registry](https://cloud.google.com/container-registry/)

    % gcloud docker push gcr.io/golang-game-of-life/pub:v1-6
    % gcloud docker push gcr.io/golang-game-of-life/web:v1-8

Create the Kubernetes services and replication controllers.

    % kubectl create -f pub/pub-service.json
    % kubectl create -f pub/pub-rc.json
    % kubectl create -f web-service.json
    % kubectl create -f web-rc.json

To delete (shutdown and remove):

    % kubectl delete -f pub/pub-service.json
    % kubectl delete -f pub/pub-rc.json
    % kubectl delete -f web-service.json
    % kubectl delete -f web-rc.json

To run a rolling update:

    % kubectl rolling-update pub-v1-v6 -f pub/pub-rc.json
    % kubectl rolling-update web-v1-v8 -f web-rc.json
