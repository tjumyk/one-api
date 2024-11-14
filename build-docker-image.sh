#!/bin/bash

set -e

sudo docker image pull node:16
sudo docker image pull golang:alpine
sudo docker image pull alpine

sudo docker build . -t one-api-lmzgc
rm -f one-api-lmzgc.tar.gz && sudo docker image save one-api-lmzgc | gzip > one-api-lmzgc.tar.gz
