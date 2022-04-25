FROM golang:1.16 AS builder
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /k8s-deploy .

FROM centos
LABEL Version=0.0.1
LABEL Name=k8s-deploy
LABEL maintainer="sunshuo <sunshuo@haoduo.vip>"

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo Asia/Shanghai > /etc/timezone

RUN mkdir -p /opt/server/config/.kube
COPY ./config/.kube/dev-k8s-config /opt/server/config/.kube/dev-k8s-config
COPY ./config/.kube/qa-k8s-config /opt/server/config/.kube/qa-k8s-config
COPY ./config/.kube/prod-k8s-config /opt/server/config/.kube/prod-k8s-config

WORKDIR /opt/server

# 拷贝可执行文件
COPY --from=builder /k8s-deploy /opt/server/k8s-deploy

EXPOSE 8080

CMD ["./k8s-deploy"]