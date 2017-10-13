package api

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func PullDockerImage(ctx context.Context, imageName string, cli MeliAPiClient) error {
	GetRegistryAuth, err := GetRegistryAuth(imageName)
	if err != nil {
		log.Println(err, " :unable to get registry credentials for image, ", imageName)
		return err
	}

	imagePullResp, err := cli.ImagePull(
		ctx,
		imageName,
		types.ImagePullOptions{RegistryAuth: GetRegistryAuth})
	if err != nil {
		log.Println(err, " :unable to pull image")
	}
	defer imagePullResp.Close()

	scanner := bufio.NewScanner(imagePullResp)
	for scanner.Scan() {
		output := strings.Replace(scanner.Text(), "u003e", ">", -1)
		log.Println(output)
	}
	err = scanner.Err()
	if err != nil {
		log.Println(err, "error in scanning")
	}

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

	registryURL, username, password, err := GetAuth(imgFromDockerfile)
	if err != nil {
		return "", &popagateError{originalErr: err}
	}
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

	scanner := bufio.NewScanner(imageBuildResponse.Body)
	for scanner.Scan() {
		output := strings.Replace(scanner.Text(), "u003e", ">", -1)
		log.Println(output)
	}
	err = scanner.Err()
	if err != nil {
		log.Println(err, "error in scanning")
	}

	return imageName, nil
}
