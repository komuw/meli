module github.com/komuw/meli

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.11 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/docker/distribution v2.7.0+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20181221150755-2cb26cfe9cbf
	github.com/docker/docker-credential-helpers v0.6.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.3.3 // indirect
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.0
	github.com/sirupsen/logrus v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20190422183909-d864b10871cd // indirect
	golang.org/x/net v0.0.0-20190424112056-4829fb13d2c6 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	golang.org/x/sys v0.0.0-20190424160641-4347357a82bc // indirect
	golang.org/x/text v0.3.1 // indirect
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/genproto v0.0.0-20181221175505-bd9b4fb69e2f // indirect
	google.golang.org/grpc v1.17.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2
	gotest.tools v2.2.0+incompatible // indirect
)

// look at: https://github.com/golang/go/issues/29376#issuecomment-449416502

// github.com/docker/engine  v18.09.1-rc1
replace github.com/docker/docker => github.com/docker/engine v0.0.0-20181203220643-82a4418f57d5

go 1.13
