#

FROM debian:jessie

# Packaged dependencies
RUN apt-get update && \
    apt-get install -y \
	curl \
	git \
	tar && \
	rm -rf /var/lib/apt/lists/*

# Install Go
# IMPORTANT: If the version of Go is updated, the Windows to Linux CI machines
#            will need updating, to avoid errors. Ping #docker-maintainers on IRC
#            with a heads-up.
ENV GO_VERSION 1.7.3
RUN curl -fsSL "https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz" \
	| tar -xzC /usr/local

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

# https://github.com/kardianos/govendor
RUN go get -u github.com/kardianos/govendor

# Set workdir
WORKDIR /go/src/github.com/lukin0110/push

# Add the entrypoint.sh
COPY contrib/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod ugo+x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

# Upload ipchecker source
COPY . /go/src/github.com/lukin0110/push

# Run bash by default
CMD ["bash"]
