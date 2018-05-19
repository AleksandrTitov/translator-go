FROM golang

RUN mkdir -p /opt/gocode/translator
ENV GOPATH=/opt/gocode/translator

COPY main.go /opt/gocode/translator

WORKDIR /opt/gocode/translator
RUN go build

EXPOSE 8080

CMD ["/opt/gocode/translator/translator"]