# meli            

[![CircleCI](https://circleci.com/gh/komuW/meli.svg?style=svg)](https://circleci.com/gh/komuW/meli)        


Meli is supposed to be a faster alternative to docker-compose. Faster in the sense that, Meli will try to pull as many services(docker containers) 
as it can in parallel.

Meli is a Swahili word meaning ship; so think of Meli as a ship carrying your docker containers safely across the treacherous container seas.

It's currently work in progress, API will remain unstable for sometime.

I only intend to support docker-compose [version 3+](https://docs.docker.com/compose/compose-file/compose-versioning/)          

Meli is NOT intended to replicate every feature of docker-compose, it is primarily intended to enable you to pull, build and run the services in your docker-compose file as fast as possible.          
If you want to exec in to a running container, use docker; if you want to run an adhoc command in a container, use docker; if you want..... you get the drift.


# Installing          
Very early test releases are available from the [releases page](https://github.com/komuW/meli/releases)          

# Usage  
`meli --help`         
```bash
Usage of meli:
  -up
    	Builds, re/creates, starts, and attaches to containers for a service.
        -d option runs containers in the background
  -f string
    	path to docker-compose.yml file. (default "docker-compose.yml")
  -v	Show version information.
  -version
    	Show version information.
```

`cat docker-compose.yml`                 
```bash
version: '3'
services:
  redis:
    image: 'redis:3.0-alpine'
    environment:
      - RACK_ENV=development
      - type=database
    ports:
      - "6300:6379"
      - "6400:22"
```           

`meli -up`          
```bash 
2017/10/07 14:30:11 {"status":"Pulling from library/redis","id":"3.0-alpine"}
2017/10/07 14:30:11 {"status":"Digest: sha256:350469b395eac82395f9e59d7b7b90f7d23fe0838965e56400739dec3afa60de"}
2017/10/07 14:30:11 {"status":"Status: Image is up to date for redis:3.0-alpine"}
...
2017/10/07 14:30:12 u2017-10-07T11:30:12.720619075Z  1:M 07 Oct 11:30:12.720 * The server is now ready to accept connections on port 6379
```

# Usage as a library
You really should be using the official [docker Go sdk](https://godoc.org/github.com/moby/moby/client)         
However, if you feel inclined to use meli;
```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/komuw/meli/api"
)

func main() {
	dc := &api.DockerContainer{
		ComposeService: api.ComposeService{Image: "busybox"},
		LogMedium:      os.Stdout,
		FollowLogs:     true}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to intialize docker client")
	}
	defer cli.Close()

	api.GetAuth() // read dockerhub info
	err = api.PullDockerImage(ctx, cli, dc)
	log.Println(err)

}

```


# Benchmarks
Aaah, everyones' favorite vanity metric yet no one seems to know how to conduct one in a consistent and scientific manner.          
Take any results you see here with a large spoon of salt; They are unscientific and reek of all that is wrong with most developer benchmarks.             

Having made that disclaimer,                 

test machine:             
[This circleCI machine](https://github.com/komuW/meli/blob/master/.circleci/config.yml#L9)


docker-compose version:         
`docker-compose --version`
```bash
docker-compose version 1.16.1, build 6d1ac219
```

Meli version:   
[version 0.0.8](https://github.com/komuW/meli/releases/tag/v0.0.8)
           

Benchmark test:           
[this docker-compose file](https://github.com/komuW/meli/blob/master/testdata/docker-compose.yml)

Benchmark script:               
for docker-compose:      
`docker ps -aq | xargs docker rm -f; docker system prune -af; /usr/bin/time -apv docker-compose up -d`        
for meli:                
`docker ps -aq | xargs docker rm -f; docker system prune -af; /usr/bin/time -apv meli -up -d`            

Benchmark results(average):                       

| tool           | Elapsed wall clock time(seconds) |
| :---           |          ---:                    |
| docker-compose |  9.389 seconds                  |
| meli           |  3.658  seconds                  |

Thus, meli appears to be 2.56 times faster than docker-compose(by wall clock time).           
You can [checkout the current benchmark results from the circleCI](https://circleci.com/gh/komuW/meli/)              
However, I'm not making a tool to take docker-compose to the races.                   

# Build                   
`git clone git@github.com:komuW/meli.git`           
`go build -o meli main.go`           
`./meli -up -f /path/to/docker-compose-file.yml`                   


# TODO
- add better documentation(godoc)
- stabilise API(maybe)

