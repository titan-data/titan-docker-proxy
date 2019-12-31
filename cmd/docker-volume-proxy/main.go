package main

import (
	"context"
    "fmt"
    "flag"
    "os"

    titan "github.com/titan-data/titan-client-go"
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

    config := titan.NewConfiguration()
    config.Host = fmt.Sprintf("%s:%d", *host, *port)
    apis := titan.NewAPIClient(config)
    volumeApi := apis.VolumesApi

    ctx := context.Background()
    volumes, _, _ := volumeApi.ListVolumes(ctx, "mongo")

    for _, volume := range volumes {
        println(volume.Name)
    }
}
