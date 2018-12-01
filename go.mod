module github.com/komuw/meli

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.11 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/docker/distribution v2.6.2+incompatible // indirect
	github.com/docker/docker v0.0.0-20170601211448-f5ec1e2936dc
	github.com/docker/docker-credential-helpers v0.6.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.3.3 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/sirupsen/logrus v1.2.0 // indirect
	golang.org/x/net v0.0.0-20181201002055-351d144fa1fc // indirect
	golang.org/x/sys v0.0.0-20181128092732-4ed8d59d0b35 // indirect
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/grpc v1.16.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2
	gotest.tools v2.2.0+incompatible // indirect
)

// github.com/docker/engine v18.09.0
replace github.com/docker/docker => github.com/docker/engine v0.0.0-20181106193140-f5749085e9cb

// github.com/docker/distribution v2.7.0-rc.0
// a proper tagged release is expected in early fall(September 2018)
// see; https://github.com/docker/distribution/issues/2693
replace github.com/docker/distribution => github.com/docker/distribution v2.7.0-rc.0+incompatible
