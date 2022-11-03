FROM golang:bullseye

WORKDIR /tmp/

RUN set -eux; \
    apt update && apt install -y wget unzip \
    && wget https://codeload.github.com/fxtaoo/qwflow/zip/refs/heads/master -O qwflow.zip \
    && unzip qwflow.zip \
    && cd qwflow-master  \
    && go mod tidy \
    && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o qwflow .


FROM alpine:latest

RUN set -eux; \
    apk add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /app

COPY --chmod=0755 --from=0 /tmp/qwflow-master/qwflow /app/qwflow
COPY template /app/template

EXPOSE 8174

CMD ["./qwflow"]
