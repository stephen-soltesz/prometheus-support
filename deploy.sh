#!/bin/bash

set -ex
kubectl run prometheus --image=soltesz/demo-kubernetes-discover:latest --port=9090
kubectl expose deployment prometheus --type="LoadBalancer"
kubectl get service prometheus

