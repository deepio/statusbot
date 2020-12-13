FROM golang:alpine as build

COPY ./statusbot.go ./testbot.go /opt/
RUN cd /opt \
  && apk add --no-cache git \
  && go get gopkg.in/yaml.v2 \
  && go build statusbot.go \
  && go build testbot.go

FROM alpine:latest

ENV PATH=/opt:/usr/bin:/bin/
COPY --from=build /opt/statusbot /opt/testbot /opt/
COPY ./config/test.yaml /tmp
COPY ./config/test.json /tmp

WORKDIR /opt