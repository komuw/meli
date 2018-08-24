module github.com/komuw/meli

require (
	github.com/Microsoft/go-winio v0.4.10 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.6.2+incompatible // indirect
	github.com/docker/docker v0.0.0-20170601211448-f5ec1e2936dc
	github.com/docker/docker-credential-helpers v0.6.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.3.3 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stevvooe/resumable v0.0.0-20170302213456-2aaf90b2ceea // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 // indirect
	golang.org/x/sys v0.0.0-20180821044426-4ea2f632f6e9 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.1
)

// github.com/docker/engine v18.06.1-ce
replace github.com/docker/docker => github.com/docker/engine v0.0.0-20180816081446-320063a2ad06

// github.com/docker/distribution master
// a proper tagged release is expected in early fall(September 2018)
// see; https://github.com/docker/distribution/issues/2693
replace github.com/docker/distribution => github.com/docker/distribution v2.6.0-rc.1.0.20180820212402-02bf4a2887a4+incompatible
