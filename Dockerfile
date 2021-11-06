FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY * ./

RUN go build -o /bc_info

EXPOSE 8080

CMD [ "/bc_info" ]