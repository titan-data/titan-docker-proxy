/*
 * Copyright The Titan Project Contributors.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/titan-data/titan-docker-proxy/internal/proxy"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: docker-volume-proxy [--host host] [--port port] socket\n")
		flag.PrintDefaults()
	}

	host := flag.String("host", "localhost", "host to connect to")
	port := flag.Int("port", 5001, "port to connect to")

	flag.Parse()

	if flag.NArg() != 1 {
		panic("missing required socket path")
	}
	path := flag.Arg(0)

	fmt.Printf("Proxying requests from %s to %s:%d\n", path, *host, *port)

	titan := proxy.Proxy(*host, *port)
	volumes := titan.ListVolumes()

	for _, vol := range volumes.Volumes {
		fmt.Println(vol.Name)
	}
}
