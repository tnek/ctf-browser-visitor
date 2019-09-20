FROM ubuntu:18.04

MAINTAINER tnek

RUN apt-get update && apt-get install -y firefox python3 python3-pip

RUN pip3 install -U pip
RUN pip3 install -r requirements.txt

COPY geckodriver /usr/local/bin

RUN groupadd -g 1000 app
RUN useradd -g app -m -u 1000 app -s /bin/bash
USER app

COPY src ./src
COPY requirements.txt ./
RUN pip install -r requirements.txt
WORKDIR ./src

EXPOSE 5000
CMD ["hypercorn", "-w", "40", "-b", "0.0.0.0:5000", "app:app"]
