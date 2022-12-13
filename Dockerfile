FROM golang

WORKDIR /service

COPY . .

RUN cd cmd/shortener/ && go build -o main .
CMD cmd/shortener/main