FROM golang:latest
WORKDIR /code
RUN apt-get update && apt-get install -y wget && \
    echo "deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/ ./" >> /etc/apt/sources.list && \
    wget https://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/Release.key -O- | apt-key add && \
    apt-get install -y libzmq3-dev

COPY entrypoint.sh go.mod ./
COPY main.go.txt ./main.go
RUN go mod tidy && go mod vendor
ENTRYPOINT ["/entrypoint.sh"]
