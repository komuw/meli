package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	imagename := "busybox"
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	imagePullResp, err := cli.ImagePull(ctx, imagename, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, imagePullResp)
	if err != nil {
		log.Println(err)
	}

}
