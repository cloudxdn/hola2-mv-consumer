FROM golang:1.22.1

WORKDIR /root

COPY hola2-mv-consumer.tar .
RUN tar xvf hola2-mv-consumer.tar

WORKDIR /root/hola2-mv-consumer

RUN go mod vendor
RUN go build -o hola2-mv-consumer main.go

EXPOSE 7443

ENTRYPOINT [ "sh", "-c" ]
CMD [ "./hola2-mv-consumer" ]