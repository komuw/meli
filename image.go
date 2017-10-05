package main

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func PullDockerImage(ctx context.Context, imageName string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println(err, "unable to intialize docker client")
	}
	defer cli.Close()
	imagePullResp, err := cli.ImagePull(
		ctx,
		imageName,
		types.ImagePullOptions{})
	if err != nil {
		log.Println(err, "unable to pull image")
	}
	defer imagePullResp.Close()
	_, err = io.Copy(os.Stdout, imagePullResp)
	if err != nil {
		log.Println(err, "unable to write to stdout")
	}

}

func BuildDockerImage(ctx context.Context, dockerFile string) string {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println(err, "unable to intialize docker client")
	}
	defer cli.Close()
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	if dockerFile == "" {
		dockerFile = "Dockerfile"
	}
	dockerFileReader, err := os.Open(dockerFile)
	if err != nil {
		log.Println(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Println(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Println(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Println(err, " :unable to write tar body")
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())
	imageName := "meli_" + strings.ToLower(dockerFile)
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			//PullParent:     true,
			//Squash:     true, currently only supported in experimenta mode
			Tags:           []string{imageName},
			Remove:         true, //remove intermediary containers after build
			ForceRemove:    true,
			SuppressOutput: false,
			Dockerfile:     dockerFile,
			Context:        dockerFileTarReader})
	if err != nil {
		log.Println(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Println(err, " :unable to read image build response")
	}

	return imageName

}
