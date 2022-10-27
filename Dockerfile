FROM golang:bullseye

WORKDIR /tmp/

RUN set -eux; \
    apt update && apt install -y wget unzip \
    && wget https://codeload.github.com/fxtaoo/qwflow/zip/refs/heads/master -O qwflow.zip \
    && unzip qwflow.zip \
    && cd qwflow-master  \
    && go mod tidy \
    && go build -o qwflow .


FROM debian:stable-slim

ENV TZ="Asia/Shanghai"

WORKDIR /app

COPY --chmod=0755 --from=0 /tmp/qwflow-master/qwflow /app/qwflow
COPY template /app/template

EXPOSE 8174

CMD ["./qwflow"]
