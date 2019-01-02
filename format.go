package meli

import (
	"os"
	"path/filepath"
	"strings"
)

func formatContainerName(containerName, curentDir string) string {
	// container names are supposed to be unique
	// we are using the docker-compose service name as well as current dir as the container name
	f := func(c rune) bool {
		return c == 58 // 58 is the ':' character
	}
	formattedContainerName := strings.FieldsFunc(containerName, f)[0]
	contName := "meli_" + formattedContainerName + curentDir

	return contName
}

func formatLabels(label string) []string {
	f := func(c rune) bool {
		if c == 58 {
			// 58 is the ':' character
			return true
		} else if c == 61 {
			//61 is '=' char
			return true
		}
		return false
	}
	// TODO: we should trim any whitespace before returning.
	// this will prevent labels like type= web
	return strings.FieldsFunc(label, f)
}

func formatPorts(port string) []string {
	f := func(c rune) bool {
		if c == 58 {
			// 58 is the ':' character
			return true
		} else if c == 61 {
			//61 is '=' char
			return true
		}
		return false
	}
	// TODO: we should trim any whitespace before returning.
	// this will prevent labels like type= web
	return strings.FieldsFunc(port, f)
}

func formatServiceVolumes(volume, dockerComposeFile string) []string {
	f := func(c rune) bool {
		return c == 58 // 58 is the ':' character
	}
	volume = os.ExpandEnv(volume)
	// TODO: we should trim any whitespace before returning.
	// this will prevent labels like type= web
	hostAndContainerPath := strings.FieldsFunc(volume, f)
	dockerComposeFileDir := filepath.Dir(dockerComposeFile)

	if strings.Contains(hostAndContainerPath[0], "./") {
		dockerComposeFilePath, _ := filepath.Abs(hostAndContainerPath[0])
		hostPath := filepath.Join(dockerComposeFilePath, dockerComposeFileDir)
		hostAndContainerPath[0] = hostPath
	} else if strings.HasPrefix(hostAndContainerPath[0], ".") {
		dockerComposeFileDirAbs, _ := filepath.Abs(dockerComposeFileDir)
		hostPath := filepath.Join(dockerComposeFileDirAbs, hostAndContainerPath[0])
		hostAndContainerPath[0] = hostPath
	}

	return hostAndContainerPath
}

func formatRegistryAuth(auth string) []string {
	f := func(c rune) bool {
		return c == 58 // 58 is the ':' character
	}
	// TODO: we should trim any whitespace before returning.
	// this will prevent labels like type= web
	return strings.FieldsFunc(auth, f)
}

func formatComposePath(path string) []string {
	f := func(c rune) bool {
		// TODO; check if this is cross platform
		return c == 47 // 47 is the '/' character
	}
	// TODO: we should trim any whitespace before returning.
	return strings.FieldsFunc(path, f)
}
