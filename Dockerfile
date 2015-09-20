FROM golang:1.5

COPY . /go/src/github.com/docker/swarm
WORKDIR /go/src/github.com/docker/swarm

ENV GOPATH /go/src/github.com/docker/swarm/Godeps/_workspace:$GOPATH
RUN CGO_ENABLED=0 go install -v -a -tags netgo -installsuffix netgo -ldflags "-w -X github.com/docker/swarm/version.GITCOMMIT `git rev-parse --short HEAD`"

ENV SWARM_HOST :2375
EXPOSE 2375

VOLUME $HOME/.swarm

CMD ["swarm"]
CMD ["--help"]

ENV SERF_VERSION 0.6.4_linux_amd64

RUN apt-get update
RUN apt-get install -y unzip wget


RUN curl -L https://dl.bintray.com/mitchellh/serf/$SERF_VERSION.zip -o serf.zip


RUN unzip serf.zip

RUN mv serf /usr/local/bin

RUN rm serf.zip

