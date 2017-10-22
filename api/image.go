package api

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
)

func PullDockerImage(ctx context.Context, cli MeliAPiClient, dc *DockerContainer) error {
	imageName := dc.ComposeService.Image
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
		return &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to pull image %s", imageName)}
	}
	defer imagePullResp.Close()

	// supplying your own buffer is perfomant than letting the system do it for you
	buff := make([]byte, 2048)
	io.CopyBuffer(os.Stdout, imagePullResp, buff)

	return nil
}

func walkFnClosure(src string, tw *tar.Writer, buf *bytes.Buffer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// todo: maybe we should return nil
			return err
		}

		tarHeader, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		// update the name to correctly reflect the desired destination when untaring
		// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
		tarHeader.Name = strings.TrimPrefix(strings.Replace(path, src, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(tarHeader)
		if err != nil {
			return err
		}
		// return on directories since there will be no content to tar
		if info.Mode().IsDir() {
			return nil
		}

		// open files for taring
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}
		readFile, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		_, err = tw.Write(readFile)
		if err != nil {
			return err
		}
		return nil
	}
}

//func BuildDockerImage(ctx context.Context, k, dockerFile string, cli MeliAPiClient) (string, error) {

func BuildDockerImage(ctx context.Context, cli MeliAPiClient, dc *DockerContainer) (string, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := dc.ComposeService.Build.Dockerfile
	if dockerFile == "" {
		dockerFile = "Dockerfile"
	}

	// TODO: we should probably use the filepath stdlib module
	// so that atleast it can guarantee us os agnotic'ness
	formattedDockerComposePath := FormatComposePath(dc.DockerComposeFile)
	if len(formattedDockerComposePath) == 0 {
		// very unlikely to hit this situation, but
		return "", fmt.Errorf(" :docker-compose file is empty %s", dc.DockerComposeFile)
	}
	pathToDockerFile := formattedDockerComposePath[0]
	if pathToDockerFile != "docker-compose.yml" {
		dockerFile = pathToDockerFile + "/" + dockerFile
	}

	dockerFilePath, err := filepath.Abs(dockerFile)
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to get path to Dockerfile %s", dockerFile)}
	}
	dockerContextPath := filepath.Dir(dockerFilePath)
	dockerFileName := filepath.Base(dockerFile)

	dockerFileReader, err := os.Open(dockerFilePath)
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

	imageName := "meli_" + strings.ToLower(dc.ServiceName)

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

	// TODO: we need to read the context passed in the docker-compose context key for a service
	// rather than assume the context is the dir the Dockerfile is in.
	err = filepath.Walk(dockerContextPath, walkFnClosure(dockerContextPath, tw, buf))
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to walk dockefile context path %s", dockerFile)}
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

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
			Dockerfile:     dockerFileName,
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
