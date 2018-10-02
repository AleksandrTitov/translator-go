FROM golang:alpine AS build-env

WORKDIR /app
ADD . /app
RUN cd /app && go build -o translator

FROM alpine
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/translator /app

EXPOSE 8080
CMD ["/app/translator"]
