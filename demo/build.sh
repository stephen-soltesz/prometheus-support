#!/bin/bash

sudo docker build -t soltesz-demo-sample .
sudo docker tag soltesz-demo-sample soltesz/demo-client-sample:latest
sudo docker push soltesz/demo-client-sample:latest

