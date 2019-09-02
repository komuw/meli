package meli

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "unable to pull image %v", imageName)
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
		if err != nil {
			return err
		}
		defer f.Close()

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
	// reset the buffer since it may contain data from a previous round
	// see issues/118
	for i := range *bufp {
		(*bufp)[i] = 0

	}
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
	// TODO: I dont like the way we are handling paths here.
	// look at dirWithComposeFile in container.go
	dockerFile := dc.ComposeService.Build.Dockerfile
	if dockerFile == "" {
		dockerFile = "Dockerfile"
	}
	dirWithComposeFile := filepath.Dir(dc.DockerComposeFile)
	dirWithComposeFileAbs, err := filepath.Abs(dirWithComposeFile)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get absolute path of %v", dirWithComposeFile)
	}
	userContext := filepath.Dir(dc.ComposeService.Build.Context + "/")
	userContextAbs := filepath.Join(dirWithComposeFileAbs, userContext)
	if filepath.IsAbs(userContext) {
		// For user Contexts that are absolute paths,
		// do NOT join them with anything. They should be used as is.
		userContextAbs = userContext
	}
	if userContextAbs == "/" {
		// ie: dc.ComposeService.Build.Context=="" because user didn't provide any
		userContextAbs = dirWithComposeFile
	}
	dockerFilePath, err := filepath.Abs(
		filepath.Join(userContextAbs, dockerFile))
	if err != nil {
		return "", errors.Wrapf(err, "unable to get path to Dockerfile %v", dockerFile)
	}

	dockerFileReader, err := os.Open(dockerFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to open Dockerfile %v", dockerFilePath)
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read dockerfile %v", dockerFile)
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

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	/*
		Context is either a path to a directory containing a Dockerfile, or a url to a git repository.
		When the value supplied is a relative path, it is interpreted as relative to the location of the Compose file.
		This directory is also the build context that is sent to the Docker daemon.
		- https://docs.docker.com/compose/compose-file/#context
	*/
	UserProvidedContextPath := filepath.Dir(userContextAbs + "/")
	err = filepath.Walk(UserProvidedContextPath, walkFnClosure(UserProvidedContextPath, tw, buf))
	if err != nil {
		return "", errors.Wrapf(err, "unable to walk user provided context path %v", UserProvidedContextPath)
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			//PullParent:     true,
			Squash:         true, //currently only supported in experimenta mode
			Tags:           []string{imageName},
			Remove:         true, //remove intermediary containers after build
			NoCache:        dc.Rebuild,
			SuppressOutput: false,
			Dockerfile:     dockerFile,
			Context:        dockerFileTarReader,
			AuthConfigs:    AuthConfigs})
	if err != nil {
		return "", errors.Wrapf(err, "unable to build docker image %v for service %v", imageName, dc.ServiceName)
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
