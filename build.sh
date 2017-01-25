#!/bin/bash

sudo docker build -t soltesz-demo-prometheus .
sudo docker tag soltesz-demo-prometheus soltesz/demo-kubernetes-discover:latest
sudo docker push soltesz/demo-kubernetes-discover:latest

