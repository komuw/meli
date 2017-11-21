package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os/user"
	"runtime"
	"strings"
	"sync"

	"github.com/docker/docker-credential-helpers/client"
)

var AuthInfo sync.Map

func useCredStore(server string) (string, string) {
	// this program is usually installed by docker(i think)
	prog := "docker-credential-secretservice"
	goos := runtime.GOOS
	// TODO: handle other Oses or just fail with an error if we encounter an OS that we do not know.
	if goos == "windows" {
		prog = "docker-credential-wincred"
	} else if goos == "darwin" {
		prog = "docker-credential-osxkeychain"
	}

	programfunc := client.NewShellProgramFunc(prog)
	cred, err := client.Get(programfunc, server)
	if err != nil {
		return "", ""
	}

	return cred.Username, cred.Secret
}

func GetAuth() {
	usr, err := user.Current()
	if err != nil {
		AuthInfo.Store("quay", map[string]string{"registryURL": "", "username": "", "password": ""})
		AuthInfo.Store("dockerhub", map[string]string{"registryURL": "", "username": "", "password": ""})
		return
	}

	// TODO: the config can be in many places
	// try to find them and use them; https://github.com/docker/docker-py/blob/e9fab1432b974ceaa888b371e382dfcf2f6556e4/docker/auth.py#L269
	dockerAuth, err := ioutil.ReadFile(usr.HomeDir + "/.docker/config.json")
	if err != nil {
		// we'll just try accessing the public access docker hubs/quay
		AuthInfo.Store("quay", map[string]string{"registryURL": "", "username": "", "password": ""})
		AuthInfo.Store("dockerhub", map[string]string{"registryURL": "", "username": "", "password": ""})
		return
	}

	type AuthData struct {
		Auths      map[string]map[string]string `json:"auths,omitempty"`
		CredsStore string                       `json:"credsStore,omitempty"`
	}
	data := &AuthData{}
	err = json.Unmarshal([]byte(dockerAuth), data)
	if err != nil {
		AuthInfo.Store("quay", map[string]string{"registryURL": "", "username": "", "password": ""})
		AuthInfo.Store("dockerhub", map[string]string{"registryURL": "", "username": "", "password": ""})
		return
	}

	// TODO: we are only checking for dockerHub and quay.io
	// registries, we should probably be exhaustive in future.
	dockerEncodedAuth := data.Auths["https://index.docker.io/v1/"]["auth"]
	dockerRegistryURL := "https://index.docker.io/v1/"
	quayEncodedAuth := data.Auths["quay.io"]["auth"]
	quayRegistryURL := "quay.io"

	if dockerEncodedAuth == "" {
		AuthInfo.Store("dockerhub", map[string]string{"registryURL": "", "username": "", "password": ""})
	}
	if quayEncodedAuth == "" {
		AuthInfo.Store("quay", map[string]string{"registryURL": "", "username": "", "password": ""})
	}
	dockerAuth, err = base64.StdEncoding.DecodeString(dockerEncodedAuth)
	if err != nil {
		AuthInfo.Store("dockerhub", map[string]string{"registryURL": "", "username": "", "password": ""})
	}
	quayAuth, err := base64.StdEncoding.DecodeString(quayEncodedAuth)
	if err != nil {
		AuthInfo.Store("quay", map[string]string{"registryURL": "", "username": "", "password": ""})
	}

	dockerUsername, dockerPassword, quayUsername, quayPassword := "PLACEHOLDER", "PLACEHOLDER", "PLACEHOLDER", "PLACEHOLDER"
	if data.CredsStore != "" {
		dockerUsername, dockerPassword = useCredStore(dockerRegistryURL)
		quayUsername, quayPassword = useCredStore(quayRegistryURL)

	} else {
		dockerUserPass := FormatRegistryAuth(string(dockerAuth))
		quayUserPass := FormatRegistryAuth(string(quayAuth))

		if len(dockerUserPass) < 2 {
			dockerUsername, dockerPassword = "", ""
		} else {
			dockerUsername = dockerUserPass[0]
			dockerPassword = dockerUserPass[1]
		}
		if len(quayUserPass) < 2 {
			quayUsername, quayPassword = "", ""
		} else {
			quayUsername = quayUserPass[0]
			quayPassword = quayUserPass[1]
		}
	}

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

	AuthInfo.Store("dockerhub", map[string]string{"registryURL": dockerRegistryURL, "username": dockerUsername, "password": dockerPassword, "RegistryAuth": dockerRegistryAuth})
	AuthInfo.Store("quay", map[string]string{"registryURL": quayRegistryURL, "username": quayUsername, "password": quayPassword, "RegistryAuth": quayRegistryAuth})

}
