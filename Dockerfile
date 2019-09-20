FROM ubuntu:18.04

MAINTAINER tnek

RUN apt-get update
RUN apt-get install -y firefox python3 python3-pip

RUN pip3 install -U pip
RUn pip3 install -r requirements.txt


