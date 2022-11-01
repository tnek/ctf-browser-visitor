FROM golang:1.18-alpine as builder
MAINTAINER tnek

RUN apk add --no-cache ca-certificates git

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

ADD keys/deploy_id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

WORKDIR /src
COPY ./go.mod ./go.sum ./
COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /ctf-browser-visitor .

FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl xvfb openjdk-11-jre software-properties-common
RUN add-apt-repository ppa:mozillateam/ppa && apt-get update && apt-get install -y firefox-esr

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /ctf-browser-visitor /ctf-browser-visitor

RUN VERSION=$(curl -sL https://api.github.com/repos/mozilla/geckodriver/releases/latest | grep tag_name | cut -d '"' -f 4) && curl -sL "https://github.com/mozilla/geckodriver/releases/download/$VERSION/geckodriver-$VERSION-linux-aarch64.tar.gz" | tar -xz -C /usr/local/bin

ENV SELENIUM_JAR_ADDR=https://github.com/SeleniumHQ/selenium/releases/download/selenium-3.141.59/selenium-server-standalone-3.141.59.jar
RUN curl -sL $SELENIUM_JAR_ADDR > /usr/local/bin/selenium-server.jar

EXPOSE 3000
CMD ["./ctf-browser-visitor", "--selenium=/usr/local/bin/selenium-server.jar", "--driver=/usr/local/bin/geckodriver", "0.0.0.0", "3000"]
