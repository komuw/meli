# meli            

[![CircleCI](https://circleci.com/gh/komuw/meli.svg?style=svg)](https://circleci.com/gh/komuw/meli)
[![codecov](https://codecov.io/gh/komuw/meli/branch/master/graph/badge.svg)](https://codecov.io/gh/komuw/meli)
[![GoDoc](https://godoc.org/github.com/komuw/meli?status.svg)](https://godoc.org/github.com/komuw/meli)
[![Go Report Card](https://goreportcard.com/badge/github.com/komuw/meli)](https://goreportcard.com/report/github.com/komuw/meli)          


Meli is supposed to be a faster, and drop in, alternative to docker-compose. Faster in the sense that, Meli will try to pull as many services(docker containers) 
as it can in parallel. 

Meli is a Swahili word meaning ship; so think of Meli as a ship carrying your docker containers safely across the treacherous container seas.

It's currently work in progress, API will remain unstable for sometime.

I only intend to support docker-compose [version 3+](https://docs.docker.com/compose/compose-file/compose-versioning/)          

Meli is NOT intended to replicate every feature of docker-compose, it is primarily intended to enable you to pull, build and run the services in your docker-compose file as fast as possible.          
If you want to exec in to a running container, use docker; if you want to run an adhoc command in a container, use docker; if you want..... you get the drift.


# Installing/Upgrading          
Download a binary release for your particular OS from the [releases page](https://github.com/komuW/meli/releases)           
We have binaries for:                
- linux(32bit and 64bit)           
- windows(32bit and 64bit)            
- darwin(32bit and 64bit)                     

Optionally, you can install using curl;       
```bash
curl -sfL https://raw.githubusercontent.com/komuw/meli/master/install.sh | sh
```

# Usage  
`meli --help`         
```bash
Usage of meli:
  -up
    	Builds, re/creates, starts, and attaches to containers for a service.
        -d option runs containers in the background
  -f string
    	path to docker-compose.yml file. (default "docker-compose.yml")
  -build
    	Rebuild services
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
redis :: Pulling from library/redis 
redis :: Pulling fs layer 
redis :: Pulling fs layer 
redis :: Downloading [======================>        ]  3.595kB/8.164kB
redis :: Downloading [==============================>]  8.164kB/8.164kB
redis :: Download complete [========================>]  8.164kB/8.164kB
redis :: The server is now ready to accept connections on port 6379
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
	"github.com/komuw/meli"
)

func main() {
	dc := &meli.DockerContainer{
		ComposeService: meli.ComposeService{Image: "busybox"},
		LogMedium:      os.Stdout,
		FollowLogs:     true}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to intialize docker client")
	}
	defer cli.Close()

	meli.GetAuth() // read dockerhub info
	err = meli.PullDockerImage(ctx, cli, dc)
	log.Println(err)

}
```


# Benchmarks
Aaah, everyones' favorite vanity metric yet no one seems to know how to conduct one in a consistent and scientific manner.          
Take any results you see here with a large spoon of salt; They are unscientific and reek of all that is wrong with most developer benchmarks.             

Having made that disclaimer,                 

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
| docker-compose |  10.411 seconds                  |
| meli           |  3.945  seconds                  |

Thus, meli appears to be 2.6 times faster than docker-compose(by wall clock time).           
You can [checkout the current benchmark results from circleCI](https://circleci.com/gh/komuW/meli/)              
However, I'm not making a tool to take docker-compose to the races.                   

# Build                   
`git clone git@github.com:komuW/meli.git`           
`go build -o meli cli/cli.go`           
`./meli -up -f /path/to/docker-compose-file.yml`                   


# TODO
- add better documentation(godoc)
- stabilise API(maybe)
