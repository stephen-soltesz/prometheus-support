#!/bin/bash

apt-get install -y python3-pip
pip3 install grafanalib

generate-dashboard -o frontend.json frontend.dashboard.py
