FROM ubuntu:latest
RUN apt-get update

RUN apt-get install -y golang nodejs npm mongodb zip git vim