FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go get -u github.com/jackc/tern

COPY . ./

RUN apk add curl
RUN apk add bash
RUN curl -sfL $(curl -s https://api.github.com/repos/powerman/dockerize/releases/latest | grep -i /dockerize-$(uname -s)-$(uname -m)\" | cut -d\" -f4) | install /dev/stdin /usr/local/bin/dockerize

RUN chmod 755 ./run.sh

CMD /usr/local/bin/dockerize -wait tcp://bc_db:5432 -wait tcp://bc_redis:6379 /bin/bash ./run.sh
