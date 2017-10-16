package api

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func PullDockerImage(ctx context.Context, imageName string, cli MeliAPiClient) error {
	result, _ := AuthInfo.Load("dockerhub")
	if strings.Contains(imageName, "quay") {
		result, _ = AuthInfo.Load("quay")
	}
	GetRegistryAuth := result.(map[string]string)["RegistryAuth"]

	imagePullResp, err := cli.ImagePull(
		ctx,
		imageName,
		types.ImagePullOptions{RegistryAuth: GetRegistryAuth})
	if err != nil {
		log.Println(err, " :unable to pull image")
	}
	defer imagePullResp.Close()

	// supplying your own buffer is perfomant than letting the system do it for you
	buff := make([]byte, 2048)
	io.CopyBuffer(os.Stdout, imagePullResp, buff)

	return nil
}

func BuildDockerImage(ctx context.Context, dockerFile string, cli MeliAPiClient) (string, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFileReader, err := os.Open(dockerFile)
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to open Dockerfile %s", dockerFile)}
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to read dockerfile")}
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to write tar header")}
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to write tar body")}
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())
	imageName := "meli_" + strings.ToLower(dockerFile)

	splitDockerfile := strings.Split(string(readDockerFile), " ")
	splitImageName := strings.Split(splitDockerfile[1], "\n")
	imgFromDockerfile := splitImageName[0]

	result, _ := AuthInfo.Load("dockerhub")
	if strings.Contains(imgFromDockerfile, "quay") {
		result, _ = AuthInfo.Load("quay")
	}
	authInfo := result.(map[string]string)
	registryURL := authInfo["registryURL"]
	username := authInfo["username"]
	password := authInfo["password"]

	AuthConfigs := make(map[string]types.AuthConfig)
	AuthConfigs[registryURL] = types.AuthConfig{Username: username, Password: password}

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
			Context:        dockerFileTarReader,
			AuthConfigs:    AuthConfigs})
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to build docker image")}
	}
	defer imageBuildResponse.Body.Close()

	buff := make([]byte, 2048)
	io.CopyBuffer(os.Stdout, imageBuildResponse.Body, buff)

	return imageName, nil
}
