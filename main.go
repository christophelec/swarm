package main

import (
	_ "github.com/christophelec/swarm/discovery/file"
	_ "github.com/christophelec/swarm/discovery/kv"
	_ "github.com/christophelec/swarm/discovery/nodes"
	_ "github.com/christophelec/swarm/discovery/token"
	_ "github.com/christophelec/swarm/discovery/serf"

	"github.com/docker/swarm/cli"
)

func main() {
	cli.Run()
}
