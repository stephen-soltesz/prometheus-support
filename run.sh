#!/bin/bash

docker run -p 80:9090 -d soltesz-demo-prometheus

# run Grafana
# sudo docker run -d -p 3000:3000 grafana/grafana
