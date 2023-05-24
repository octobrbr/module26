FROM golang:latest AS compiling_stage
RUN mkdir -p /go/src/mod26
WORKDIR /go/src/mod26
ADD main.go .
ADD go.mod .
RUN go install .
 
FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="Lutkov<v.lutkoff@gmail.com>"
WORKDIR /root/
COPY --from=compiling_stage /go/bin/mod26 .
ENTRYPOINT ./mod26