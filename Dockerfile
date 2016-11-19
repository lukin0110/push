FROM golang:1.7.3-wheezy

# Packaged dependencies
RUN apt-get update && \
    apt-get install -y \
	curl \
	git \
	tar && \
	rm -rf /var/lib/apt/lists/*

# Package manager
RUN go get -u github.com/kardianos/govendor

# Set workdir
WORKDIR /go/src/github.com/lukin0110/push

# Add the entrypoint.sh
COPY contrib/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod ugo+x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

# Copy the source
COPY . /go/src/github.com/lukin0110/push

# Run bash by default
CMD ["bash"]
