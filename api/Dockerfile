FROM golang:1.12.1

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/app
COPY . .

RUN go get github.com/pilu/fresh
RUN go get ./...

CMD [ "fresh" ]
