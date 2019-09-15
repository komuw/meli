module github.com/komuw/meli

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect

	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/grpc v1.23.0 // indirect
	// TODO: tag with a proper version after https://github.com/go-yaml/yaml/issues/487 is fixed
	gopkg.in/yaml.v3 v3.0.0-20190709130402-674ba3eaed22
	gotest.tools v2.2.0+incompatible // indirect
)

// look at: https://github.com/golang/go/issues/29376#issuecomment-449416502

// github.com/docker/engine v19.03.1
replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190725163905-fa8dd90ceb7b
