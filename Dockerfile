FROM golang:1.21

ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /server
COPY . /server

RUN go mod tidy
RUN go build -v -o app

# webhook服务端口
EXPOSE 8000
# 业务端口
EXPOSE 9000

ENTRYPOINT ./app