# Build Archlinux packages with drone
#
#     docker build --rm=true -t mikkeloscar/drone-pkgbuild .

FROM gliderlabs/alpine:3.2
MAINTAINER Mikkel Oscar Lyderik <mikkeloscar@gmail.com>

RUN apk-install python3
RUN mkdir -p /opt/drone
COPY requirements.txt /opt/drone/
WORKDIR /opt/drone
RUN pip3 install -r requirements.txt
COPY plugin /opt/drone/

ENTRYPOINT ["python3", "/opt/drone/plugin/main.py"]
