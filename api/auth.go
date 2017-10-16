package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os/user"
	"strings"
)

var AuthInfo = make(map[string]map[string]string)

func GetAuth() {
	usr, err := user.Current()
	if err != nil {
		// unable to find user
		AuthInfo["quay"] = map[string]string{"registryURL": "", "username": "", "password": ""}
		AuthInfo["dockerhub"] = map[string]string{"registryURL": "", "username": "", "password": ""}
		return
	}

	// TODO: the config can be in many places
	// try to find them and use them; https://github.com/docker/docker-py/blob/e9fab1432b974ceaa888b371e382dfcf2f6556e4/docker/auth.py#L269
	dockerAuth, err := ioutil.ReadFile(usr.HomeDir + "/.docker/config.json")
	if err != nil {
		// we'll just try accessing the public access docker hubs/quay
		AuthInfo["quay"] = map[string]string{"registryURL": "", "username": "", "password": ""}
		AuthInfo["dockerhub"] = map[string]string{"registryURL": "", "username": "", "password": ""}
		return
	}

	type AuthData struct {
		Auths map[string]map[string]string `json:"auths,omitempty"`
	}
	data := &AuthData{}
	err = json.Unmarshal([]byte(dockerAuth), data)
	if err != nil {
		AuthInfo["quay"] = map[string]string{"registryURL": "", "username": "", "password": ""}
		AuthInfo["dockerhub"] = map[string]string{"registryURL": "", "username": "", "password": ""}
		return
	}

	// TODO: we are only checking for dockerHub and quay.io
	// registries, we should probably be exhaustive in future.
	dockerEncodedAuth := data.Auths["https://index.docker.io/v1/"]["auth"]
	dockerRegistryURL := "https://index.docker.io/v1/"
	quayEncodedAuth := data.Auths["quay.io"]["auth"]
	quayRegistryURL := "quay.io"

	if dockerEncodedAuth == "" {
		AuthInfo["dockerhub"] = map[string]string{"registryURL": "", "username": "", "password": ""}
	}
	if quayEncodedAuth == "" {
		AuthInfo["quay"] = map[string]string{"registryURL": "", "username": "", "password": ""}
	}
	dockerAuth, err = base64.StdEncoding.DecodeString(dockerEncodedAuth)
	if err != nil {
		AuthInfo["dockerhub"] = map[string]string{"registryURL": "", "username": "", "password": ""}
	}
	quayAuth, err := base64.StdEncoding.DecodeString(quayEncodedAuth)
	if err != nil {
		AuthInfo["quay"] = map[string]string{"registryURL": "", "username": "", "password": ""}
	}

	dockerUserPass := FormatRegistryAuth(string(dockerAuth))
	dockerUsername := dockerUserPass[0]
	dockerPassword := dockerUserPass[1]

	quayUserPass := FormatRegistryAuth(string(quayAuth))
	quayUsername := quayUserPass[0]
	quayPassword := quayUserPass[1]

	dockerStringRegistryAuth := `{"username": "DOCKERUSERNAME", "password": "DOCKERPASSWORD", "email": null, "serveraddress": "DOCKERREGISTRYURL"}`
	dockerStringRegistryAuth = strings.Replace(dockerStringRegistryAuth, "DOCKERUSERNAME", dockerUsername, 1)
	dockerStringRegistryAuth = strings.Replace(dockerStringRegistryAuth, "DOCKERPASSWORD", dockerPassword, 1)
	dockerStringRegistryAuth = strings.Replace(dockerStringRegistryAuth, "DOCKERREGISTRYURL", dockerRegistryURL, 1)
	dockerRegistryAuth := base64.URLEncoding.EncodeToString([]byte(dockerStringRegistryAuth))

	quayStringRegistryAuth := `{"username": "quayUSERNAME", "password": "quayPASSWORD", "email": null, "serveraddress": "quayREGISTRYURL"}`
	quayStringRegistryAuth = strings.Replace(quayStringRegistryAuth, "quayUSERNAME", quayUsername, 1)
	quayStringRegistryAuth = strings.Replace(quayStringRegistryAuth, "quayPASSWORD", quayPassword, 1)
	quayStringRegistryAuth = strings.Replace(quayStringRegistryAuth, "quayREGISTRYURL", quayRegistryURL, 1)
	quayRegistryAuth := base64.URLEncoding.EncodeToString([]byte(quayStringRegistryAuth))

	AuthInfo["dockerhub"] = map[string]string{"registryURL": dockerRegistryURL, "username": dockerUsername, "password": dockerPassword, "RegistryAuth": dockerRegistryAuth}
	AuthInfo["quay"] = map[string]string{"registryURL": quayRegistryURL, "username": quayUsername, "password": quayPassword, "RegistryAuth": quayRegistryAuth}

}
