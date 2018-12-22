module github.com/komuw/meli

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.10 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.6.2+incompatible // indirect
	github.com/docker/docker v0.0.0-20170601211448-f5ec1e2936dc
	github.com/docker/docker-credential-helpers v0.6.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.3.3 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/onsi/ginkgo v1.6.0 // indirect
	github.com/onsi/gomega v1.4.1 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.0.6 // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac // indirect
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 // indirect
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f // indirect
	golang.org/x/sys v0.0.0-20180821044426-4ea2f632f6e9 // indirect
	golang.org/x/text v0.3.0 // indirect
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 // indirect
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 // indirect
	google.golang.org/grpc v1.14.0 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.1
	gotest.tools v2.1.0+incompatible // indirect
)

// look at: https://github.com/golang/go/issues/29376#issuecomment-449416502

// github.com/docker/engine v18.06.1-ce
replace github.com/docker/docker => github.com/docker/engine v0.0.0-20180816081446-320063a2ad06

// github.com/docker/distribution master
// a proper tagged release is expected in early fall(September 2018)
// see; https://github.com/docker/distribution/issues/2693
replace github.com/docker/distribution => github.com/docker/distribution v2.6.0-rc.1.0.20180820212402-02bf4a2887a4+incompatible
