# This file describes the standard way to build the ipchecker
#

FROM debian:jessie

# Packaged dependencies
RUN apt-get update && apt-get install -y \
	curl \
	tar

# Install Go
# IMPORTANT: If the version of Go is updated, the Windows to Linux CI machines
#            will need updating, to avoid errors. Ping #docker-maintainers on IRC
#            with a heads-up.
ENV GO_VERSION 1.7.3
RUN curl -fsSL "https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz" \
	| tar -xzC /usr/local

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

# Set workdir
WORKDIR /go/src/github.com/lukin0110/push

# Upload ipchecker source
COPY . /go/src/github.com/lukin0110/push
