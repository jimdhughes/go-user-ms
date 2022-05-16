# Go-Users-Ms

## About

This is the repo to accompany my blog at: https://blog.jimdhughes.com/go-micro-services-and-embedded-databases/ and https://blog.jimdhughes.com/go-micro-services-and-embedded-databases-part-2/

It was built as a way for me to play with Go, Docker, and BoltDB - an embedded key/value store.

## Use

Clone the repo and navigate to the repo
`go build`

`./go-user-ms`

By default we will set up an http server. to change this, set the environment variable `USERMS_SERVER_TYPE=[grpc|http]`

## Docker

Make a docker image using `make docker` investigate the makefile for more info

## Running on docker

`make docker`

`docker run -itd go-user-ms`
