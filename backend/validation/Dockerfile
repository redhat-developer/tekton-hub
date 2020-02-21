FROM golang:1.13.0-stretch

RUN apt-get update && apt-get install -y \
    python3-pip

WORKDIR /app
COPY requirements.txt ./
RUN pip3 install -r requirements.txt

USER 1000
