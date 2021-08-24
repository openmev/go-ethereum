FROM golang:1.16-stretch AS builder

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y -qq apt-utils expect git git-extras software-properties-common \
    inetutils-tools wget ca-certificates build-essential libssl-dev make

ADD . /go/src/github.com/ethereum/go-ethereum
WORKDIR /go/src/github.com/ethereum/go-ethereum

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o

FROM debian:buster-20210816 as final

COPY --from=builder /go/src/github.com/ethereum/go-ethereum/build/bin /usr/local/bin

RUN apt-get update && apt-get install -y -qq iperf3 openssh-server iputils-ping tmux software-properties-common && apt-get clean -y -qq

#RUN add-apt-repository ppa:ethereum/ethereum && apt-get update && \
#    apt-get install -y solc

WORKDIR /
#ENV PATH /go-ethereum/build/bin:${PATH}
EXPOSE 8080 8545 8180 3030 8546

ENTRYPOINT ["/bin/bash"]
#CMD exec $SHELL
