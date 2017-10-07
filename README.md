# meli

Meli is supposed to be a faster alternative to docker-compose. Faster in the sense that, Meli will try to pull as many services(docker containers) 
as it can in parallel.

Meli is a Swahili word meaning ship; so think of Meli as a ship carrying your docker contatiners safely across the treacherous container seas.

It's currently work in progress.

I only intend to support docker-compose version 3+; https://docs.docker.com/compose/compose-file/compose-versioning/


# Usage             
```shell
meli --help
Usage of meli:
  -up
    	Builds, re/creates, starts, and attaches to containers for a service.
        -d option runs containers in the background
  -v	Show version information.
  -version
    	Show version information.
```

```shell
meli -up 
2017/10/07 14:30:11 {"status":"Pulling from library/redis","id":"3.0-alpine"}
2017/10/07 14:30:11 {"status":"Digest: sha256:350469b395eac82395f9e59d7b7b90f7d23fe0838965e56400739dec3afa60de"}
2017/10/07 14:30:11 {"status":"Status: Image is up to date for redis:3.0-alpine"}
...
2017/10/07 14:30:12 u2017-10-07T11:30:12.720619075Z  1:M 07 Oct 11:30:12.720 * The server is now ready to accept connections on port 6379
```

# Build                   
`go build -o meli *.go`           
`cp meli testdata/`                 
`cd testdata`               
`./meli`                