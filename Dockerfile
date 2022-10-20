FROM ubuntu:latest
MAINTAINER tnek
RUN apt-get update && apt-get install -y firefox python3 python3-pip curl

RUN VERSION=$(curl -sL https://api.github.com/repos/mozilla/geckodriver/releases/latest | grep tag_name | cut -d '"' -f 4) && curl -sL "https://github.com/mozilla/geckodriver/releases/download/$VERSION/geckodriver-$VERSION-linux-aarch64.tar.gz" | tar -xz -C /usr/local/bin

RUN pip3 install -U pip
COPY requirements.txt ./
RUN pip install -r requirements.txt

RUN groupadd -g 1000 app
RUN useradd -g app -m -u 1000 app -s /bin/bash
USER app

WORKDIR /src
COPY src .

EXPOSE 8080
CMD ["hypercorn", "-w", "40", "-b", "0.0.0.0:8080", "app:app"]
