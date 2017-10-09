package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/user"
	"strings"
)

func GetRegistryAuth(imageName string) (string, error) {
	registryURL, username, password, err := GetAuth(imageName)
	if err != nil {
		return "", &popagateError{originalErr: err}
	}

	stringRegistryAuth := `{"username": "DOCKERUSERNAME", "password": "DOCKERPASSWORD", "email": null, "serveraddress": "DOCKERREGISTRYURL"}`

	stringRegistryAuth = strings.Replace(stringRegistryAuth, "DOCKERUSERNAME", username, 1)
	stringRegistryAuth = strings.Replace(stringRegistryAuth, "DOCKERPASSWORD", password, 1)
	stringRegistryAuth = strings.Replace(stringRegistryAuth, "DOCKERREGISTRYURL", registryURL, 1)
	RegistryAuth := base64.URLEncoding.EncodeToString([]byte(stringRegistryAuth))

	return RegistryAuth, nil
}

func GetAuth(imageName string) (string, string, string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", "", "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to find current user")}
	}
	// TODO: the config can be in many places
	// try to find them and use them; https://github.com/docker/docker-py/blob/e9fab1432b974ceaa888b371e382dfcf2f6556e4/docker/auth.py#L269
	dockerAuth, err := ioutil.ReadFile(usr.HomeDir + "/.docker/config.json")
	if err != nil {
		return "", "", "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to read docker auth file, ~/.docker/config.json")}
	}

	type AuthData struct {
		Auths map[string]map[string]string `json:"auths,omitempty"`
	}
	data := &AuthData{}
	err = json.Unmarshal([]byte(dockerAuth), data)
	if err != nil {
		return "", "", "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to unmarshal auth info")}
	}

	encodedAuth := "placeholder"
	registryURL := "placeholder"
	// TODO: we are only checking for dockerHub and quay.io
	// registries, we should probably be exhaustive in future.
	if strings.Contains(imageName, "quay") {
		// quay
		encodedAuth = data.Auths["quay.io"]["auth"]
		registryURL = "quay.io"
	} else {
		encodedAuth = data.Auths["https://index.docker.io/v1/"]["auth"]
		registryURL = "https://index.docker.io/v1/"
	}

	if encodedAuth == "" {
		return "", "", "", &popagateError{
			newErr: errors.New(" :unable to find any auth info in ~/.docker/config.json")}
	}

	yourAuth, err := base64.StdEncoding.DecodeString(encodedAuth)
	if err != nil {
		return "", "", "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to base64 decode auth info")}
	}
	userPass := fomatRegistryAuth(string(yourAuth))
	username := userPass[0]
	password := userPass[1]

	return registryURL, username, password, nil

}
