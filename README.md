# meli

Meli is supposed to be a faster alternative to docker-compose. Faster in the sense that, Meli will try to pull as many services(docker containers) 
as it can in parallel.

Meli is a Swahili word meaning ship; so think of Meli as a ship carrying your docker containers safely across the treacherous container seas.

It's currently work in progress.

I only intend to support docker-compose version 3+; https://docs.docker.com/compose/compose-file/compose-versioning/           

Meli is NOT intended to replicate every feature of docker-compose, it is primarily intended to enable you to pull, build and run the services in your docker-compose file as fast as possible.          
If you want to exec in to a running container, use docker; if you want to run an adhoc command in a container, use docker; if you want..... you get the drift.


# Installing          
Very early test releases are available from the [releases page](https://github.com/komuW/meli/releases)          
Currently we only have a linux 64bit release, but you can build your own. See the build section below.

# Usage  
`meli --help`         
```bash
Usage of meli:
  -up
    	Builds, re/creates, starts, and attaches to containers for a service.
        -d option runs containers in the background
  -v	Show version information.
  -version
    	Show version information.
```

`meli -up`          
```bash 
2017/10/07 14:30:11 {"status":"Pulling from library/redis","id":"3.0-alpine"}
2017/10/07 14:30:11 {"status":"Digest: sha256:350469b395eac82395f9e59d7b7b90f7d23fe0838965e56400739dec3afa60de"}
2017/10/07 14:30:11 {"status":"Status: Image is up to date for redis:3.0-alpine"}
...
2017/10/07 14:30:12 u2017-10-07T11:30:12.720619075Z  1:M 07 Oct 11:30:12.720 * The server is now ready to accept connections on port 6379
```

# Build                   
`git clone git@github.com:komuW/meli.git`           
`go build -o meli *.go`           
`cp meli /dir/with/docker-compose-file/`                 
`cd /dir/with/docker-compose-file/`               
`./meli -up`                


# Benchmarks
Aaah, everyones' favorite vanity metric yet no one seems to know how to conduct one in a consistent and scientific manner.          
Take any results you see here with a large spoon of salt; They are unscientific and reek of all that is wrong with most developer benchmarks.             

Having made that disclaimer,                 

test machine:             
`lsb_release -a`
```bash
Distributor ID:	Ubuntu
Description:	Ubuntu 16.04 LTS
Release:	16.04
Codename:	xenial
No LSB modules are available.
```
`uname -ra`
```bash
4.4.0-96-generic #119-Ubuntu SMP Tue Sep 12 14:59:54 UTC 2017 x86_64 x86_64 x86_64 GNU/
```

docker-compose version:         
`docker-compose --version`
```bash
docker-compose version 1.16.1, build 6d1ac219
```

Meli version:   
```bash
git master branch head as of commit https://github.com/komuW/meli/commit/c27cb879a62ef4a61f7aef261b6b3d437090e4cc
``` 
NB: I haven't started making Meli versions since it's still early days.           

Benchmark test:           
[this docker-compose file](https://github.com/komuW/meli/blob/master/testdata/docker-compose.yml)

Benchmark script:               
for docker-compose:      
`docker ps -aq | xargs docker rm -f; docker system prune -af; /usr/bin/time -apv docker-compose up -d`        
for meli:                
`docker ps -aq | xargs docker rm -f; docker system prune -af; /usr/bin/time -apv meli -up -d`            

the above scripts were ran 3 times for each tool and an average taken. 

Benchmark results(average):                       

| tool           | User time(seconds) | Elapsed wall clock time(seconds) |
| :---           |     :---:          |          ---:                    |
| docker-compose | 1.61 sec           | 63.57 sec                        |
| meli           | 0.04 sec           | 28.43 sec                        |

Thus, meli appears to be 2.2 times faster than docker-compose(by wall clock time).       
There are still some low hanging fruits, performance wise, that I'll target in future.        
But I'm not making a tool to take docker-compose to the races.

