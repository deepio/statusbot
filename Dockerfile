FROM golang:alpine as build

COPY ./statusbot.go ./testbot.go /opt/
RUN cd /opt \
  && go build statusbot.go \
  && go build testbot.go

FROM alpine:latest

ENV PATH=/opt:/usr/bin:/bin/
COPY --from=build /opt/statusbot /opt/testbot /opt/
COPY ./config/test.json /tmp

WORKDIR /opt