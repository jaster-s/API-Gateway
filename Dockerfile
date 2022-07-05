FROM golang:1.17-bullseye

WORKDIR /usr/src/app

COPY . .

RUN mkdir -p assets/resize

RUN go get -d -v
RUN go build -o fli-gateway-api .

CMD ["./fli-gateway-api"]