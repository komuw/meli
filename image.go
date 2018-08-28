package meli

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
)

// PullDockerImage pulls a docker from a registry via docker daemon
func PullDockerImage(ctx context.Context, cli APIclient, dc *DockerContainer) error {
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

	var imgProg imageProgress
	scanner := bufio.NewScanner(imagePullResp)
	for scanner.Scan() {
		_ = json.Unmarshal(scanner.Bytes(), &imgProg)
		fmt.Fprintln(dc.LogMedium, dc.ServiceName, "::", imgProg.Status, imgProg.Progress)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(" :unable to log output for image", imageName, err)
	}

	imagePullResp.Close()
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
		if src == "." {
			// see: issues/74
			tarHeader.Name = strings.TrimPrefix(path, string(filepath.Separator))
		}

		err = tw.WriteHeader(tarHeader)
		if err != nil {
			return err
		}
		// return on directories since there will be no content to tar
		if info.Mode().IsDir() {
			return nil
		}
		// return on non-regular files since there will be no content to tar
		if !info.Mode().IsRegular() {
			// non regular files are like symlinks etc; https://golang.org/src/os/types.go?h=ModeSymlink#L49
			return nil
		}

		// open files for taring
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}

		tr := io.TeeReader(f, tw)
		_, err = poolReadFrom(tr)
		if err != nil {
			return err
		}

		return nil
	}
}

// this is taken from io.util
var blackHolePool = sync.Pool{
	New: func() interface{} {
		// TODO: change this size accordingly
		// we could find the size of the file we want to tar
		// then pass that in as the size. That way we will
		// always create a right sized slice and not have to incure cost of slice regrowth(if any)
		b := make([]byte, 512)
		return &b
	},
}

// this is taken from io.util
func poolReadFrom(r io.Reader) (n int64, err error) {
	bufp := blackHolePool.Get().(*[]byte)
	readSize := 0
	for {
		readSize, err = r.Read(*bufp)
		n += int64(readSize)
		if err != nil {
			blackHolePool.Put(bufp)
			if err == io.EOF {
				return n, nil
			}
			return
		}
	}
}

// BuildDockerImage builds a docker image via docker daemon
func BuildDockerImage(ctx context.Context, cli APIclient, dc *DockerContainer) (string, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// TODO: I dont like the way we are handling paths here.
	// look at dirWithComposeFile in container.go
	dockerFile := dc.ComposeService.Build.Dockerfile
	if dockerFile == "" {
		dockerFile = "Dockerfile"
	}
	formattedDockerComposePath := formatComposePath(dc.DockerComposeFile)
	if len(formattedDockerComposePath) == 0 {
		// very unlikely to hit this situation, but
		return "", fmt.Errorf(" :docker-compose file is empty %s", dc.DockerComposeFile)
	}
	pathToDockerFile := formattedDockerComposePath[0]
	if pathToDockerFile != "docker-compose.yml" {
		dockerFile = filepath.Join(pathToDockerFile, dockerFile)
	}

	dockerFilePath, err := filepath.Abs(dockerFile)
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to get path to Dockerfile %s", dockerFile)}
	}

	dockerFileContextPath := filepath.Dir(dockerFile)
	/*
		Context is either a path to a directory containing a Dockerfile, or a url to a git repository.
		When the value supplied is a relative path, it is interpreted as relative to the location of the Compose file.
		This directory is also the build context that is sent to the Docker daemon.
		- https://docs.docker.com/compose/compose-file/#context

		So it looks like, we only need to send one of
		UserContext or dockerFileContextPath to docker server and not two.
	*/
	UserProvidedContextPath := filepath.Dir(dc.ComposeService.Build.Context + "/")
	if dc.ComposeService.Build.Context == "." {
		// context will be the directory containing the compose file
		UserProvidedContextPath = filepath.Dir(dc.DockerComposeFile)
	} else if dc.ComposeService.Build.Context == "" {
		// context will be the directory containing the compose file
		UserProvidedContextPath = filepath.Dir(dc.DockerComposeFile)
	}

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

	// TODO: stop calling filepath.Walk twice; once for dockerfile context and then for uer context
	// Both the user provided context and that of the DockerFile needs to be in the tar file.
	err = filepath.Walk(dockerFileContextPath, walkFnClosure(dockerFileContextPath, tw, buf))
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to walk dockefile context path %s", dockerFile)}
	}
	err = filepath.Walk(UserProvidedContextPath, walkFnClosure(UserProvidedContextPath, tw, buf))
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to walk user provided context path %s", dockerFile)}
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	dockerFileName := filepath.Base(dockerFile)
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			//PullParent:     true,
			//Squash:     true, currently only supported in experimenta mode
			Tags:           []string{imageName},
			Remove:         true, //remove intermediary containers after build
			NoCache:        dc.Rebuild,
			SuppressOutput: false,
			Dockerfile:     dockerFileName,
			Context:        dockerFileTarReader,
			AuthConfigs:    AuthConfigs})
	if err != nil {
		return "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to build docker image")}
	}

	var imgProg imageProgress
	scanner := bufio.NewScanner(imageBuildResponse.Body)
	for scanner.Scan() {
		_ = json.Unmarshal(scanner.Bytes(), &imgProg)
		fmt.Fprint(
			dc.LogMedium,
			dc.ServiceName,
			"::",
			imgProg.Status,
			imgProg.Progress,
			imgProg.Stream)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(" :unable to log output for image", imageName, err)
	}

	imageBuildResponse.Body.Close()
	return imageName, nil
}
