package api

import (
	"fmt"
	"path/filepath"
	"strings"
)

func FormatContainerName(containerName string) string {
	// container names are supposed to be unique
	// we are using the docker-compose service as the container name
	f := func(c rune) bool {
		if c == 58 {
			// 58 is the ':' character
			return true
		}
		return false
	}
	return strings.FieldsFunc(containerName, f)[0]
}

func FormatLabels(label string) []string {
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

func FormatPorts(port string) []string {
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

func FormatServiceVolumes(volume string) []string {
	f := func(c rune) bool {
		if c == 58 {
			// 58 is the ':' character
			return true
		}
		return false
	}
	// TODO: we should trim any whitespace before returning.
	// this will prevent labels like type= web
	x := strings.FieldsFunc(volume, f)
	if strings.Contains(x[0], "./") {
		fmt.Println(x[0])
		path, err := filepath.Abs(x[0])
		fmt.Println("x0, path, err", x[0], path, err)
	}

	return x
}

func FormatRegistryAuth(auth string) []string {
	f := func(c rune) bool {
		if c == 58 {
			// 58 is the ':' character
			return true
		}
		return false
	}
	// TODO: we should trim any whitespace before returning.
	// this will prevent labels like type= web
	return strings.FieldsFunc(auth, f)
}

func FormatComposePath(path string) []string {
	f := func(c rune) bool {
		// TODO; check if this is cross platform
		if c == 47 {
			// 47 is the '/' character
			return true
		}
		return false
	}
	// TODO: we should trim any whitespace before returning.
	return strings.FieldsFunc(path, f)
}

type popagateError struct {
	originalErr error
	newErr      error
}

func (p *popagateError) Error() string {
	return fmt.Sprintf("originalErr:: %s \nThisErr:: %s", p.originalErr, p.newErr)
}
