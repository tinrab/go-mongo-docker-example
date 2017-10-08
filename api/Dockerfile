FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go get github.com/pilu/fresh
RUN go-wrapper download
RUN go-wrapper install

CMD [ "fresh" ]
