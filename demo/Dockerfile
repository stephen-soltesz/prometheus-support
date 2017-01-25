FROM ubuntu:16.04
MAINTAINER soltesz@google.com
RUN apt-get update
RUN apt-get install -y ca-certificates python-setuptools python-dev build-essential
RUN easy_install pip
RUN pip install prometheus_client
ADD sample.py /sample.py
EXPOSE 8000
ENTRYPOINT ["/sample.py"]
