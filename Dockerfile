FROM golang:1.20.4-bullseye

USER root
WORKDIR /data/app

RUN \
    apt-get update && apt-get install -y curl nano wget telnet net-tools dnsutils iputils-ping redis-tools mariadb-client \
    python3 python3-pip && \
    apt-get install -y locales && sed -ie 's/# zh_CN.UTF-8 UTF-8/zh_CN.UTF-8 UTF-8/g' /etc/locale.gen && locale-gen && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && pip3 install jupyterlab requests && rm -rf ~/.cache/pip

ENV TZ=Asia/Shanghai \
    LANG=zh_CN.UTF-8 \
    LANGUAGE=zh_CN.UTF-8 \
    LC_ALL=zh_CN.UTF-8

RUN curl -LO https://dl.k8s.io/release/v1.26.0/bin/linux/amd64/kubectl && chmod +x kubectl && mv ./kubectl /usr/local/bin/

COPY ./easyai-platform-linux-amd64 /data/app/easyai-platform